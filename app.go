package elio

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// elio option constants
const (
	// defaultLogLevel default log level
	defaultLogLevel = "debug"
	// defaultLogJson default json flag
	defaultLogJson = false
	// defaultLogShortCaller default short caller flag
	defaultLogShortCaller = false
	// defaultLogNoColor default log nocolor
	defaultLogNoColor = false
)

type run struct {
	svc    Service
	ctx    context.Context
	cancel context.CancelFunc
	ioHost *IoHost
}

// App app
type App struct {
	ctx      context.Context
	cancel   context.CancelFunc
	config   *Config
	envApp   EnvApp
	wg       sync.WaitGroup
	metrics  AppMetrics
	services []Service
	runs     []run
}

//// public scope
// String object to string
func (a *App) String() string {
	return fmt.Sprintf("App::%p", a)
}

// End call cancel callback
func (a *App) End() {
	a.cancel()

	for _, r := range a.runs {
		r.ioHost.End()
		r.cancel()
	}
}

// Wait wait app
func (a *App) Wait() {
	a.wg.Wait()
}

// Run run
const defaultAppIntervalMs = 100 * time.Millisecond

func (a *App) Run() {
	a.wg.Add(1)

	for _, s := range a.services {
		c, f := context.WithCancel(context.Background())
		if nil == s.OnInit(c, f) {
			a.runs = append(a.runs, run{
				svc:    s,
				ctx:    c,
				cancel: f,
			})
		}
	}

	intervalMs := defaultAppIntervalMs
	if true == a.config.Exists("elio.intervalMs") {
		intervalMs = time.Duration(a.config.GetInt("elio.intervalMs")) * time.Millisecond
	}

	a.wg.Add(len(a.runs))

	for i := 0; i < len(a.runs); i++ {
		//for _, r := range a.runs {
		// TODO: load service config in here r.svc.Name()
		c := ProvideConfigIo(a.runs[i].svc.Name(), a.config)
		//c.InURL = "0.0.0.0:7000"
		fmt.Printf("in.url:%v\n", c.InURL)

		var ioHost *IoHost
		var err error
		if ioHost, err = Host(c, a.runs[i].svc); nil != err {
			AppFatal().Str(LogObject, a.String()).
				Err(err).Msg("failed to echo service")
		} else {
			AppInfo().Str(LogObject, a.String()).
				Msgf("serve app:%s with url:%s", a.runs[i].svc.Name(), ioHost.Url)
			a.runs[i].ioHost = ioHost
		}

		go func(r run, i time.Duration) {
			defer func() {
				a.wg.Done()
				r.svc.OnExit()
			}()

			tickPrev := time.Now()
			ticker := time.NewTicker(i)
			for {
				select {
				case tick := <-ticker.C:
					r.svc.OnLoop(r.ioHost, tick, tick.Sub(tickPrev))
					tickPrev = tick
				case <-r.ctx.Done():
					return
				}
			}
		}(a.runs[i], intervalMs)
	}

	go func(i time.Duration) {
		defer func() {
			a.wg.Done()
		}()

		ticker := time.NewTicker(i)
		for {
			select {
			case <-ticker.C:
				// TODO: work in here
				//fmt.Printf(".")
			case <-a.ctx.Done():
				return
			}
		}
	}(intervalMs)
}

// Register register service
func (a *App) Register(s Service) {
	a.services = append(a.services, s)
}

// GetMetrics get app metrics
func (a *App) Metrics() *AppMetrics {
	return &a.metrics
}

// MetricCounter get metric counter
func (a *App) MetricCounter(name string) metrics.Counter {
	return a.Metrics().GetOrRegisterCounter(name)
}

// MetricGauge get metric gauge
func (a *App) MetricGauge(name string) metrics.Gauge {
	return a.Metrics().GetOrRegisterGauge(name)
}

// MetricMeter get metric meter
func (a *App) MetricMeter(name string) metrics.Meter {
	return a.Metrics().GetOrRegisterMeter(name)
}

//// private scope
// getContext get context
func (a *App) getContext() context.Context {
	return a.ctx
}

// init init
func (a *App) init() {
	a.config = NewConfig()
	path := flag.String("c", "", "a config path")
	flag.Parse()

	a.loadConfig(*path)

	a.initLog()

	a.ctx, a.cancel = context.WithCancel(context.Background())
}

// loadConfig load config
func (a *App) loadConfig(path string) {
	a.config.Load(path)

	ll, _ := a.config.GetStringOrDefault("elio.log.level", "debug")
	if l, e := zerolog.ParseLevel(strings.ToLower(ll)); nil == e {
		a.envApp.logLevel = l
	} else {
		a.envApp.logLevel = zerolog.DebugLevel
	}

	outs, _ := a.config.GetOrDefault("elio.log.out", "")
	switch outs.(type) {
	case []interface{}:
		{
			for _, v := range outs.([]interface{}) {
				a.envApp.logOuts = append(a.envApp.logOuts, v.(string))
			}
		}
	case string:
		os := strings.Split(outs.(string), ",")
		for _, v := range os {
			a.envApp.logOuts = append(a.envApp.logOuts, v)
		}
	default:
	}

	a.envApp.logJson, _ = a.config.GetBoolOrDefault("elio.log.json", defaultLogJson)
	a.envApp.logShortCaller, _ = a.config.GetBoolOrDefault("elio.log.shortCaller", defaultLogShortCaller)
	a.envApp.logNoColor, _ = a.config.GetBoolOrDefault("elio.log.nocolor", defaultLogNoColor)
}

// initLog init log
func (a *App) initLog() {
	zerolog.TimestampFieldName = "logtime"
	zerolog.TimeFieldFormat = time.RFC3339Nano
	if true == a.envApp.logShortCaller {
		zerolog.CallerMarshalFunc = func(file string, line int) string {
			f := filepath.Base(file)
			return f + ":" + strconv.Itoa(line)
		}
	}

	var writers []io.Writer
	if 0 < len(a.envApp.logOuts) {
		for _, o := range a.envApp.logOuts {
			if "stdout" == o {
				if true == a.envApp.logJson {
					writers = append(writers, os.Stdout)
				} else {
					o := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: a.envApp.logNoColor}
					// o.FormatLevel = func(i interface{}) string {
					// 	return strings.ToUpper(fmt.Sprintf("%-6s", i))
					// }
					// o.FormatMessage = func(i interface{}) string {
					// 	return fmt.Sprintf("* %s *", i)
					// }
					// o.FormatFieldName = func(i interface{}) string {
					// 	return fmt.Sprintf("%s:", i)
					// }
					//o.FormatFieldValue = func(i interface{}) string {
					//	return strings.ToUpper(fmt.Sprintf("%s", i))
					//}
					// if true == a.envApp.logShortCaller {
					// 	o.FormatCaller = func(i interface{}) string {
					// 		t := fmt.Sprintf("%s", i)
					// 		s := strings.Split(t, ":")
					// 		if 2 != len(s) {
					// 			return t
					// 		}
					// 		f := filepath.Base(s[0])
					// 		return f + ":" + s[1]
					// 	}
					// }
					writers = append(writers, o)
				}

			} else if "stderr" == o {
				writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339, NoColor: a.envApp.logNoColor})
			} else {
				if f, err := os.OpenFile(o, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); nil == err {
					// Add file and line number to log
					writers = append(writers, f)
				}
			}
		}

	} else {
		//writers = append(writers, zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
		o := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		// o.FormatLevel = func(i interface{}) string {
		// 	return strings.ToUpper(fmt.Sprintf("%-6s", i))
		// }
		// o.FormatMessage = func(i interface{}) string {
		// 	return fmt.Sprintf("* %s *", i)
		// }
		// o.FormatFieldName = func(i interface{}) string {
		// 	return fmt.Sprintf("%s:", i)
		// }
		// o.FormatFieldValue = func(i interface{}) string {
		// 	return strings.ToUpper(fmt.Sprintf("%s", i))
		// }
		// if true == a.envApp.logShortCaller {
		// 	o.FormatCaller = func(i interface{}) string {
		// 		t := fmt.Sprintf("%s", i)
		// 		s := strings.Split(t, ":")
		// 		if 2 != len(s) {
		// 			return t
		// 		}
		// 		f := filepath.Base(s[0])
		// 		return f + ":" + s[1]
		// 	}
		// }
		writers = append(writers, o)
	}

	zerolog.SetGlobalLevel(a.envApp.logLevel)

	log.Logger = log.With().Caller().Logger()
	if 0 < len(writers) {
		log.Logger = log.Output(io.MultiWriter(writers...))
	}
}

// NewApp new app
func NewApp() *App {
	a := new(App)
	if nil != a {
		a.init()
	}
	return a
}
