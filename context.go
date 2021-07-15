package elio

import "fmt"

// Contexter interface
type Contexter interface {
	GetSession() *Session
	String() string
}

// Context implementation
type Context struct {
	Session *Session
}

// NewContext create new context
func NewContext(n *Session) *Context {
	c := new(Context)
	c.Session = n
	return c
}

// GetSession get session
func (c *Context) GetSession() *Session {
	return c.Session
}

// SetSession set session
func (c *Context) SetSession(s *Session) {
	c.Session = s
}

func (c *Context) String() string {
	return fmt.Sprintf("Context::%p", c)
}
