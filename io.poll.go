// This code has a lot of references to 'github.com/tidwall/evio' - cppis

// +build linux

package elio

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// PollEvent poll event
type PollEvent struct {
	unix.EpollEvent
	Session *Session
}

// Poll event poll object
type Poll struct {
	fd     int    // epoll fd
	wfd    int    // wake fd
	wfdBuf []byte // wfd buffer to read packet
	jobs   SafeSlice
}

// String object to string
func (p *Poll) String() string {
	return fmt.Sprintf("Poll::%p", p)
}

// NewPoll create poll
func NewPoll() (p *Poll) {
	p = new(Poll)
	return p
}

// Begin begin poll
func (p *Poll) Begin() (err error) {
	defer func() {
		if nil != err {
			if 0 != p.fd {
				unix.Close(p.fd)
			}
			if 0 != p.wfd {
				unix.Close(p.wfd)
			}
		}
	}()

	if p.fd, err = unix.EpollCreate1(0); nil != err {
		AppError().Err(err).Msg("failed to create epoll")
		return err
	}

	r0, _, e0 := unix.Syscall(unix.SYS_EVENTFD2, 0, 0, 0)
	if 0 != e0 {
		err = fmt.Errorf("failed to init fd:%v", r0)
	} else {
		p.wfd = int(r0)
		p.wfdBuf = make([]byte, 8)
		p.ControlAdd(p.wfd, unix.EPOLLIN|unix.EPOLLERR)
	}

	AppInfo().Str(LogObject, p.String()).Msg("begin epoll")

	return err
}

// End ...
func (p *Poll) End() {
	AppInfo().Str(LogObject, p.String()).Msgf("end epoll")

	var err error

	if 0 != p.wfd {
		if err = unix.Close(p.wfd); nil != err {
			AppError().Str(LogObject, p.String()).Err(err).
				Msgf("failed to close epoll")
		}

		p.wfd = 0
	}

	if 0 != p.fd {
		AppInfo().Str(LogObject, p.String()).Msgf("poll:%s closing", p.String())
		if err = unix.Close(p.fd); nil != err {
			AppError().Str(LogObject, p.String()).Err(err).
				Msgf("failed to close epoll")
		}

		p.fd = 0
	}
}

// WaitResult wait result
type WaitResult struct {
	evt unix.EpollEvent
	err error
}

// // WaitAndRun ...
// func (p *Poll) WaitAndRun(count int, iter func(e unix.EpollEvent) error) (errored []WaitResult, err error) {
// 	events := make([]unix.EpollEvent, count)
// 	var results []WaitResult

// 	for {
// 		var n int
// 		n, err = unix.EpollWait(p.fd, events, -1)
// 		if nil != err && unix.EINTR != err {
// 			AppError().Str(LogObject, p.String()).Err(err).
// 				Msg("failed to poll.wait.run")
// 			return results, err
// 		}

// 		AppDebug().Str(LogObject, p.String()).
// 			Msgf("poll.wait.run awake with %v events", n)

// 		for i := 0; i < n; i++ {
// 			if e := iter(events[i]); nil != e {
// 				// AppError().Str(LogObject, p.String()).
// 				//  Msgf("failed to poll.wait.run.iter with error:'%v'", e)
// 				results = append(results, WaitResult{evt: events[i], err: e})
// 			}
// 		}
// 	}
// }

// Wait ...
func (p *Poll) Wait(events []unix.EpollEvent, timeout int) (waits int, err error) {
	waits, err = unix.EpollWait(p.fd, events, timeout)
	if nil != err && unix.EINTR != err {
		//AppError().Str(LogObject, p.String()).Err(err).
		//	Msg("failed to poll.wait")
		return waits, err
	}

	//AppDebug().Str(LogObject, p.String()).
	//	Msgf("poll.wait awake with %v events", waits)

	return waits, nil
}

// Make the endianness of bytes compatible with more linux OSs under different processor-architectures,
// according to http://man7.org/linux/man-pages/man2/eventfd.2.html.
var (
	x uint64 = 1
	b        = (*(*[8]byte)(unsafe.Pointer(&x)))[:]
)

// Trigger ...
func (p *Poll) Trigger(job interface{}) error {
	p.jobs.Append(job)

	_, err := syscall.Write(p.wfd, b)
	return err
}

// ControlAdd control add
func (p *Poll) ControlAdd(fd int, events uint32) error {
	return unix.EpollCtl(p.fd, unix.EPOLL_CTL_ADD, fd,
		&unix.EpollEvent{
			Fd:     int32(fd),
			Events: events,
		},
	)
}

// ControlMod control mod
func (p *Poll) ControlMod(fd int, events uint32) error {
	return unix.EpollCtl(p.fd, unix.EPOLL_CTL_MOD, fd,
		&unix.EpollEvent{
			Fd:     int32(fd),
			Events: events,
		},
	)
}

// ControlDel control del
func (p *Poll) ControlDel(fd int, events uint32) error {
	return unix.EpollCtl(p.fd, unix.EPOLL_CTL_DEL, fd,
		&unix.EpollEvent{
			Fd:     int32(fd),
			Events: events,
		},
	)
}
