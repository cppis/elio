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
	s := new(Echo)
	if nil != s {
		s.app = app
	}
	return s
}

// String string
func (s *Echo) String() string {
	return fmt.Sprintf("Echo::%p", s)
}

func (s *Echo) Name() string {
	return "echo"
}

func (s *Echo) OnInit(ctx context.Context, cancel context.CancelFunc) error {
	s.ctx = ctx
	s.cancel = cancel

	elio.AppDebug().Str(elio.LogObject, s.String()).Msg("on init")
	return nil
}

func (s *Echo) OnExit() {
	elio.AppDebug().Str(elio.LogObject, s.String()).Msg("on exit")
}

func (s *Echo) OnOpen(sn *elio.Session) error {
	fmt.Printf("o")

	elio.AppDebug().Str(elio.LogObject, s.String()).
		Str(elio.LogSession, sn.String()).Msgf("service:%s on.open", s.Name())

	return nil
}

func (s *Echo) OnClose(sn *elio.Session, err error) {
	fmt.Printf("c")

	elio.AppDebug().Str(elio.LogObject, s.String()).
		Str(elio.LogSession, sn.String()).Msgf("service:%s on.close", s.Name())
}

func (s *Echo) OnError(sn *elio.Session, err error) {
	//fmt.Printf("e")

	elio.AppError().Str(elio.LogObject, s.String()).
		Str(elio.LogSession, sn.String()).Msgf("service:%s on.error", s.Name())
}

func (s *Echo) OnRead(sn *elio.Session, in []byte) (processed int) {
	fmt.Printf("+%d", len(in))

	sn.Write(in)

	if 'q' == in[0] {
		elio.Elio().End()
	} else if '?' == in[0] {
		processed = -1
	}

	return processed
}

func (s *Echo) OnWrite(sn *elio.Session, out []byte) {
	fmt.Printf("-%d", len(out))
}

const (
	// defaultFetchLimit default fetch limit
	defaultFetchLimit int = 2000
)

func (s *Echo) OnLoop(host *elio.IoHost, t time.Time, d time.Duration) {
	//fmt.Printf("on loop with delta:%v\n", d)
	//if t.Sub(e.prev) > 10*time.Second {
	//	fmt.Printf("e")
	//}

	//host.RunDivision(t, r.callbackDivision)

	_, _ = host.Dispatching(t, defaultFetchLimit)

	s.prev = t
}
