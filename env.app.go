package elio

// EnvApp env app
//	ELIO_LOG_LEVEL
//	ELIO_LOG_OUTS
//	ELIO_LOG_JSON
//	ELIO_LOG_NOCOLOR
//	ELIO_IN_METRIC
//	ELIO_IN_PPROF
type EnvApp struct {
	logLevel string
	logOuts  []string
	logJson  bool
	logShortCallser bool
	logNoColor  bool
	inMetric    int
	inPprof     int
}

func (e *EnvApp) init() {
}
