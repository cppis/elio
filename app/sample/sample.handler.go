package main

import "github.com/cppis/elio"

// Query query
func (s *Sample) Query(n *elio.Session, c string) {
	// add event to queue
	// make event
	n.GetIoCore().Host.PostToQueue(&EventQuery{
		context: n.GetContext().(*AppContext),
		command: c,
		host:    n.GetIoCore().Host,
	})
}

// Echo echo
func (s *Sample) Echo(n *elio.Session, m string) {
	// add event to queue
	// make event
	v := &EventPacket{
		session: n,
		context: n.GetContext().(*AppContext),
		message: m,
	}
	n.GetIoCore().Host.PostToQueue(v)
}
