package main

import "github.com/cppis/elio"

// Query query
func (s *Sample) Query(n *elio.Session, c string) {
	// add event to queue
	// make event
	n.GetIo().Host.PostToQueue(&EventQuery{
		context: n.GetContext().(*AppContext),
		command: c,
		host:    n.GetIo().Host,
	})
}

// Sample sample
func (s *Sample) Echo(n *elio.Session, m string) {
	// add event to queue
	// make event
	v := &EventPacket{
		session: n,
		context: n.GetContext().(*AppContext),
		message: m,
	}
	n.GetIo().Host.PostToQueue(v)
}
