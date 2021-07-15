package elio

import (
	"os"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// FlagNone default none flag
	FlagNone int = iota
	// FlagTerminate terminate flag
	FlagTerminate
	// FlagSafeTerminate safe terminate flag
	FlagSafeTerminate
	// FlagStop stop listen flag
	FlagStop
	// FlagAllout disconnect all flag
	FlagAllout
	// FlagTerminateAfter5m exit after 5 minute flag
	FlagTerminateAfter5m
)

const (
	// Duration5m duration 5 minute
	Duration5m time.Duration = 5 * time.Second
)

// logParams logParams object
type logParams struct {
	project  atomic.Value
	stage    atomic.Value
	service  atomic.Value
	suid     atomic.Value
	hostName atomic.Value
	cpuProf  atomic.Value
	memProf  atomic.Value
}

const (
	// ProdLogger production logger
	ProdLogger = "prod"
	// DevLogger development logger
	DevLogger = "dev"
	// ExamLogger example logger
	ExamLogger = "dev"
)

// init init
func (a *logParams) init() {
	a.project.Store("")
	a.stage.Store("")
	a.service.Store("")
	a.suid.Store("")
	a.hostName.Store("")
	a.cpuProf.Store("")
	a.memProf.Store("")
}

// GetHostname get hostname
func (p *logParams) GetHostname() (hostName string) {
	hostName = p.hostName.Load().(string)
	if "" == hostName {
		var err error
		if hostName, err = os.Hostname(); nil == err {
			p.hostName.Store(hostName)
		}
	}

	return hostName
}

// GetProject get project
func (p *logParams) GetProject() string {
	return p.project.Load().(string)
}

// SetProject set project
func (p *logParams) SetProject(v string) {
	p.project.Store(v)
}

// GetStage get stage
func (p *logParams) GetStage() string {
	return p.stage.Load().(string)
}

// SetStage set stage
func (p *logParams) SetStage(v string) {
	p.stage.Store(v)
}

// GetService get service
func (p *logParams) GetService() string {
	return p.service.Load().(string)
}

// SetService set service
func (p *logParams) SetService(v string) {
	p.service.Store(v)
}

// GetSuid get server unique id
func (p *logParams) GetSuid() string {
	return p.suid.Load().(string)
}

// SetSuid set server unique id
func (p *logParams) SetSuid(v string) {
	p.suid.Store(v)
}

// GetCPUProf get cpu prof file
func (p *logParams) GetCPUProf() string {
	return p.cpuProf.Load().(string)
}

// SetCPUProf set cpu prof file
func (p *logParams) SetCPUProf(v string) {
	p.cpuProf.Store(v)
}

// GetMemProf get mem prof file
func (p *logParams) GetMemProf() string {
	return p.memProf.Load().(string)
}

// SetMemProf set cpu prof file
func (p *logParams) SetMemProf(v string) {
	p.memProf.Store(v)
}

var instanceAp *logParams
var onceAp sync.Once

// LogParams get log params
func LogParams() *logParams {
	onceAp.Do(func() {
		instanceAp = &logParams{}
		instanceAp.init()
	})
	return instanceAp
}
