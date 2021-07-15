package elio

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-metrics/exp"
)

const (
	// CaptureInterval capture interval
	CaptureInterval = 100 * time.Millisecond
	// StatsDInterval capture interval
	StatsDInterval = 2500 * time.Millisecond
	// StatsDAddress statsd address
	StatsDAddress = "127.0.0.1:8125"
)

// MetricsApp app metrics
type MetricsApp struct {
	Registry metrics.Registry
}

// String object to string
func (m *MetricsApp) String() string {
	return fmt.Sprintf("MetricsApp::%p", m)
}

// Init init
func (m *MetricsApp) Init(inMetric int) {
	m.Registry = metrics.NewRegistry() //metrics.DefaultRegistry	//

	// debug gc stats
	metrics.RegisterDebugGCStats(m.Registry)
	go metrics.CaptureDebugGCStats(m.Registry, CaptureInterval)

	// runtime mem stats
	metrics.RegisterRuntimeMemStats(m.Registry)
	go metrics.CaptureRuntimeMemStats(m.Registry, CaptureInterval)

	exp.Exp(m.Registry)

	//if addr, err := net.ResolveUDPAddr("udp", StatsD_Address); nil != err {
	//	AppError().Str(elf.LogObject, m.String()).
	//	Err(err).Msg("failed to resolve statd address")
	//
	//} else {
	//	p := fmt.Sprintf("%s-%s-%s", elf.GetAppParams().GetProject(), elf.GetAppParams().GetStage(), elf.GetAppParams().GetSuid())

	//	go statsd.StatsD(m.Registry, StatsD_Interval, p, addr)
	//}

	if 0 != inMetric {
		// metric web service:
		//	http://localhost:{in.metric}/debug/metrics
		go http.ListenAndServe(fmt.Sprintf(":%d", inMetric), nil)
	}
}

// GetOrRegisterCounter get or register counter
/*//
// Counters hold an int64 value that can be incremented and decremented.
type Counter interface {
 	Clear()
 	Count() int64
 	Dec(int64)
 	Inc(int64)
 	Snapshot() Counter
}
//*/
func (m *MetricsApp) GetOrRegisterCounter(name string) metrics.Counter {
	return m.Registry.GetOrRegister(name, metrics.NewCounter).(metrics.Counter)
}

// GetOrRegisterEwma get or register EWMA
/*//
// EWMAs continuously calculate an exponentially-weighted moving average
// based on an outside source of clock ticks.
type EWMA interface {
	Rate() float64
	Snapshot() EWMA
	Tick()
	Update(int64)
}
//*/
func (m *MetricsApp) GetOrRegisterEwma(name string) metrics.EWMA {
	return m.Registry.GetOrRegister(name, metrics.NewEWMA).(metrics.EWMA)
}

// GetOrRegisterGauge get or register gauge
/*//
// Gauges hold an int64 value that can be set arbitrarily.
type Gauge interface {
	Snapshot() Gauge
	Update(int64)
	Value() int64
}
//*/
func (m *MetricsApp) GetOrRegisterGauge(name string) metrics.Gauge {
	return m.Registry.GetOrRegister(name, metrics.NewGauge).(metrics.Gauge)
}

// GetOrRegisterGaugeFloat64 get or register gauge float64
/*//
// GaugeFloat64s hold a float64 value that can be set arbitrarily.
type GaugeFloat64 interface {
	Snapshot() GaugeFloat64
	Update(float64)
	Value() float64
}
//*/
func (m *MetricsApp) GetOrRegisterGaugeFloat64(name string) metrics.GaugeFloat64 {
	return m.Registry.GetOrRegister(name, metrics.NewGaugeFloat64()).(metrics.GaugeFloat64)
}

// // GetOrRegisterHealthcheck get or register health check
// /*//
// // Healthchecks hold an error value describing an arbitrary up/down status.
// type Healthcheck interface {
// 	Check()
// 	Error() error
// 	Healthy()
// 	Unhealthy(error)
// }
// //*/
// func (m *MetricsApp) GetOrRegisterHealthcheck(name string, check func(h metrics.Healthcheck)) metrics.Healthcheck {
// 	return m.Registry.GetOrRegister(name, metrics.NewHealthcheck(check)).(metrics.Healthcheck)
// }

// GetOrRegisterHistogram get or register histogram
/*//
// Histograms calculate distribution statistics from a series of int64 values.
type Histogram interface {
	Clear()
	Count() int64
	Max() int64
	Mean() float64
	Min() int64
	Percentile(float64) float64
	Percentiles([]float64) []float64
	Sample() Sample
	Snapshot() Histogram
	StdDev() float64
	Sum() int64
	Update(int64)
	Variance() float64
}
//*/
func (m *MetricsApp) GetOrRegisterHistogram(name string, s metrics.Sample) metrics.Histogram {
	return m.Registry.GetOrRegister(name, func() metrics.Histogram { return metrics.NewHistogram(s) }).(metrics.Histogram)
}

// GetOrRegisterMeter get or register meter
/*//
// Meters count events to produce exponentially-weighted moving average rates
// at one-, five-, and fifteen-minutes and a mean rate.
type Meter interface {
	Count() int64
	Mark(int64)
	Rate1() float64
	Rate5() float64
	Rate15() float64
	RateMean() float64
	Snapshot() Meter
	Stop()
}
//*/
func (m *MetricsApp) GetOrRegisterMeter(name string) metrics.Meter {
	return m.Registry.GetOrRegister(name, metrics.NewMeter).(metrics.Meter)
}

// GetOrRegisterTimer get or register timer
/*//
// Timers capture the duration and rate of events.
type Timer interface {
	Count() int64
	Max() int64
	Mean() float64
	Min() int64
	Percentile(float64) float64
	Percentiles([]float64) []float64
	Rate1() float64
	Rate5() float64
	Rate15() float64
	RateMean() float64
	Snapshot() Timer
	StdDev() float64
	Stop()
	Sum() int64
	Time(func())
	Update(time.Duration)
	UpdateSince(time.Time)
	Variance() float64
}
//*/
func (m *MetricsApp) GetOrRegisterTimer(name string) metrics.Timer {
	return m.Registry.GetOrRegister(name, metrics.NewTimer).(metrics.Timer)
}
