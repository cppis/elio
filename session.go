package elio

import (
	"errors"
	"fmt"
	"net"
	"reflect"

	"go.uber.org/atomic"
)

const invalidFd = -1

// Session session
type Session struct {
	refCount   atomic.Int32 // reference counting
	fd         int          // file descriptor
	sa         Sockaddr     // remote socket address
	conn       net.Conn
	context    interface{}
	stats      interface{}
	ioCore     *IoCore
	addrLocal  net.Addr // local addre
	addrRemote net.Addr // remote addr
	ipRemote   string
	buffer     []byte
	bufferKeep *Buffer
	outQueue   SafeSlice
}

// NewSession new session
func NewSession(f int, c net.Conn, i *IoCore) *Session {
	n := new(Session) //poolSession.Get().(*Session) //
	if nil != n {
		AppDebug().Str(LogSession, n.String()).Msg("new session")

		n.fd = f
		n.conn = c
		n.ioCore = i
		n.buffer = make([]byte, i.Config.ReadBufferLen)
		n.bufferKeep = &Buffer{}
		n.ipRemote, _ = n.GetAddr()

		n.IncRef()
		return n
	}
	return n
}

// DeleteSession delete session
func DeleteSession(n *Session) {
	AppDebug().Str(LogSession, n.String()).Msg("delete session")

	//n.init()
	//poolSession.Put(n)
}

// GetFdFromLisener this function is linux only
func GetFdFromLisener(listener net.Listener) int {
	fdValue := reflect.Indirect(reflect.Indirect(reflect.ValueOf(listener)).FieldByName("fd"))
	//return uintptr(fdValue.FieldByName("sysfd").Int())
	return int(fdValue.FieldByName("sysfd").Int())
}

// GetFdFromConn this function is linux only
func GetFdFromConn(conn net.Conn) int {
	//tls := reflect.TypeOf(conn.UnderlyingConn()) == reflect.TypeOf(&tls.Conn{})
	// Extract the file descriptor associated with the connection
	//connVal := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn").Elem()
	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	//if tls {
	//	tcpConn = reflect.Indirect(tcpConn.Elem())
	//}
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")

	return int(pfdVal.FieldByName("sysfd").Int())
}

// String ...
func (n *Session) String() string {
	return fmt.Sprintf("Session::%p;%d;%v", n, n.GetFd(), n.GetRemoteIP())
}

// init init
func (n *Session) init() {
	n.fd = 0
	n.conn = nil
	n.context = nil
	n.stats = nil
	n.ioCore = nil
	//n.buffer = nil
	n.bufferKeep.Clear(-1)
	n.addrLocal = nil
	n.addrRemote = nil
	n.ipRemote = ""
	n.bufferKeep = nil
	//_ = n.outQueue.Fetch()
}

// IncRef increase reference count
func (n *Session) IncRef() int32 {
	return n.refCount.Inc()
}

// SubRef subscribe reference count
func (n *Session) SubRef(ref int32) int32 {
	r := n.refCount.Sub(ref)
	if 0 == r {
		DeleteSession(n)
	}

	return r
}

// DecRef decrease reference count
func (n *Session) DecRef() int32 {
	r := n.refCount.Dec()
	if 0 == r {
		DeleteSession(n)
	}

	return r
}

// GetFd get fd
func (n *Session) GetFd() int {
	return n.fd
}

// GetConn get conn
func (n *Session) GetConn() net.Conn {
	return n.conn
}

// GetService get service
func (n *Session) GetIoCore() *IoCore {
	return n.ioCore
}

// GetContext get context
func (n *Session) GetContext() interface{} {
	return n.context
}

// SetContext set context
func (n *Session) SetContext(c interface{}) {
	n.context = c
}

// GetStats get stats
func (n *Session) GetStats() interface{} {
	return n.stats
}

// SetStats set stats
func (n *Session) SetStats(s interface{}) {
	n.stats = s
}

// GetRemoteIP get remote ip
func (n *Session) GetRemoteIP() string {
	return n.ipRemote
}

// GetKeepBuffer get keep buffer
func (n *Session) GetKeepBuffer() *Buffer {
	return n.bufferKeep
}

// Close close
func (n *Session) Close() error {
	return n.ioCore.io.Close(n)
}

// Shutdown shutdown
func (n *Session) Shutdown(how int) error {
	return n.ioCore.io.Shutdown(n, how)
}

// Write write
func (n *Session) Write(out []byte) (sent int, err error) {
	defer func() {
		if r := recover(); r != nil {
			AppError().Str(LogObject, n.String()).
				Msgf("panic in session.Write fd:%d with recover:%s", n.fd, r)

			sent = -1
			err = errors.New("elf.net failed to session.Write with panic")
		}
	}()

	return n.ioCore.io.Write(n, out)
}

// PostWrite post write
func (n *Session) PostWrite(out []byte) (sent int, err error) {
	return n.ioCore.io.PostWrite(n, out)
}

// CountOutQueue count out queue
func (n *Session) CountOutQueue() int {
	return n.outQueue.Count()
}

// set keep alive
// func SetKeepAlive(fd, secs int) error {
// 	if err := unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_KEEPALIVE, 1); err != nil {
// 		return err
// 	}
// 	if err := unix.SetsockoptInt(fd, unix.IPPROTO_TCP, unix.TCP_KEEPINTVL, secs); err != nil {
// 		return err
// 	}
// 	return unix.SetsockoptInt(fd, unix.IPPROTO_TCP, unix.TCP_KEEPIDLE, secs)
// }
