package elio

import "github.com/rs/zerolog"

// EnvApp env app
//	ELIO_LOG_LEVEL
//	ELIO_LOG_OUTS
//	ELIO_LOG_JSON
//	ELIO_LOG_NOCOLOR
//	ELIO_IN_METRIC
//	ELIO_IN_PPROF
type EnvApp struct {
	logLevel       zerolog.Level
	logOuts        []string
	logJson        bool
	logShortCaller bool
	logNoColor     bool
	inMetric       int
	inPprof        int
}

func (e *EnvApp) init() {
}
