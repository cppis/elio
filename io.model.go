package elio

// Io I/O model
type IoModel interface {
	GetIo() *Io
	SetIo(c *Io)
	Listen(addr string) bool
	Run() bool
	Shut() // close listen
	End()  // end service
	Read(n *Session, in []byte) (receipt int, err error)
	Write(n *Session, out []byte) (sent int, err error)
	PostWrite(n *Session, out []byte) (sent int, err error)
	Trigger(job interface{}) error
	Close(n *Session) error
	Shutdown(n *Session, how int) error
	CloseAll()
}
