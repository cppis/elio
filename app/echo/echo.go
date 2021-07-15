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

func (e *Echo) Name() string {
	return "echo"
}

func (e *Echo) OnListen(i *elio.IoCore) {
	fmt.Printf("%s on listen\n", e.Name())
}

func (e *Echo) OnShut(i *elio.IoCore) {
	fmt.Printf("%s on shut\n", e.Name())
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
	fmt.Printf("%s on open\n", e.Name())
	return nil
}

func (e *Echo) OnClose(s *elio.Session, err error) {
	fmt.Printf("%s on close\n", e.Name())
}

func (e *Echo) OnError(s *elio.Session, err error) {
	fmt.Printf("on error\n")
}

func (e *Echo) OnRead(s *elio.Session, in []byte) int {
	fmt.Printf("on read %d\n", len(in))
	e.app.End()
	//s.GetIoCore().End()
	return len(in)
}

func (e *Echo) OnWrite(s *elio.Session, out []byte) {
	fmt.Printf("on write\n")
}

func (e *Echo) OnLoop(t time.Time, d time.Duration) {
	//fmt.Printf("on loop with delta:%v\n", d)
	if t.Sub(e.prev) > 3*time.Second {
		fmt.Printf("e")
	}
	e.prev = t
}
