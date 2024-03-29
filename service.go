package elio

import (
	"context"
	"time"
)

const (
	OnReadInvalidLen = -1
)

// Service service
type Service interface {
	Name() string
	OnInit(ctx context.Context, cancel context.CancelFunc) error
	OnExit()
	//OnListen(i *IoCore)
	//OnShut(i *IoCore)
	OnOpen(s *Session) error
	OnClose(s *Session, err error)
	OnError(s *Session, err error)
	OnRead(s *Session, in []byte) int
	OnWrite(s *Session, out []byte)
	OnLoop(host *IoHost, t time.Time, d time.Duration)
}
