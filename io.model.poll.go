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
	ioCore     *IoCore
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

func (m *ioPoll) GetIoCore() *IoCore {
	return m.ioCore
}

func (m *ioPoll) SetIoCore(c *IoCore) {
	m.ioCore = c
	c.SetIo(m)
}

// Listen listen
func (m *ioPoll) Listen(addr string) bool {
	m.GetIoCore().Host.Wg.Add(1)
	defer func() {
		m.GetIoCore().Host.Wg.Done()
	}()

	var err error
	if err = m.GetIoCore().Listener.Listen("tcp", addr); nil == err {
		//AppInfo().Str(LogObject, m.String()).Msgf("succeed to listen with url:%s", addr)
		//fmt.Printf("succeed to listen with url '%s'\n", m.GetIoCore().Addr.String())

		if 1 == m.GetIoCore().InCount.Add(1) {
			m.GetIoCore().InAddr.Store(addr)
		}

		//m.ctxAccept, m.cancelAccept = context.WithCancel(context.Background())
		//m.ctxIo, m.cancelIo = context.WithCancel(context.Background())

		go m.loopAccept()

		return true
	}

	AppError().Str(LogObject, m.String()).
		Err(err).Msgf("listen service url:'%s' failed", addr)

	return false
}

// Run run
func (m *ioPoll) Run() bool {
	m.GetIoCore().Host.Wg.Add(1)

	go m.loopIo()

	return true
}

// Shut shut listen
func (m *ioPoll) Shut() {
	//fmt.Printf("poll:%s shut\n", m.String())
	addr := m.GetIoCore().InAddr.Load()

	AppDebug().Str(LogObject, m.String()).
		Msgf("shut service url:'%s' poll.accept:%s", addr, m.pollAccept.String())

	if nil != m.GetIoCore().Listener.listener {
		m.GetIoCore().Listener.Close()
	}

	m.pollAccept.End()
}

// End end
func (m *ioPoll) End() {
	m.CloseAll()
	m.Shut()
	//m.pollAccept.End()
	m.pollIo.End()

	//m.GetIoCore().Host.Wg.Done()
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

			//fmt.Printf("session:%s on.read with %d bytes\n", n.String(), receipt)
			_ = m.GetIoCore().Service.OnRead(n, n.buffer[:r])
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
				Err(err).Msgf("failed to write with fd:%v written:%d/%d count.outqueue:%d", n.fd, written, l)
		}
	}

	if 0 < written {
		AppTrace().Str(LogObject, m.String()).Str(LogSession, n.String()).
			Msgf("succeed to write with fd:%v out:%d/%d", n.fd, written, l)

		m.GetIoCore().Service.OnWrite(n, out[:written])
	}

	if nil != err {
		if (unix.EAGAIN != err) && (unix.EWOULDBLOCK != err) {
			n.ioCore.Service.OnError(n, err)
			n.ioCore.io.Close(n)
		}
	}

	return written, err
}

// PostWrite post write
func (m *ioPoll) PostWrite(n *Session, out []byte) (written int, err error) {
	if nil != out {
		// co := n.CountOutQueue()
		// if 0 == co {
		// 	written, err = m.Write(n, out)
		// 	if written == len(out) {
		// 		return written, err
		// 	}

		// 	if (unix.EAGAIN != err) && (unix.EWOULDBLOCK != err) {
		// 		AppError().Str(LogObject, m.String()).Str(LogSession, n.String()).
		// 			Err(err).Msgf("failed to write %d", written)
		// 		return written, err
		// 	}

		// 	AppDebug().Str(LogObject, m.String()).Str(LogSession, n.String()).
		// 		Err(err).Msgf("unable to write %d", written)
		// }

		n.IncRef()
		b := GetByteBuffer()
		b.Write(out[written:])
		n.outQueue.Append(b)

		//written = len(out)
	}

	//fmt.Printf("poll:%s write %d B to fd:%d\n", m.String(), sent, n.fd)

	if 0 < n.CountOutQueue() {
		AppTrace().Str(LogObject, m.String()).Str(LogSession, n.String()).
			Msgf("post.write trigger out:%d count.outqueue:%d", len(out), n.CountOutQueue())

		err = m.pollIo.ControlMod(n.fd, unix.EPOLLET|unix.EPOLLONESHOT|unix.EPOLLIN|unix.EPOLLOUT)
		if nil != err {
			AppError().Str(LogObject, m.String()).Str(LogSession, n.String()).
				Err(err).Msg("failed to write")
			//fmt.Printf("poll:%s failed to write to fd:%d with error:%v\n", m.String(), n.fd, err)
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

	m.GetIoCore().sessionCmap.IterCb(c)
}

const waitTimeoutMsec = -1 //10
const pollAcceptCount = 1  //4

// loopAccept loop
func (m *ioPoll) loopAccept() {
	m.GetIoCore().Host.Wg.Add(1)

	addr := m.GetIoCore().InAddr.Load()

	AppInfo().Str(LogObject, m.String()).
		Msgf("succeed to listen with url:%s", addr)
	// if zerolog.InfoLevel < zerolog.GlobalLevel() {
	// 	fmt.Printf("succeed to listen with url:'%s'", addr.String())
	// }

	var err error

	defer func() {
		if nil != m.pollAccept {
			m.pollAccept.End()
			m.pollAccept = nil
		}
		m.GetIoCore().Host.Wg.Done()
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
		if err = m.pollAccept.ControlAdd(m.GetIoCore().Listener.ToFd(), unix.EPOLLIN|unix.EPOLLEXCLUSIVE); nil != err {
			//fmt.Println("poll add.r failed with error:", err.Error())
			AppError().Str(LogObject, m.String()).Err(err).
				Msgf("failed to ready listen with poll.accept:%s", m.pollAccept.String())
		}

		//fmt.Printf("poll:%s ready with listen.fd:%d\n", p.String(), p.listenFd)
		//AppDebug().Str(LogObject, m.String()).
		//	Msgf("succeed to ready with poll.accept:%s listen.fd:%d",
		//		m.pollAccept.String(), m.GetIoCore().Listener.ToFd())

		if err = unix.SetNonblock(m.GetIoCore().Listener.ToFd(), true); nil != err {
			//fmt.Println("poll set.nonblock failed with error:", err.Error())
			AppError().Str(LogObject, m.String()).Err(err).
				Msgf("failed to set.nonblock poll.accept:%s", m.pollAccept.String())
		}

	} else {
		AppError().Str(LogObject, m.String()).
			Msgf("failed to begin poll:%s", m.pollAccept.String())
		return
	}

	//fmt.Println("succeed to poll.ready")
	m.GetIoCore().Service.OnListen(m.GetIoCore())

	AppInfo().Str(LogObject, m.String()).
		Msgf("poll.accept.wait with %d events", m.GetIoCore().Config.InWaitCount)

	var wg sync.WaitGroup
	wg.Add(pollAcceptCount)

	for i := 0; i < pollAcceptCount; i++ {
		go func(w *sync.WaitGroup) {
			defer w.Done()

			var events []unix.EpollEvent
			events = make([]unix.EpollEvent, m.GetIoCore().Config.InWaitCount)

			for {
				var waits int
				if waits, err = m.pollAccept.Wait(events, waitTimeoutMsec); nil != err {
					AppError().Str(LogObject, m.String()).Err(err).
						Msgf("poll.accept:%s failed to wait", m.pollAccept.String())
					//fmt.Printf("poll:%s wait exit...\n", m.poll.String())
					break
				}

				failed := 0

				for i := 0; i < waits; i++ {
					//AppDebug().Str(LogObject, m.String()).
					//Msgf("poll:%s wait event:%v with fd:%d", m.poll.String(), events[i].Events, events[i].Fd)
					//fmt.Printf("poll:%s wait event:%v with fd:%d\n", m.poll.String(), events[i].Events, events[i].Fd)

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
	m.GetIoCore().Service.OnShut(m.GetIoCore())

	AppInfo().Str(LogObject, m.String()).Msgf("service:%s close poll.accpet", m.String())
}

const pollIoCount = 1

// loopIo loop io
func (m *ioPoll) loopIo() {
	//addr := m.GetIoCore().InAddr.Load().(*net.TCPAddr)

	// AppInfo().Str(LogObject, m.String()).
	// 	Msgf("succeed to listen with url:'%s'", addr.String())
	// if zerolog.InfoLevel < zerolog.GlobalLevel() {
	// 	fmt.Printf("succeed to listen with url:'%s'", addr.String())
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
		Msgf("poll.io.wait with %d events", m.GetIoCore().Config.InWaitCount)

	var events []unix.EpollEvent
	events = make([]unix.EpollEvent, m.GetIoCore().Config.InWaitCount)

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
					//fmt.Printf("poll:%s wait exit...\n", m.poll.String())
					break
				}

				//fmt.Printf("awake wait... %d events\n", waits)
				//fmt.Printf("\n\n\npoll.wait: waits %d events\n\n\n", waits)

				for i := 0; i < waits; i++ {
					//AppDebug().Str(LogObject, m.String()).
					//Msgf("poll:%s wait event:%v with fd:%d", m.poll.String(), events[i].Events, events[i].Fd)
					//fmt.Printf("poll:%s wait event:%v with fd:%d\n", m.poll.String(), events[i].Events, events[i].Fd)

					fd := int(events[i].Fd)
					if fd == m.pollIo.wfd {
						AppDebug().Str(LogObject, m.String()).
							Msgf("poll:%s wake up event:%v with fd:%d", m.pollIo.String(), events[i].Events, events[i].Fd)

						_, _ = unix.Read(m.pollIo.wfd, m.pollIo.wfdBuf)

						//jobs := m.pollIo.jobs.Fetch()
						//for _, job := range jobs {
						//	job.(Job).Work()
						//}

					} else {
						var n *Session

						if result, ok := m.GetIoCore().sessionCmap.Get(fmt.Sprintf("0x%08x", fd)); true == ok {
							n = result.(*Session)

						} else {
							AppError().Str(LogObject, m.String()).
								Msgf("failed to get session by fd:%d", fd)
							continue
						}

						if unix.EPOLLOUT == events[i].Events&unix.EPOLLOUT {
							AppTrace().Str(LogObject, m.String()).
								Msgf("poll:%s awake fd:%d with events.epollout:%X", m.pollIo.String(), events[i].Fd, events[i].Events)
							//fmt.Printf("events.epollout:%X with fd:%d", events[i].Events, events[i].Fd)

							m.handleWrite(n)
						}
						if unix.EPOLLIN == events[i].Events&unix.EPOLLIN {
							AppTrace().Str(LogObject, m.String()).
								Msgf("poll:%s awake fd:%d with events.epollin:%X", m.pollIo.String(), events[i].Fd, events[i].Events)
							//fmt.Printf("events.epollin:%X with fd:%d", events[i].Events, events[i].Fd)

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
	//fmt.Printf("awake accepting session:%d with EPOLLIN\n", fd)
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

// 	//fmt.Printf("get error events by fd:%d\n", n.fd)

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
		m.GetIoCore().Service.OnError(n, err)
	}

	m.GetIoCore().Service.OnClose(n, err)

	m.pollIo.ControlDel(n.fd, 0)
	m.Close(n)

	m.GetIoCore().sessionCmap.Remove(fmt.Sprintf("0x%08x", n.fd))

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
				//fmt.Printf("poll accept fd:%d failed with error(unix.EAGAIN):%v\n", fd, err)
				//AppDebug().Str(LogObject, m.String()).
				//	Msgf("accepting would block with fd:%d error(EAGAIN):%v", fd, err)
				break
			}
			//fmt.Printf("poll accept fd:%d failed with error:'%v'\n", fd, err)
			AppError().Str(LogObject, m.String()).
				Err(err).Msgf("accept to fd:%d failed", fd)
			break
		}

		//if false == m.GetIoCore().sessionCmap.Has(fmt.Sprintf("0x%08x", nfd)) {
		if err := unix.SetNonblock(nfd, true); nil != err {
			//fmt.Printf("accepting setnonblock failed with error:'%v'\n", err)
			AppError().Str(LogObject, m.String()).
				Err(err).Msgf("accept set.nonblock to fd:%d failed", nfd)
			break
		}

		if true == m.GetIoCore().Config.InNoDelay {
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

		if 0 != m.GetIoCore().Config.InRcvBuff {
			if err = unix.SetsockoptInt(nfd, syscall.SOL_SOCKET, syscall.SO_RCVBUF, m.GetIoCore().Config.InRcvBuff*1024); err != nil {
				AppError().Str(LogObject, m.String()).
					Err(err).Msgf("accept set.rcvbuf:%d(KB) to fd:%d failed", m.GetIoCore().Config.InRcvBuff, nfd)
				break
			} else {
				//AppDebug().Str(LogObject, m.String()).
				//	Msgf("succeed to accepting set.nodelay with fd:%d true'", nfd)
			}
		}

		if 0 != m.GetIoCore().Config.InSndBuff {
			if err = unix.SetsockoptInt(nfd, syscall.SOL_SOCKET, syscall.SO_SNDBUF, m.GetIoCore().Config.InSndBuff*1024); err != nil {
				AppError().Str(LogObject, m.String()).
					Err(err).Msgf("accept set.sndbuf:%d(KB) to fd:%d failed", m.GetIoCore().Config.InSndBuff, nfd)
				break
			} else {
				//AppDebug().Str(LogObject, m.String()).
				//	Msgf("succeed to accepting set.nodelay with fd:%d true'", nfd)
			}
		}

		/*//
		if 0 < m.GetIoCore().Config.InRecvTimeo {
			tv := unix.NsecToTimeval(int64(m.GetIoCore().Config.InRecvTimeo) * int64(time.Millisecond))
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

		n := NewSession(nfd, nil, m.GetIoCore())

		m.GetIoCore().Service.OnOpen(n)

		m.GetIoCore().sessionCmap.Set(fmt.Sprintf("0x%08x", n.fd), n)
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
// 			if nil != m.GetIoCore().Config.Events.OnRead {
// 				AppTrace().Str(LogObject, m.String()).Str(LogSession, n.String()).
// 					Msgf("succeed to read with fd:%v in:%d", n.fd, r)

// 				//fmt.Printf("session:%s on.read with %d bytes\n", n.String(), receipt)
// 				_ = m.GetIoCore().Config.Events.OnRead(n, n.buffer[:r])
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
