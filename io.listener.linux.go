// +build linux

package elio

import (
	"net"

	"github.com/libp2p/go-reuseport"
	"golang.org/x/sys/unix"
)

// Listen listen
func (l *Listener) Listen(network, address string) (err error) {
	l.Lock()
	defer l.Unlock()

	if nil == l.listener {
		if true == l.reuse {
			l.listener, err = reuseport.Listen(network, address)
		} else {
			l.listener, err = net.Listen(network, address)
		}
		if nil == err {
			if IoDefault != l.ios {
				if l.listenF, err = l.listener.(*net.TCPListener).File(); nil == err {
					l.listenFd = int(l.listenF.Fd())
				}
			}
		}
	}

	return err
}

// Close close
func (l *Listener) Close() {
	l.Lock()
	defer l.Unlock()

	if nil != l.listener {
		if IoDefault != l.ios {
			if nil != l.listenF {
				//fmt.Printf("poll:%s close listen.f:%p ...\n", m.String(), m.listenF)
				l.listenF.Close()
				l.listenF = nil
			}

			AppDebug().Str(LogObject, l.String()).
				Msgf("listener:'%s' close listen.fd:%d", l.String(), l.listenFd)
			unix.Close(l.listenFd)
		}

		//fmt.Printf("service:%s close listener\n", m.String())
		if e := l.listener.Close(); nil != e {
			AppDebug().Str(LogObject, l.String()).
				Msgf("failed to close listener with error:'%v'", e)
		}

		l.listener = nil
	}
}
