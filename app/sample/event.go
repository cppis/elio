package main

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/cppis/elio"
)

// EventOpen event open
type EventOpen struct {
	session *elio.Session
	host    *elio.IoHost
}

// String string
func (e *EventOpen) String() string {
	return fmt.Sprintf("EventOpen::%p", e)
}

// Handle handle
func (e *EventOpen) Handle() {
	c := NewAppContext(e.session, e.host)
	e.session.SetContext(c)

	c.SetState(Opened)

	division, ok := e.host.Register(unsafe.Pointer(e.session), c)
	if false == ok {
		elio.AppError().Str(elio.LogObject, e.String()).
			Str(elio.LogSession, e.session.String()).
			Str(elio.LogContext, c.String()).
			Msgf("failed to register division:%p:%d", unsafe.Pointer(c), division)

	} else {
		elio.AppDebug().Str(elio.LogObject, e.String()).
			Str(elio.LogSession, e.session.String()).
			Str(elio.LogContext, c.String()).
			Msgf("succeed to register key:%p value:%p division:%d", unsafe.Pointer(e.session), unsafe.Pointer(c), division)

		c.DivIndex.Store(division)
	}

	elio.AppDebug().Str(elio.LogObject, e.String()).
		Str(elio.LogSession, e.session.String()).
		Str(elio.LogContext, c.String()).Msg("event.open")
}

// EventClose event close
type EventClose struct {
	session *elio.Session
	host    *elio.IoHost
}

// String string
func (e *EventClose) String() string {
	return fmt.Sprintf("EventClose::%p", e)
}

// Handle handle
func (e *EventClose) Handle() {
	var c *AppContext

	defer func() {
		if nil != c {
			c.GetSession().SetContext(nil)
			//c.GetSession().DecRef()

			c.Release() //SubRef()
		}
	}()

	if nil != e.session.GetContext() {
		c = e.session.GetContext().(*AppContext)
		division := c.DivIndex.Load()
		if false == e.host.Unregister(unsafe.Pointer(e.session), division) {
			elio.AppError().Str(elio.LogObject, e.String()).
				Str(elio.LogSession, e.session.String()).
				Str(elio.LogContext, c.String()).
				Msgf("failed to unregister division:%d", division)

		} else {
			elio.AppDebug().Str(elio.LogObject, e.String()).
				Str(elio.LogSession, e.session.String()).
				Str(elio.LogContext, c.String()).
				Msgf("succeed to unregister key:%p division:%d", unsafe.Pointer(e.session), division)

			c.DivIndex.Store(elio.InvalidDivIndex)
		}

	} else {
		elio.AppDebug().Str(elio.LogObject, e.String()).
			Str(elio.LogSession, e.session.String()).
			Msgf("event.close with no context")

	}
}

// EventQuery event query
type EventQuery struct {
	context *AppContext
	command string
	host    *elio.IoHost
}

// String string
func (e *EventQuery) String() string {
	return fmt.Sprintf("EventQuery::%p", e)
}

// Handle handle
func (e *EventQuery) Handle() {
	s := e.context.GetSession()

	elio.AppDebug().Str(elio.LogObject, e.String()).
		Str(elio.LogSession, s.String()).
		Str(elio.LogContext, e.context.String()).Msgf("event.query %s", e.command)

	var out string
	out = fmt.Sprintf("query\ncommand: %s", e.command)

	words := strings.Fields(e.command)
	switch words[0] {
	case "countuser":
		// countuser
		c := e.host.CountUser()
		out += fmt.Sprintf("\ncount: %d", c)
	case "listuser":
		// listuser
		l := e.host.ListUser()
		out += fmt.Sprintf("\nlist: %s", l)
	case "finduser":
		// finduser {uid}
		// countuser
		var uid elio.UID
		uid.FromString(words[1])
		u, ok := e.host.FindUser(uid)
		if ok {
			out += fmt.Sprintf("\nfound: %+v", *u.(*AppContext))
		} else {
			out += fmt.Sprintf("\nfound: not")
		}
	case "countroom":
		// countroom
	case "listroom":
		// listroom
	case "findroom":
		// findroom {room id}
	default:
		out = fmt.Sprintf("query\nresult: invalid command%s\n\n", words[0])
	}

	out += fmt.Sprintf("\n\n")
	e.context.Write([]byte(out))
}

// EventPacket event packet
type EventPacket struct {
	session *elio.Session
	context *AppContext
	message string
}

// String string
func (e *EventPacket) String() string {
	return fmt.Sprintf("EventPacket::%p", e)
}

// Echo echo
func (e *EventPacket) Echo(s *elio.Session) (err error) {
	elio.AppInfo().Str(elio.LogObject, e.String()).
		Str(elio.LogSession, s.String()).
		Str(elio.LogContext, e.context.String()).
		Msgf("uid:%v receive echo", e.context.UID.ToInt())

	if elio.DebugEnabled() {
		elio.AppDebug().Str(elio.LogSession, s.String()).
			Str(elio.LogContext, e.context.String()).
			Msg("succeed to unpack echo")
	}

	//return e.packet.Pack(protocol.CodeEchoResponse.ToUint16(), e.packet.Flags, e.packet.Body)
	//return &req, err

	out := fmt.Sprintf("echo\n%s\n\n", e.message)
	e.context.Write([]byte(out))
	return err
}
