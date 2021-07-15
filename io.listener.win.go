// +build windows

package elio

import "github.com/libp2p/go-reuseport"

// Listen listen
func (l *Listener) Listen(network, address string) (err error) {
	l.Lock()
	defer l.Unlock()

	if nil == l.listener {
		l.listener, err = reuseport.Listen(network, address)
	}

	return err
}

// Close close
func (l *Listener) Close() {
	l.Lock()
	defer l.Unlock()

	if nil != l.listener {
		//fmt.Printf("service:%s close listener\n", m.String())
		if e := l.listener.Close(); nil != e {
			AppDebug().Str(LogObject, l.String()).
				Msgf("failed to close listener with error:'%v'", e)
		}

		l.listener = nil
	}
}
