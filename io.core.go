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

// IoCore IoCore
type IoCore struct {
	Listener    *Listener //net.Listener
	Config      ConfigIo
	Host        *IoHost
	Service     Service
	io          Io
	InAddr      atomic.String //*net.TCPAddr
	InCount     atomic.Int32
	sessionCmap cmap.ConcurrentMap
	ctx         context.Context
	cancel      context.CancelFunc
	action      Action
}

// NewIoCore new server
func NewIoCore(c ConfigIo, s Service) *IoCore {
	if core := new(IoCore); nil != core {
		core.Service = s
		core.ctx, core.cancel = context.WithCancel(context.Background())
		core.Listener = newListener(c.InModel, c.InReusePort)

		if i := GenIO(c.InModel); nil != i {
			i.SetIoCore(core)
			core.Config = c
			core.Init()
			return core
		}
	}

	return nil
}

func (c *IoCore) String() string {
	return fmt.Sprintf("IoCore::%p", c)
}

// Run run IoCore
func (c *IoCore) Run(addr *net.TCPAddr) (ok bool) {
	AppInfo().Str(LogObject, c.String()).
		Msgf("IoCore.run url:%s begin", c.Config.InURL)

	if ok = c.Listen(addr); ok {
		c.io.Run()
	}

	return ok
}

// Listen listen IoCore
func (c *IoCore) Listen(addr *net.TCPAddr) (ok bool) {
	return c.io.Listen(addr.String())
}

// Shut shut
func (c *IoCore) Shut() {
	c.io.Shut()
}

// End end
func (c *IoCore) End() {
	defer func() {
		c.Host.Wg.Done()
	}()

	//c.cancel()
	c.io.End()

	//c.Host.Terminate()
}

// Shutdown shutdown
func (c *IoCore) Shutdown(n *Session, how int) error {
	return c.io.Shutdown(n, how)
}

// Terminate terminate
func (c *IoCore) Terminate() {
	AppInfo().Str(LogObject, c.String()).Msg("terminate IoCore")

	c.io.Shut()
	c.io.End()
}

// SafeTerminate safe terminate
//func (c *IoCore) SafeTerminate() {
//	c.flagSafeExit.Store(true)
//}

// Init init
func (c *IoCore) Init() {
	c.sessionCmap = cmap.New()

	//c.flagExit.Store(false)
	//c.flagSafeExit.Store(false)
}

// GetIO get IO
func (c *IoCore) GetIo() Io {
	return c.io
}

// SetIO set IO
func (c *IoCore) SetIo(io Io) {
	c.io = io
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
