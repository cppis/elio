package elio

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strings"

	"time"

	cmap "github.com/orcaman/concurrent-map"
	"go.uber.org/atomic"
)

const (
	// DefaultInterval default interval
	DefaultInterval time.Duration = 20 * time.Millisecond
)

// https://blog.alexellis.io/inject-build-time-vars-golang/

// Action is an action that occurs after the completion of an event.
type Action int

const (
	// None indicates that no action should occur following an event.
	None Action = iota
	// Listen listen.
	Listen
	// Shut close the listen.
	Shut
	// SafeTermiate close the listen and exit when no user.
	SafeTermiate
	// Termiate terminate the server.
	Termiate
	// TermiateAfter5m shutdown after 5 minute.
	TermiateAfter5m
	// CloseAll disconnect all connections.
	CloseAll
	// Close disconnect connection.
	Close
)

// Io IO
type Io struct {
	Listener    *Listener //net.Listener
	Config      ConfigIo
	Host        *IoHost
	Service     Service
	ioModel     IoModel
	InAddr      atomic.String //*net.TCPAddr
	InCount     atomic.Int32
	sessionCmap cmap.ConcurrentMap
	ctx         context.Context
	cancel      context.CancelFunc
	action      Action
}

// NewIo new server
func NewIo(h *IoHost, c ConfigIo, s Service) *Io {
	if io := new(Io); nil != io {
		io.Host = h
		io.Service = s
		io.ctx, io.cancel = context.WithCancel(context.Background())
		io.Listener = newListener(c.InModel, c.InReusePort)

		if i := GenIO(c.InModel); nil != i {
			i.SetIo(io)
			io.Config = c
			io.Init()
			return io
		}
	}

	return nil
}

func (i *Io) String() string {
	return fmt.Sprintf("Io::%p", i)
}

// Run run Io
func (i *Io) Run(addr *net.TCPAddr) (ok bool) {
	AppInfo().Str(LogObject, i.String()).
		Msgf("Io.run url:%s begin", i.Config.InURL)

	if ok = i.Listen(addr); ok {
		i.ioModel.Run()
	}

	return ok
}

// Listen listen Io
func (i *Io) Listen(addr *net.TCPAddr) (ok bool) {
	return i.ioModel.Listen(addr.String())
}

// Shut shut
func (i *Io) Shut() {
	i.ioModel.Shut()
}

// End end
func (i *Io) End() {
	defer func() {
		i.Host.Wg.Done()
	}()

	//i.cancel()
	i.ioModel.End()

	//i.Host.Terminate()
}

// Shutdown shutdown
func (i *Io) Shutdown(n *Session, how int) error {
	return i.ioModel.Shutdown(n, how)
}

// Terminate terminate
func (i *Io) Terminate() {
	AppInfo().Str(LogObject, i.String()).Msg("terminate Io")

	i.ioModel.Shut()
	i.ioModel.End()
}

// Init init
func (i *Io) Init() {
	i.sessionCmap = cmap.New()
}

// GetIoModel get IO model
func (i *Io) GetIoModel() IoModel {
	return i.ioModel
}

// SetIoModel set IO model
func (i *Io) SetIoModel(ioModel IoModel) {
	i.ioModel = ioModel
}

// GetBaseAndConfig get base and config
func GetBaseAndConfig() (base string, config string) {
	//if 1 < len(os.Args) {
	//	base = os.Args[0]
	//	config = GetBasename(os.Args[1])
	//} else {
	base = GetBasename(os.Args[0])
	var builder strings.Builder
	builder.WriteString(base)
	builder.WriteString(".json")
	config = builder.String()
	//}

	return base, config
}

// StartUp start up
func StartUp() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	if err := SetLimit(); nil != err {
		AppError().Err(err).Msg("failed to set.limit")
	}
}

// StopDown stop down
func StopDown(exit func()) {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, SigTerm)
	signal.Notify(gracefulStop, SigInt)

	go func() {
		sig := <-gracefulStop
		AppInfo().Msgf("caught sig:%+v\n", sig)
		exit()
	}()
}
