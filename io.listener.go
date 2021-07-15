package elio

import (
	"fmt"
	"net"
	"os"
	"sync"
)

// Listener listener
type Listener struct {
	sync.RWMutex
	listenFd int // listen fd
	listenF  *os.File
	listener net.Listener
	ios      IOs
	reuse    bool
}

// newListener new listener
func newListener(io string, reuse bool) *Listener {
	l := new(Listener)
	if nil != l {
		if IoPoll == GetCurrentIO(io) {
			l.ios = IoPoll
		} else {
			l.ios = IoDefault
		}

		l.reuse = reuse
	}
	return l
}

func (l *Listener) String() string {
	return fmt.Sprintf("Listener::%p", l)
}

// ToFile to file
func (l *Listener) ToFile() (err error) {
	if l.listenF, err = l.listener.(*net.TCPListener).File(); nil == err {
		l.listenFd = int(l.listenF.Fd())
	}

	return err
}

// ToFd to file descriptor
func (l *Listener) ToFd() int {
	return l.listenFd
}
