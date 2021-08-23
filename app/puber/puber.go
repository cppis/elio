package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cppis/elio"
)

// Puber service
type Puber struct {
	app    *elio.App
	ctx    context.Context
	cancel context.CancelFunc
	prev   time.Time
}

func NewPuber(app *elio.App) *Puber {
	s := new(Puber)
	if nil != s {
		s.app = app
	}
	return s
}

// String string
func (s *Puber) String() string {
	return fmt.Sprintf("Puber::%p", s)
}

func (s *Puber) Name() string {
	return "puber"
}

func (s *Puber) OnInit(ctx context.Context, cancel context.CancelFunc) error {
	s.ctx = ctx
	s.cancel = cancel

	mqttUrl, _ := s.app.Config().GetStringOrDefault(fmt.Sprintf("%s.mqtt.url", s.Name()), "")
	fmt.Printf("%s", mqttUrl)

	elio.AppDebug().Str(elio.LogObject, s.String()).Msg("on init")
	return nil
}

func (s *Puber) OnExit() {
	elio.AppDebug().Str(elio.LogObject, s.String()).Msg("on exit")
}

func (s *Puber) OnOpen(sn *elio.Session) error {
	fmt.Printf("o")

	elio.AppDebug().Str(elio.LogObject, s.String()).
		Str(elio.LogSession, sn.String()).Msgf("service:%s on.open", s.Name())

	return nil
}

func (s *Puber) OnClose(sn *elio.Session, err error) {
	fmt.Printf("c")

	elio.AppDebug().Str(elio.LogObject, s.String()).
		Str(elio.LogSession, sn.String()).Msgf("service:%s on.close", s.Name())
}

func (s *Puber) OnError(sn *elio.Session, err error) {
	//fmt.Printf("e")

	elio.AppError().Str(elio.LogObject, s.String()).
		Str(elio.LogSession, sn.String()).Msgf("service:%s on.error", s.Name())
}

func (s *Puber) OnRead(sn *elio.Session, in []byte) (processed int) {
	fmt.Printf("+%d", len(in))

	sn.Write(in)

	if 'q' == in[0] {
		elio.Elio().End()
	}

	return processed
}

func (s *Puber) OnWrite(sn *elio.Session, out []byte) {
	fmt.Printf("-%d", len(out))
}

const (
	// defaultFetchLimit default fetch limit
	defaultFetchLimit int = 2000
)

func (s *Puber) OnLoop(host *elio.IoHost, t time.Time, d time.Duration) {
	//fmt.Printf("on loop with delta:%v\n", d)
	//if t.Sub(e.prev) > 10*time.Second {
	//	fmt.Printf("e")
	//}

	//host.RunDivision(t, r.callbackDivision)

	_, _ = host.Dispatching(t, defaultFetchLimit)

	s.prev = t
}
