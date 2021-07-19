// +build linux

package elio

import (
	"fmt"
	"io"
	"sync"
	"syscall"

	"golang.org/x/sys/unix"
)

// ioPoll service implementation
type ioPoll struct {
	io         *Io
	pollAccept *Poll
	pollIo     *Poll
}

// NewIoPoll new io poll
func NewIoPoll() *ioPoll {
	io := new(ioPoll)
	return io
}

func (m *ioPoll) String() string {
	return fmt.Sprintf("ioPoll::%p", m)
}

func (m *ioPoll) GetIo() *Io {
	return m.io
}

func (m *ioPoll) SetIo(io *Io) {
	m.io = io
	io.SetIoModel(m)
}

// Listen listen
func (m *ioPoll) Listen(addr string) bool {
	m.GetIo().Host.Wg.Add(1)
	defer func() {
		m.GetIo().Host.Wg.Done()
	}()

	var err error
	if err = m.GetIo().Listener.Listen("tcp", addr); nil == err {
		//AppInfo().Str(LogObject, m.String()).Msgf("succeed to listen with url:%s", addr)

		if 1 == m.GetIo().InCount.Add(1) {
			m.GetIo().InAddr.Store(addr)
		}

		go m.loopAccept()

		return true
	}

	AppError().Str(LogObject, m.String()).
		Err(err).Msgf("listen service url:'%s' failed", addr)

	return false
}

// Run run
func (m *ioPoll) Run() bool {
	m.GetIo().Host.Wg.Add(1)

	go m.loopIo()

	return true
}

// Shut shut listen
func (m *ioPoll) Shut() {
	addr := m.GetIo().InAddr.Load()

	AppDebug().Str(LogObject, m.String()).
		Msgf("shut service url:'%s' poll.accept:%s", addr, m.pollAccept.String())

	if nil != m.GetIo().Listener.listener {
		m.GetIo().Listener.Close()
	}

	m.pollAccept.End()
}

// End end
func (m *ioPoll) End() {
	m.CloseAll()
	m.Shut()

	m.pollIo.End()
}

// Read read from session
func (m *ioPoll) Read(n *Session, in []byte) (receipt int, err error) {
	for {
		var r int
		r, err = unix.Read(n.fd, in)
		if 0 < r {
			receipt += r

			AppTrace().Str(LogObject, m.String()).Str(LogSession, n.String()).
				Msgf("succeed to read with fd:%v in:%d/%d", n.fd, r, receipt)

			_ = m.GetIo().Service.OnRead(n, n.buffer[:r])
		}

		if 0 == r || nil != err {
			if nil == err {
				err = io.EOF
			}

			if (unix.EAGAIN == err) || (unix.EWOULDBLOCK == err) {
				//AppTrace().Str(LogObject, m.String()).Str(LogSession, n.String()).
				//	Err(err).Msgf("unable to read fd:%v in:%d", n.fd, r)

			} else {
				if io.EOF == err {
					AppTrace().Str(LogObject, m.String()).Str(LogSession, n.String()).
						Err(err).Msgf("closed to read with fd:%d receipt:%d", n.fd, r)

				} else {
					AppError().Str(LogObject, m.String()).Str(LogSession, n.String()).
						Err(err).Msgf("failed to read with fd:%d failed receipt:%d", n.fd, r)
				}
			}

			break
		}
	}

	return receipt, err
}

// Write write
func (m *ioPoll) Write(n *Session, out []byte) (written int, err error) {
	l := len(out)
	written, err = n.sysWriteAll(out)
	if nil != err {
		if (unix.EAGAIN == err) || (unix.EWOULDBLOCK == err) {
			//AppTrace().Str(LogObject, m.String()).Str(LogSession, n.String()).
			//	Err(err).Msgf("unable to write fd:%v out:%d/%d", n.fd, written, l)

		} else {
			AppError().Str(LogObject, m.String()).Str(LogSession, n.String()).
				Err(err).Msgf("failed to write with fd:%v written:%d/%d", n.fd, written, l)
		}
	}

	if 0 < written {
		AppTrace().Str(LogObject, m.String()).Str(LogSession, n.String()).
			Msgf("succeed to write with fd:%v out:%d/%d", n.fd, written, l)

		m.GetIo().Service.OnWrite(n, out[:written])
	}

	if nil != err {
		if (unix.EAGAIN != err) && (unix.EWOULDBLOCK != err) {
			n.io.Service.OnError(n, err)
			n.io.ioModel.Close(n)
		}
	}

	return written, err
}

// PostWrite post write
func (m *ioPoll) PostWrite(n *Session, out []byte) (written int, err error) {
	if nil != out {
		n.IncRef()
		b := GetByteBuffer()
		b.Write(out[written:])
		n.outQueue.Append(b)
	}

	if 0 < n.CountOutQueue() {
		AppTrace().Str(LogObject, m.String()).Str(LogSession, n.String()).
			Msgf("post.write trigger out:%d count.outqueue:%d", len(out), n.CountOutQueue())

		err = m.pollIo.ControlMod(n.fd, unix.EPOLLET|unix.EPOLLONESHOT|unix.EPOLLIN|unix.EPOLLOUT)
		if nil != err {
			AppError().Str(LogObject, m.String()).Str(LogSession, n.String()).
				Err(err).Msg("failed to write")

			m.Close(n)
		}
	}

	return written, err
}

func (m *ioPoll) Trigger(j interface{}) (err error) {
	return m.pollIo.Trigger(j)
}

func (m *ioPoll) Close(n *Session) (err error) {
	//_ = n.writeQueue.Fetch()

	//AppDebug().Str(LogObject, m.String()).Str(LogSession, n.String()).Msgf("session.event close")
	err = unix.Close(n.fd)

	return err
}

func (m *ioPoll) Shutdown(n *Session, how int) error {
	return unix.Shutdown(n.fd, how)
}

func (m *ioPoll) CloseAll() {
	c := func(key string, v interface{}) {
		s := v.(*Session)
		m.Close(s)
	}

	m.GetIo().sessionCmap.IterCb(c)
}

const waitTimeoutMsec = -1 //10
const pollAcceptCount = 1  //4

// loopAccept loop
func (m *ioPoll) loopAccept() {
	m.GetIo().Host.Wg.Add(1)

	addr := m.GetIo().InAddr.Load()

	AppInfo().Str(LogObject, m.String()).
		Msgf("succeed to listen with url:%s", addr)
	// if zerolog.InfoLevel < zerolog.GlobalLevel() {
	//	AppDebug().Str(LogObject, m.String()).
	//		Msgf("succeed to listen with url:'%s'", addr.String())
	// }

	var err error

	defer func() {
		if nil != m.pollAccept {
			m.pollAccept.End()
			m.pollAccept = nil
		}
		m.GetIo().Host.Wg.Done()
	}()

	if m.pollAccept = NewPoll(); nil == m.pollAccept {
		AppError().Str(LogObject, m.String()).
			Msgf("failed to create poll.accept:%s", m.pollAccept.String())
		return
	}

	AppDebug().Str(LogObject, m.String()).
		Msgf("create poll.accept:%s", m.pollAccept.String())

	if err = m.pollAccept.Begin(); nil == err {
		// accepting is no need to be oneshot
		if err = m.pollAccept.ControlAdd(m.GetIo().Listener.ToFd(), unix.EPOLLIN|unix.EPOLLEXCLUSIVE); nil != err {
			AppError().Str(LogObject, m.String()).Err(err).
				Msgf("failed to ready listen with poll.accept:%s", m.pollAccept.String())
		}

		//AppDebug().Str(LogObject, m.String()).
		//	Msgf("succeed to ready with poll.accept:%s listen.fd:%d",
		//		m.pollAccept.String(), m.GetIo().Listener.ToFd())

		if err = unix.SetNonblock(m.GetIo().Listener.ToFd(), true); nil != err {
			AppError().Str(LogObject, m.String()).Err(err).
				Msgf("failed to set.nonblock poll.accept:%s", m.pollAccept.String())
		}

	} else {
		AppError().Str(LogObject, m.String()).
			Msgf("failed to begin poll:%s", m.pollAccept.String())
		return
	}

	//m.GetIo().Service.OnListen(m.GetIo())

	AppInfo().Str(LogObject, m.String()).
		Msgf("poll.accept.wait with %d events", m.GetIo().Config.InWaitCount)

	var wg sync.WaitGroup
	wg.Add(pollAcceptCount)

	for i := 0; i < pollAcceptCount; i++ {
		go func(w *sync.WaitGroup) {
			defer w.Done()

			var events []unix.EpollEvent
			events = make([]unix.EpollEvent, m.GetIo().Config.InWaitCount)

			for {
				var waits int
				if waits, err = m.pollAccept.Wait(events, waitTimeoutMsec); nil != err {
					AppError().Str(LogObject, m.String()).Err(err).
						Msgf("poll.accept:%s failed to wait", m.pollAccept.String())
					break
				}

				failed := 0

				for i := 0; i < waits; i++ {
					//AppDebug().Str(LogObject, m.String()).
					//	Msgf("poll:%s wait event:%v with fd:%d", m.poll.String(), events[i].Events, events[i].Fd)

					fd := int(events[i].Fd)

					if unix.EPOLLIN == events[i].Events&unix.EPOLLIN {
						if false == m.handleAccept(fd) {
							failed++
						}
					}
				}

				if 0 < failed {
					break
				}
			}

		}(&wg)
	}

	wg.Wait()

	AppInfo().Str(LogObject, m.String()).
		Msgf("service:%s poll:%s on.shut", m.String(), m.pollAccept.String())
	//m.GetIo().Service.OnShut(m.GetIo())

	AppInfo().Str(LogObject, m.String()).Msgf("service:%s close poll.accpet", m.String())
}

const pollIoCount = 1

// loopIo loop io
func (m *ioPoll) loopIo() {
	//addr := m.GetIo().InAddr.Load().(*net.TCPAddr)

	// AppInfo().Str(LogObject, m.String()).
	// 	Msgf("succeed to listen with url:'%s'", addr.String())
	// if zerolog.InfoLevel < zerolog.GlobalLevel() {
	//	 AppInfo().Str(LogObject, m.String()).
	//	 	Msgf("succeed to listen with url:'%s'", addr.String())
	// }

	var err error

	defer func() {
		if nil != m.pollIo {
			m.CloseAll()
			m.pollIo.End()
			m.pollIo = nil
		}
	}()

	// Start epoll io
	if m.pollIo = NewPoll(); nil == m.pollIo {
		AppError().Str(LogObject, m.String()).Msg("failed to create poll")
		return
	}

	AppDebug().Str(LogObject, m.String()).
		Msgf("create poll.io:%s", m.pollIo.String())

	if err = m.pollIo.Begin(); nil == err {
	} else {
		AppError().Str(LogObject, m.String()).
			Msgf("failed to begin poll:%s", m.pollIo.String())
		return
	}

	AppInfo().Str(LogObject, m.String()).
		Msgf("poll.io.wait with %d events", m.GetIo().Config.InWaitCount)

	var events []unix.EpollEvent
	events = make([]unix.EpollEvent, m.GetIo().Config.InWaitCount)

	var wg sync.WaitGroup
	wg.Add(pollIoCount)

	for i := 0; i < pollIoCount; i++ {
		go func(w *sync.WaitGroup) {
			defer w.Done()

			for {
				var waits int
				if waits, err = m.pollIo.Wait(events, waitTimeoutMsec); nil != err {
					AppError().Str(LogObject, m.String()).Err(err).
						Msgf("poll.io:%s failed to wait", m.pollIo.String())
					break
				}

				for i := 0; i < waits; i++ {
					//AppDebug().Str(LogObject, m.String()).
					//Msgf("poll:%s wait event:%v with fd:%d", m.poll.String(), events[i].Events, events[i].Fd)

					fd := int(events[i].Fd)
					if fd == m.pollIo.wfd {
						AppDebug().Str(LogObject, m.String()).
							Msgf("poll:%s wake up event:%v with fd:%d", m.pollIo.String(), events[i].Events, events[i].Fd)

						_, _ = unix.Read(m.pollIo.wfd, m.pollIo.wfdBuf)

					} else {
						var n *Session

						if result, ok := m.GetIo().sessionCmap.Get(fmt.Sprintf("0x%08x", fd)); true == ok {
							n = result.(*Session)

						} else {
							AppError().Str(LogObject, m.String()).
								Msgf("failed to get session by fd:%d", fd)
							continue
						}

						if unix.EPOLLOUT == events[i].Events&unix.EPOLLOUT {
							AppTrace().Str(LogObject, m.String()).
								Msgf("poll:%s awake fd:%d with events.epollout:%X", m.pollIo.String(), events[i].Fd, events[i].Events)

							m.handleWrite(n)
						}
						if unix.EPOLLIN == events[i].Events&unix.EPOLLIN {
							AppTrace().Str(LogObject, m.String()).
								Msgf("poll:%s awake fd:%d with events.epollin:%X", m.pollIo.String(), events[i].Fd, events[i].Events)

							m.handleRead(n)
						}
						//if unix.EPOLLERR == events[i].Events&unix.EPOLLERR {
						//	m.handleError(n)
						//}
					}
				}
			}
		}(&wg)
	}

	wg.Wait()

	AppInfo().Str(LogObject, m.String()).Msgf("service:%s close poll.io", m.String())
}

// canRepost can repost?
func canRepost(err error) bool {
	if (nil == err) || (unix.EAGAIN == err) || (unix.EWOULDBLOCK == err) {
		return true
	}

	return false
}

func (m *ioPoll) handleAccept(fd int) bool {
	err := m.accept(fd)
	if true == canRepost(err) {
		//AppDebug().Str(LogObject, m.String()).Str(LogSession, n.String()).Err(err).
		//Msgf("session.event control.mod receipt:%d", receipt)

		// control mod is not necessary for epoll exclusive flags
		//if err = m.pollAccept.ControlMod(fd, unix.EPOLLET|unix.EPOLLONESHOT|unix.EPOLLIN); nil != err {
		//	AppError().Str(LogObject, m.String()).
		//		Msgf("failed to poll.mod with listen.fd:%d error:'%v'", fd, err)
		//	return false
		//}

	} else {
		AppError().Str(LogObject, m.String()).
			Msgf("failed to accept by listen.fd:%d with error:'%v'", fd, err)

	}

	return true
}

func (m *ioPoll) handleRead(n *Session) {
	_, err := m.Read(n, n.buffer)
	if true == canRepost(err) {
		//AppDebug().Str(LogObject, m.String()).Str(LogSession, n.String()).
		//	Err(err).Msgf("handle.read canRepost count.outqueue:%d", co)

		if 0 == n.CountOutQueue() {
			if err = m.pollIo.ControlMod(n.fd, unix.EPOLLET|unix.EPOLLONESHOT|unix.EPOLLIN); nil != err {
				AppError().Str(LogObject, m.String()).Str(LogSession, n.String()).
					Err(err).Msgf("failed to poll.mod.r to fd:%d", n.fd)
			}

		} else {
			if err = m.pollIo.ControlMod(n.fd, unix.EPOLLET|unix.EPOLLONESHOT|unix.EPOLLIN|unix.EPOLLOUT); nil != err {
				AppError().Str(LogObject, m.String()).Str(LogSession, n.String()).
					Err(err).Msgf("failed to poll.mod.rw to fd:%d", n.fd)
			}
		}

	} else {
		AppTrace().Str(LogObject, m.String()).Str(LogSession, n.String()).Err(err).
			Msgf("handle.read release session")

		// if io.EOF == err {
		// 	AppDebug().Str(LogObject, m.String()).Str(LogSession, n.String()).
		// 		Err(err).Msgf("succeed to close with fd:%d", n.fd)
		// } else {
		// 	AppError().Str(LogObject, m.String()).Str(LogSession, n.String()).
		// 		Err(err).Msgf("failed to read with fd:%d", n.fd)
		// }

		m.releaseSession(n, err)
	}
}

func (m *ioPoll) handleWrite(n *Session) {
	co := n.CountOutQueue()
	if 0 < co {
		AppTrace().Str(LogObject, m.String()).Str(LogSession, n.String()).
			Msgf("handle.write awake with fd:%d count.outqueue:%d", n.fd, co)

		written, err := m.write(n)
		if true == canRepost(err) {
			if nil != err {
				AppDebug().Str(LogObject, m.String()).Str(LogSession, n.String()).
					Err(err).Msgf("handle.write can repost though error with written:%d count.outqueue:%d", written, n.CountOutQueue())
			}

			if 0 == n.CountOutQueue() {
				if err = m.pollIo.ControlMod(n.fd, unix.EPOLLET|unix.EPOLLONESHOT|unix.EPOLLIN); nil != err {
					AppError().Str(LogObject, m.String()).Str(LogSession, n.String()).
						Err(err).Msgf("failed to poll.mod.r to fd:%d", n.fd)
				}

			} else {
				if err = m.pollIo.ControlMod(n.fd, unix.EPOLLET|unix.EPOLLONESHOT|unix.EPOLLIN|unix.EPOLLOUT); nil != err {
					AppError().Str(LogObject, m.String()).Str(LogSession, n.String()).
						Err(err).Msgf("failed to poll.mod.rw to fd:%d", n.fd)
				}
			}

		} else {
			AppError().Str(LogObject, m.String()).Str(LogSession, n.String()).
				Err(err).Msgf("failed to write with fd:%d written:%d", n.fd, written)

			//m.releaseSession(n, err)
		}
	}
}

// func (m *ioPoll) handleError(n *Session) {
// 	//n, e := getsockoptIntFunc(fd, unix.SOL_SOCKET, unix.SO_ERROR)
// 	//if e != nil {
// 	//	err = os.NewSyscallError("getsockopt", e)
// 	//}

// 	t := fmt.Sprintf("run.error fd:%d with EPOLLERR", n.fd)
// 	AppDebug().Str(LogObject, m.String()).
// 		Str(LogSession, n.String()).Msg(t)
// 	err := errors.New(t)

// 	m.releaseSession(n, err)
// }

// func (m *ioPoll) handleHup(n *Session) {
// 	//eno, enoOk := err.(syscall.Errno)
// 	//t := fmt.Sprintf("run.hup fd:%d with EPOLLHUP errno:%v, ok:%v", fd, eno, enoOk)
// 	err := fmt.Errorf("handle.hup EPOLLHUP with fd:%d", n.fd)

// 	AppError().Str(LogObject, m.String()).Str(LogSession, n.String()).
// 		Err(err).Msgf("handle.hup fd:%d", n.fd)

// 	m.releaseSession(n, err)
// }

// func (m *ioPoll) handleRdHup(n *Session) {
// 	err := fmt.Errorf("handle.rdhup EPOLLRDHUP with fd:%d", n.fd)

// 	AppError().Str(LogObject, m.String()).Str(LogSession, n.String()).Str(LogSession, n.String()).
// 		Err(err).Msgf("handle.rdhup fd:%d", n.fd)

// 	m.releaseSession(n, err)
// }

func (m *ioPoll) releaseSession(n *Session, err error) {
	if io.EOF != err {
		m.GetIo().Service.OnError(n, err)
	}

	m.GetIo().Service.OnClose(n, err)

	m.pollIo.ControlDel(n.fd, 0)
	m.Close(n)

	m.GetIo().sessionCmap.Remove(fmt.Sprintf("0x%08x", n.fd))

	n.DecRef()
}

// accept accept socket
func (m *ioPoll) accept(fd int) (err error) {
	for {
		var nfd int
		var sa unix.Sockaddr // remote socket address
		nfd, sa, err = unix.Accept(fd)
		if nil != err {
			if (unix.EAGAIN == err) || (unix.EWOULDBLOCK == err) {
				//AppDebug().Str(LogObject, m.String()).
				//	Msgf("accepting would block with fd:%d error(EAGAIN):%v", fd, err)
				break
			}

			AppError().Str(LogObject, m.String()).
				Err(err).Msgf("accept to fd:%d failed", fd)
			break
		}

		//if false == m.GetIo().sessionCmap.Has(fmt.Sprintf("0x%08x", nfd)) {
		if err := unix.SetNonblock(nfd, true); nil != err {
			AppError().Str(LogObject, m.String()).
				Err(err).Msgf("accept set.nonblock to fd:%d failed", nfd)
			break
		}

		if true == m.GetIo().Config.InNoDelay {
			// This should disable Nagle's algorithm in all accepted sockets by default.
			// Users may enable it with net.TCPConn.SetNoDelay(false).
			if err = unix.SetsockoptInt(nfd, syscall.IPPROTO_TCP, syscall.TCP_NODELAY, 1); err != nil {
				AppError().Str(LogObject, m.String()).
					Err(err).Msgf("accept set.nodelay to fd:%d failed", nfd)
				break
			} else {
				//AppDebug().Str(LogObject, m.String()).
				//	Msgf("succeed to accepting set.nodelay with fd:%d true'", nfd)
			}
		}

		if 0 != m.GetIo().Config.InRcvBuff {
			if err = unix.SetsockoptInt(nfd, syscall.SOL_SOCKET, syscall.SO_RCVBUF, m.GetIo().Config.InRcvBuff*1024); err != nil {
				AppError().Str(LogObject, m.String()).
					Err(err).Msgf("accept set.rcvbuf:%d(KB) to fd:%d failed", m.GetIo().Config.InRcvBuff, nfd)
				break
			} else {
				//AppDebug().Str(LogObject, m.String()).
				//	Msgf("succeed to accepting set.nodelay with fd:%d true'", nfd)
			}
		}

		if 0 != m.GetIo().Config.InSndBuff {
			if err = unix.SetsockoptInt(nfd, syscall.SOL_SOCKET, syscall.SO_SNDBUF, m.GetIo().Config.InSndBuff*1024); err != nil {
				AppError().Str(LogObject, m.String()).
					Err(err).Msgf("accept set.sndbuf:%d(KB) to fd:%d failed", m.GetIo().Config.InSndBuff, nfd)
				break
			} else {
				//AppDebug().Str(LogObject, m.String()).
				//	Msgf("succeed to accepting set.nodelay with fd:%d true'", nfd)
			}
		}

		/*//
		if 0 < m.GetIo().Config.InRecvTimeo {
			tv := unix.NsecToTimeval(int64(m.GetIo().Config.InRecvTimeo) * int64(time.Millisecond))
			if err = unix.SetsockoptTimeval(nfd, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &tv); err != nil {
				AppDebug().Str(LogObject, m.String()).
					Err(err).Msgf("accept set.recvtimeout to fd:%d failed", nfd)
				break
			} else {
				AppDebug().Str(LogObject, m.String()).
					Msgf("succeed to accepting set.recvtimeout with fd:%d tv:%v", nfd, tv.Sec)
			}
		}
		//*/

		n := NewSession(nfd, nil, m.GetIo())

		m.GetIo().Service.OnOpen(n)

		m.GetIo().sessionCmap.Set(fmt.Sprintf("0x%08x", n.fd), n)
		n.sa = sa
		if err = m.pollIo.ControlAdd(n.fd, unix.EPOLLET|unix.EPOLLONESHOT|unix.EPOLLIN); nil != err {
			AppError().Str(LogObject, m.String()).Str(LogSession, n.String()).
				Err(err).Msgf("failed to accepting poll.add.r with fd:%d", n.fd)

		} else {
			//AppDebug().Str(LogObject, m.String()).Str(LogSession, n.String()).
			//	Msgf("fd:%d succeed to accept by listen.fd:%d", nfd, fd)

		}

		//} else {
		//    // skip accepting already logged in session
		//    AppWarn().Str(LogObject, m.String()).
		//    Msgf("skip to accepting fd:%d with already accepted", nfd)
		//}
	}

	return err
}

// // read read
// func (m *ioPoll) read(n *Session) (err error) {
// 	for {
// 		var r int
// 		r, err = m.Read(n, n.buffer)

// 		if 0 < r {
// 			if nil != m.GetIo().Config.Events.OnRead {
// 				AppTrace().Str(LogObject, m.String()).Str(LogSession, n.String()).
// 					Msgf("succeed to read with fd:%v in:%d", n.fd, r)

// 				_ = m.GetIo().Config.Events.OnRead(n, n.buffer[:r])
// 			}
// 		}

// 		if 0 == r || nil != err {
// 			if nil == err {
// 				err = io.EOF
// 			}

// 			if (unix.EAGAIN == err) || (unix.EWOULDBLOCK == err) {
// 				AppTrace().Str(LogObject, m.String()).Str(LogSession, n.String()).
// 					Err(err).Msgf("unable to read fd:%v in:%d", n.fd, r)

// 			} else {
// 				if io.EOF == err {
// 					AppTrace().Str(LogObject, m.String()).Str(LogSession, n.String()).
// 						Err(err).Msgf("closed to read with fd:%d in:%d", n.fd, r)

// 				} else {
// 					AppError().Str(LogObject, m.String()).Str(LogSession, n.String()).
// 						Err(err).Msgf("failed to read with fd:%d failed in:%d", n.fd, r)
// 				}
// 			}

// 			break
// 		}
// 	}

// 	return err
// }

// write write
func (m *ioPoll) write(n *Session) (written int, err error) {
	var out []byte

	outs := n.outQueue.Fetch()

	var ref int32
	ref = int32(len(outs))
	if 0 < ref {
		for _, o := range outs {
			out = append(out, o.(*ByteBuffer).Bytes()...)
			PutByteBuffer(o.(*ByteBuffer))
		}

		l := len(out)
		written, err = m.Write(n, out)
		if written < l {
			b := GetByteBuffer()
			b.Write(out[written:])
			n.outQueue.Prepend(b)
			n.SubRef(ref - 1)

		} else {
			n.SubRef(ref)
		}

	} else {
		AppDebug().Str(LogObject, m.String()).Str(LogSession, n.String()).
			Msgf("no writable out with fd:%d", n.fd)

	}

	return written, err
}
