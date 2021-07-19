package elio

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"
	//"github.com/rs/zerolog/log"
)

const writingInterval = 10 * time.Millisecond

// ioDefault io model implementation
type ioDefault struct {
	io     *Io
	jobs   SafeSlice
	ctx    context.Context
	cancel context.CancelFunc
}

// NewIoDefault new io default
func NewIoDefault() *ioDefault {
	io := new(ioDefault)
	return io
}

func (m *ioDefault) String() string {
	return fmt.Sprintf("ioDefault::%p", m)
}

func (m *ioDefault) GetIo() *Io {
	return m.io
}

func (m *ioDefault) SetIo(io *Io) {
	m.io = io
	io.SetIoModel(m)
}

// Listen listen
func (m *ioDefault) Listen(addr string) (ok bool) {
	m.GetIo().Host.Wg.Add(1)

	defer func() {
		m.GetIo().Host.Wg.Done()
	}()

	var err error
	if err = m.GetIo().Listener.Listen("tcp", addr); nil == err {
		// AppInfo().Str(LogObject, m.String()).
		// 	Msgf("succeed to listen with url '%s'", m.GetIo().Addr.String())

		if 1 == m.GetIo().InCount.Add(1) {
			m.GetIo().InAddr.Store(addr)
		}

		m.ctx, m.cancel = context.WithCancel(context.Background())

		go m.Running()

		return true
	}

	AppError().Str(LogObject, m.String()).
		Msgf("service url:'%s' failed to listen with error:%s", addr, err.Error())

	return false
}

// Run run
func (m *ioDefault) Run() (ok bool) {
	m.GetIo().Host.Wg.Add(1)

	go m.Writing()

	return true
}

// Read read from session
func (m *ioDefault) Read(n *Session, in []byte) (int, error) {
	return n.conn.Read(in)
}

// writeAll write all
func (m *ioDefault) writeAll(c net.Conn, out []byte) (written int, err error) {
	var w int
	l := len(out)
	for (written < l) && (nil == err) {
		if w, err = c.Write(out[written:]); nil != err {
			break
		}

		written += w
	}

	return written, err
}

func (m *ioDefault) Write(n *Session, out []byte) (written int, err error) {
	var w int
	l := len(out)
	for (written < l) && (nil == err) {
		if w, err = n.GetConn().Write(out[written:]); nil != err {
			break
		}

		written += w
	}

	if 0 < written {
		m.GetIo().Service.OnWrite(n, out[:written])
	}

	if nil != err {
		m.GetIo().Service.OnError(n, err)
		m.Close(n)
	}

	return written, err
}

// PostWrite post write
func (m *ioDefault) PostWrite(n *Session, out []byte) (written int, err error) {
	if nil != out {
		n.IncRef()

		b := GetByteBuffer()
		b.Write(out[written:])
		n.outQueue.Append(b)

		m.jobs.Append(&WriteJob{
			session: n,
		})
		//return m.Write(n, out)

		return len(out), nil
	}

	return 0, nil
}

func (m *ioDefault) Trigger(job interface{}) (err error) {
	m.jobs.Append(job)
	return nil
}

// Shut shut listen
func (m *ioDefault) Shut() {
	if nil != m.GetIo().Listener {
		m.GetIo().Listener.Close()
	}
}

// End end
func (m *ioDefault) End() {
	m.cancel()

	m.CloseAll()

	m.GetIo().Host.Wg.Done()
}

func (m *ioDefault) Close(n *Session) (err error) {
	err = n.GetConn().Close()

	return err
}

func (m *ioDefault) Shutdown(n *Session, how int) error {
	return n.GetConn().Close()
}

func (m *ioDefault) CloseAll() {
	c := func(key string, v interface{}) {
		s := v.(*Session)
		m.Close(s)
	}

	m.GetIo().sessionCmap.IterCb(c)
}

// Running running io
func (m *ioDefault) Running() {
	m.GetIo().Host.Wg.Add(1)

	defer func() {
		m.GetIo().Host.Wg.Done()
	}()

	addr := m.GetIo().InAddr.Load()

	AppInfo().Str(LogObject, m.String()).
		Msgf("succeed to listen with url:%s", addr)
	// if zerolog.InfoLevel < zerolog.GlobalLevel() {
	//	AppInfo().Str(LogObject, m.String()).
	//		Msgf("succeed to listen with url:%s", addr)
	// }

	//m.GetIo().Service.OnListen(m.GetIo())

	for {
		conn, err := m.GetIo().Listener.listener.Accept()
		if nil != err {
			AppError().Str(LogObject, m.String()).Err(err).
				Msgf("service url:'%s' failed to accept", addr)
			break
		}
		conn.(*net.TCPConn).SetNoDelay(true)
		conn.(*net.TCPConn).SetLinger(0)
		//conn.SetKeepAlive(true)
		//conn.SetKeepAlivePeriod(2 * time.Second)
		// set SetReadDeadline
		// err := conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		// if err != nil {
		//     log.Println("SetReadDeadline failed:", err)
		// 	// do something else, for example create new conn
		// }

		session := NewSession(invalidFd, conn, m.GetIo())

		go m.Reading(session)
	}

	//m.GetIo().Service.OnShut(m.GetIo())
}

// Reading reading
func (m *ioDefault) Reading(n *Session) (err error) {
	var bufferLen = m.GetIo().Config.ReadBufferLen
	b := make([]byte, bufferLen)
	var readBytes int

	m.GetIo().sessionCmap.Set(n.String(), n)
	defer func() {
		m.GetIo().Service.OnClose(n, err)
		//n.DisposeNoWait()
		m.GetIo().sessionCmap.Remove(n.String())

		n.DecRef()
	}()

	if err = m.GetIo().Service.OnOpen(n); nil != err {
		return err
	}

	AppDebug().Str(LogObject, m.String()).
		Str(LogSession, n.String()).
		Msg("session opened")

	for {
		readBytes, err = m.Read(n, b)
		if readBytes != 0 {
			m.GetIo().Service.OnRead(n, b[:readBytes])
		}

		if err != nil {
			if err != io.EOF {
				AppError().Str(LogObject, m.String()).Str(LogSession, n.String()).Err(err).
					Msgf("failed to read buffer.len:%d read.len:%d", len(b), readBytes)

				m.GetIo().Service.OnError(n, err)
				m.Close(n)

			} else {
				AppDebug().Str(LogObject, m.String()).
					Str(LogSession, n.String()).
					Msg("session closed")
			}

			break
		}
	}

	return err
}

func (m *ioDefault) Writing() {
	m.GetIo().Host.Wg.Add(1)

	ticker := time.NewTicker(writingInterval)
	defer func() {
		ticker.Stop()
		m.GetIo().Host.Wg.Done()
	}()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			jobs := m.jobs.Fetch()
			for _, job := range jobs {
				job.(Job).Work()
			}
		}
	}
}

// // Writing writing
// func (m *ioDefault) Writing(n *Session) {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			AppPanic().Msgf("failed to writing with runtime.panic:%v", r)
// 			//debug.PrintStack()
// 		}

// 		m.Close(n)
// 	}()

// 	for {
// 		select {
// 		case b, ok := <-n.chanOut:
// 			if !ok {
// 				return
// 			}

// 			// TODO: cppis - HEX dump
// 			if e := log.Debug(); e.Enabled() {
// 				//log.HexDump(b.Bytes(), "dump.out.sending session:%s len:%d", n.String(), b.Len())
// 			}

// 			lenBuffer := b.Len()

// 			sent, err := EventOnWrite(n, b.Bytes())
// 			if (nil == err) && (sent == lenBuffer) {
// 				if nil != m.GetIo().Events.OnWrite {
// 					m.GetIo().Events.OnWrite(n, b.Bytes())
// 				}
// 			}
// 		}
// 	}
// }
