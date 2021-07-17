package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cppis/elio"
)

// Echo echo service
type Echo struct {
	app    *elio.App
	ctx    context.Context
	cancel context.CancelFunc
	prev   time.Time
}

func NewEcho(app *elio.App) *Echo {
	e := new(Echo)
	if nil != e {
		e.app = app
	}
	return e
}

// String string
func (e *Echo) String() string {
	return fmt.Sprintf("Echo::%p", e)
}

func (e *Echo) Name() string {
	return "echo"
}

func (e *Echo) OnInit(ctx context.Context, cancel context.CancelFunc) error {
	e.ctx = ctx
	e.cancel = cancel

	fmt.Printf("%s on init\n", e.Name())
	return nil
}

func (e *Echo) OnExit() {
	fmt.Printf("%s on exit\n", e.Name())
}

func (e *Echo) OnOpen(s *elio.Session) error {
	fmt.Printf("o")

	elio.AppDebug().Str(elio.LogObject, e.String()).
		Str(elio.LogSession, s.String()).Msgf("service:%s on.open", e.Name())

	return nil
}

func (e *Echo) OnClose(s *elio.Session, err error) {
	fmt.Printf("c")

	elio.AppDebug().Str(elio.LogObject, e.String()).
		Str(elio.LogSession, s.String()).Msgf("service:%s on.close", e.Name())
}

func (e *Echo) OnError(s *elio.Session, err error) {
	fmt.Printf("e")

	elio.AppError().Str(elio.LogObject, e.String()).
		Str(elio.LogSession, s.String()).Msgf("service:%s on.error", e.Name())
}

func (e *Echo) OnRead(s *elio.Session, in []byte) (processed int) {
	fmt.Printf("+%d", len(in))

	s.Write(in)

	if 'q' == in[0] {
		elio.Elio().End()
	}

	return processed
}

func (e *Echo) OnWrite(s *elio.Session, out []byte) {
	fmt.Printf("-%d", len(out))
}

const (
	// defaultFetchLimit default fetch limit
	defaultFetchLimit int = 2000
)

func (e *Echo) OnLoop(host *elio.IoHost, t time.Time, d time.Duration) {
	//fmt.Printf("on loop with delta:%v\n", d)
	if t.Sub(e.prev) > 10*time.Second {
		fmt.Printf("e")
	}

	//host.RunDivision(t, r.callbackDivision)

	_, _ = host.Dispatching(t, defaultFetchLimit)

	e.prev = t
}
