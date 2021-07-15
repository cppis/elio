package elio

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	// LogDump dump
	LogDump = "dump"
	// LogName name
	LogName = "name"
	// LogProject project
	LogProject = "project"
	// LogService service
	LogService = "service"
	// LogStage stage
	LogStage = "stage"
	// LogSuid server unique id
	LogSuid = "suid"
	// LogHost host
	LogHost = "host"
	// LogObject object
	LogObject = "object"
	// LogSession session
	LogSession = "session"
	// LogUuid user unique id
	LogUuid = "uuid"
	// LogIP ip
	LogIP = "ip"
	// LogCode code
	LogCode = "code"
	// LogMatchKey match key
	LogMatchKey = "matchkey"
	// LogState state
	LogState = "state"
	// LogWalltime walltime
	LogWalltime = "walltime"
	// LogErrorCode errorcode
	LogErrorCode = "errorcode"
)

// InitLog init log
func InitLog(level string, outs []string, json bool) {
	zerolog.TimestampFieldName = "logtime"
	zerolog.TimeFieldFormat = time.RFC3339Nano

	//fmt.Printf("app.env log level: %s\n", h.envVar.logLevel)
	if l, e := zerolog.ParseLevel(strings.ToLower(level)); nil == e {
		zerolog.SetGlobalLevel(l)
	}

	var writers []io.Writer
	if 0 < len(outs) {
		for _, o := range outs {
			if "stdout" == o {
				if true == json {
					writers = append(writers, os.Stdout)
				} else {
					writers = append(writers, zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
				}

			} else if "stderr" == o {
				writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
			} else {
				if f, err := os.OpenFile(o, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); nil == err {
					// Add file and line number to log
					writers = append(writers, f)
				}
			}
		}

	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

		o := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: false}
		//o.FormatLevel = func(i interface{}) string {
		//	return strings.ToUpper(fmt.Sprintf("%-6s", i))
		//}
		//o.FormatMessage = func(i interface{}) string {
		//	return fmt.Sprintf("* %s *", i)
		//}
		//o.FormatFieldName = func(i interface{}) string {
		//	return fmt.Sprintf("%s:", i)
		//}
		o.FormatFieldValue = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("%s", i))
		}
		o.FormatCaller = func(i interface{}) string {
			t := fmt.Sprintf("%s", i)
			s := strings.Split(t, ":")
			if 2 != len(s) {
				return t
			}
			f := filepath.Base(s[0])
			return f + ":" + s[1]
		}
		writers = append(writers, o)
	}

	log.Logger = log.With().Caller().Logger()
	if 0 < len(writers) {
		log.Logger = log.Output(io.MultiWriter(writers...))
	}
}

// WithLevel with level
func WithLevel(level zerolog.Level) *zerolog.Event {
	return log.WithLevel(level)
}

// PanicEnabled panic enabled
func PanicEnabled() bool {
	return log.Panic().Enabled()
}

// FatalEnabled fatal enabled
func FatalEnabled() bool {
	return log.Fatal().Enabled()
}

// ErrorEnabled error enabled
func ErrorEnabled() bool {
	return log.Error().Enabled()
}

// WarnEnabled debug enabled
func WarnEnabled() bool {
	return log.Warn().Enabled()
}

// InfoEnabled info enabled
func InfoEnabled() bool {
	return log.Info().Enabled()
}

// DebugEnabled debug enabled
func DebugEnabled() bool {
	return log.Debug().Enabled()
}

// TraceEnabled debug enabled
func TraceEnabled() bool {
	return log.Trace().Enabled()
}

// LogPanic app log panic
func LogPanic(name string) *zerolog.Event {
	return log.Panic().
		Str(LogName, name).
		Str(LogProject, LogParams().GetProject()).
		Str(LogService, LogParams().GetService()).
		Str(LogHost, LogParams().GetHostname())
}

// LogFatal app log fatal
func LogFatal(name string) *zerolog.Event {
	return log.Fatal().
		Str(LogName, name).
		Str(LogProject, LogParams().GetProject()).
		Str(LogService, LogParams().GetService()).
		Str(LogHost, LogParams().GetHostname())
}

// LogError app log error
func LogError(name string) *zerolog.Event {
	return log.Error().
		Str(LogName, name).
		Str(LogProject, LogParams().GetProject()).
		Str(LogService, LogParams().GetService()).
		Str(LogHost, LogParams().GetHostname())
}

// LogWarn app log warn
func LogWarn(name string) *zerolog.Event {
	return log.Warn().
		Str(LogName, name).
		Str(LogProject, LogParams().GetProject()).
		Str(LogService, LogParams().GetService()).
		Str(LogHost, LogParams().GetHostname())
}

// LogInfo app log info
func LogInfo(name string) *zerolog.Event {
	return log.Info().
		Str(LogName, name).
		Str(LogProject, LogParams().GetProject()).
		Str(LogService, LogParams().GetService()).
		Str(LogHost, LogParams().GetHostname())
}

// LogDebug app log debug
func LogDebug(name string) *zerolog.Event {
	return log.Debug().
		Str(LogName, name).
		Str(LogProject, LogParams().GetProject()).
		Str(LogService, LogParams().GetService()).
		Str(LogHost, LogParams().GetHostname())
}

// LogTrace app log trace
func LogTrace(name string) *zerolog.Event {
	return log.Trace().
		Str(LogName, name).
		Str(LogProject, LogParams().GetProject()).
		Str(LogService, LogParams().GetService()).
		Str(LogHost, LogParams().GetHostname())
}

// AppPanic app log panic
func AppPanic() *zerolog.Event { return LogPanic("app") }

// AppFatal app log fatal
func AppFatal() *zerolog.Event { return LogFatal("app") }

// AppError app log error
func AppError() *zerolog.Event { return LogError("app") }

// AppWarn app log warn
func AppWarn() *zerolog.Event { return LogWarn("app") }

// AppInfo app log info
func AppInfo() *zerolog.Event { return LogInfo("app") }

// AppDebug app log debug
func AppDebug() *zerolog.Event { return LogDebug("app") }

// AppTrace app log trace
func AppTrace() *zerolog.Event { return LogTrace("app") }

// PacketPanic packet log panic
func PacketPanic() *zerolog.Event { return LogPanic("packet") }

// PacketFatal packet log fatal
func PacketFatal() *zerolog.Event { return LogFatal("packet") }

// PacketError packet log error
func PacketError() *zerolog.Event { return LogError("packet") }

// PacketWarn packet log warn
func PacketWarn() *zerolog.Event { return LogWarn("packet") }

// PacketInfo packet log info
func PacketInfo() *zerolog.Event { return LogInfo("packet") }

// PacketDebug packet log debug
func PacketDebug() *zerolog.Event { return LogDebug("packet") }

// PacketTrace packet log trace
func PacketTrace() *zerolog.Event { return LogTrace("packet") }

// DumpPanic packet log panic
func DumpPanic() *zerolog.Event { return LogPanic("dump") }

// DumpFatal packet log fatal
func DumpFatal() *zerolog.Event { return LogFatal("dump") }

// DumpError packet log error
func DumpError() *zerolog.Event { return LogError("dump") }

// DumpWarn packet log warn
func DumpWarn() *zerolog.Event { return LogWarn("dump") }

// DumpInfo packet log info
func DumpInfo() *zerolog.Event { return LogInfo("dump") }

// DumpDebug packet log debug
func DumpDebug() *zerolog.Event { return LogDebug("dump") }

// DumpTrace packet log trace
func DumpTrace() *zerolog.Event { return LogTrace("dump") }
