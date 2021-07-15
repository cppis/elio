package elio

import (
	"fmt"
	"net"
	"sync"
)

// IoHost host
type IoHost struct {
	IoCores []*IoCore
	Wg      *sync.WaitGroup
	Url     *net.TCPAddr
}

// NewIoHost new host
func NewIoHost() (s *IoHost) {
	if s = new(IoHost); nil != s {
		s.Wg = new(sync.WaitGroup)
	}

	return s
}

// String string
func (s *IoHost) String() string {
	return fmt.Sprintf("IoHost::%p", s)
}

// Wait wait
func (s *IoHost) Wait() (err error) {
	s.Wg.Wait()

	AppInfo().Str(LogObject, s.String()).Msgf("all services end")

	return err
}

// End end
func (s *IoHost) End() {
	for _, c := range s.IoCores {
		c.End()
	}
}

// Terminate terminate
func (s *IoHost) Terminate(safe bool) {
	for _, c := range s.IoCores {
		//if true == safe {
		//	c.SafeTerminate()
		//} else {
		c.Terminate()
		//}
	}
}

// Host host
func Host(config ConfigIo, service Service) (host *IoHost, err error) {
	host = NewIoHost()
	host.IoCores = make([]*IoCore, config.InCount)

	for i := 0; i < config.InCount; i++ {
		core := ProvideIoCore(config, service)
		host.IoCores[i] = core

		core.Host = host

		var addr *net.TCPAddr
		if addr, err = net.ResolveTCPAddr("tcp", core.Config.InURL); nil == err {
			core.Run(addr)
			host.Url = addr
		}
	}

	return host, err
}
