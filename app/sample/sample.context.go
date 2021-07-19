package main

import (
	"encoding/hex"
	"fmt"

	"github.com/cppis/elio"
	"go.uber.org/atomic"
)

// SessionState session state type definition
type SessionState int

const (
	// Closed closed
	Closed SessionState = iota
	// Opened opened
	Opened
	// IOErrored IO Errored
	IOErrored
)

// AppContext app context
type AppContext struct {
	elio.Context

	IoHost       *elio.IoHost
	UID          elio.UID
	SessionState SessionState
	DivIndex     atomic.Uint32 // division index
	Sleep        bool
}

// NewAppContext create new context
func NewAppContext(s *elio.Session, h *elio.IoHost) *AppContext {
	c := new(AppContext)
	c.Init(s, h)
	return c
}

// Init init
func (c *AppContext) Init(s *elio.Session, h *elio.IoHost) {
	c.IoHost = h
	c.Session = s
	c.UID = elio.UIDInvalid
	c.SessionState = Closed
	c.DivIndex.Store(elio.InvalidDivIndex)
}

// Release release
func (c *AppContext) Release() {
	s := c.GetSession()
	if nil != s {
		//fmt.Printf("release\n")
		s.SetContext(nil)
	}

	//c.SetSession(nil)
}

// Copy copy
func (c *AppContext) Copy(ctx *AppContext) {
	//c.UID = ctx.UID
	c.Sleep = ctx.Sleep
	c.SessionState = ctx.SessionState
}

// IsValidState is valid state
func (c *AppContext) IsValidState() bool {
	if Opened == c.SessionState {
		return true
	}

	return false
}

// GetState get state
func (c *AppContext) GetState() SessionState {
	return c.SessionState
}

// SetState set state
func (c *AppContext) SetState(state SessionState) {
	c.SessionState = state
}

// String string
func (c *AppContext) String() string {
	return fmt.Sprintf("AppContext::%p:%d", c, c.UID.ToInt())
}

// Fnv64
func (c *AppContext) Fnv64() elio.Fnv64 {
	var fnv elio.Fnv64
	fnv.FromString(fmt.Sprintf("%p", c))
	return fnv
}

// SetUID set user id
func (c *AppContext) SetUID(uid elio.UID) {
	c.UID = uid
}

// Write write
func (c *AppContext) Write(o []byte) {
	if c.IsValidState() {
		var w int
		var err error
		if w, err = c.GetSession().Write(o); nil != err {
			elio.AppError().Str(elio.LogObject, c.String()).
				Str(elio.LogSession, c.GetSession().String()).
				Int64(elio.LogUuid, c.UID.ToInt()).Err(err).
				Msg("failed to write")

			if retErr := c.GetSession().GetIo().Shutdown(c.GetSession(), elio.ShutRd); nil != retErr {
				//elio.AppError().Str(elf.LogObject, c.String()).
				//	Str(elf.LogSession, c.GetSession().String()).
				//	Int64(elf.LogUuid, c.UID.ToInt()).Err(retErr).
				//	Msg("failed to shutdown of write")
			}
			c.SetState(IOErrored)

		} else {
			if elio.TraceEnabled() {
				elio.AppTrace().Str(elio.LogObject, c.String()).
					Str(elio.LogSession, c.GetSession().String()).
					Int64(elio.LogUuid, c.UID.ToInt()).
					Str(elio.LogPayload, hex.Dump(o)).
					Msgf("succeed to write %d bytes", w)
			} else if elio.DebugEnabled() {
				elio.AppDebug().Str(elio.LogObject, c.String()).
					Str(elio.LogSession, c.GetSession().String()).
					Int64(elio.LogUuid, c.UID.ToInt()).
					Msgf("succeed to write %d bytes", w)
			}
		}
	}
}

// PostWrite post write
func (c *AppContext) PostWrite(o []byte) {
	if c.IsValidState() {
		var w int
		var err error
		if w, err = c.GetSession().PostWrite(o); nil != err {
			elio.AppError().Str(elio.LogObject, c.String()).
				Str(elio.LogSession, c.GetSession().String()).
				Int64(elio.LogUuid, c.UID.ToInt()).Err(err).
				Msg("failed to post.write")

			if retErr := c.GetSession().GetIo().Shutdown(c.GetSession(), elio.ShutRd); nil != retErr {
				//elio.AlogError().Str(elio.LogObject, c.String()).
				//	Str(elf.LogSession, c.GetSession().String()).
				//	Int64(elio.LogUuid, c.UID.ToInt()).Err(retErr).
				//	Msg("failed to shutdown of write")
			}
			c.SetState(IOErrored)

		} else {
			if elio.TraceEnabled() {
				elio.AppTrace().Str(elio.LogObject, c.String()).
					Str(elio.LogSession, c.GetSession().String()).
					Int64(elio.LogUuid, c.UID.ToInt()).
					Str(elio.LogPayload, hex.Dump(o)).
					Msgf("succeed to post.write %d bytes", w)
			} else if elio.DebugEnabled() {
				elio.AppDebug().Str(elio.LogObject, c.String()).
					Str(elio.LogSession, c.GetSession().String()).
					Int64(elio.LogUuid, c.UID.ToInt()).
					Msgf("succeed to post.write %d bytes", w)
			}
		}
	}
}

// DirectOut direct out
func (c *AppContext) DirectOut(o []byte) (w int, err error) {
	if c.IsValidState() {
		if w, err = c.GetSession().Write(o); nil != err {
			elio.AppError().Str(elio.LogObject, c.String()).
				Str(elio.LogSession, c.GetSession().String()).
				Int64(elio.LogUuid, c.UID.ToInt()).Err(err).
				Msg("failed to write")

			if retErr := c.GetSession().GetIo().Shutdown(c.GetSession(), elio.ShutRd); nil != retErr {
				elio.AppError().Str(elio.LogObject, c.String()).
					Str(elio.LogSession, c.GetSession().String()).
					Int64(elio.LogUuid, c.UID.ToInt()).Err(retErr).
					Msg("failed to shutdown of write")
			}
			c.SetState(IOErrored)

		} else {
			if elio.DebugEnabled() {
				elio.AppDebug().Str(elio.LogObject, c.String()).
					Str(elio.LogSession, c.GetSession().String()).
					Int64(elio.LogUuid, c.UID.ToInt()).
					Str(elio.LogPayload, hex.Dump(o)).
					Msgf("succeed to write:%d", w)
			}
		}
	}

	return w, err
}

// CheckValid check valid
func (c *AppContext) CheckValid() bool {
	if c.SessionState == Opened &&
		c.Sleep == false &&
		nil != c.GetSession() {
		return true
	}

	return false
}
