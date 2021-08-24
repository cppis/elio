package elio

import (
	"fmt"
	"net"
	"sync"
	"time"
	"unsafe"
)

// DivCallback metric callback
type DivCallback func(t time.Time, v interface{}) error

// IoHost host
type IoHost struct {
	Ios        []*Io
	Wg         *sync.WaitGroup
	Url        *net.TCPAddr
	Users      map[UID]interface{}
	Rooms      map[string]interface{} //*Room
	EventQueue *EventQueue
	Divs       *DivMap
	DivChecked uint32
}

// NewIoHost new host
func NewIoHost(divs int) (h *IoHost) {
	if h = new(IoHost); nil != h {
		h.Wg = new(sync.WaitGroup)
		h.Users = make(map[UID]interface{})
		h.Rooms = make(map[string]interface{})
		h.EventQueue = NewEventQueue()
		h.Divs = NewDivMap(uint32(divs))
	}

	return h
}

// String string
func (h *IoHost) String() string {
	return fmt.Sprintf("IoHost::%p", h)
}

// End end
func (h *IoHost) End() {
	for _, c := range h.Ios {
		c.End()
	}
}

// Terminate terminate
func (h *IoHost) Terminate(safe bool) {
	for _, c := range h.Ios {
		//if true == safe {
		//	c.SafeTerminate()
		//} else {
		c.Terminate()
		//}
	}
}

// FindUser find user
func (h *IoHost) FindUser(uid UID) (interface{}, bool) {
	c, ok := h.Users[uid]
	return c, ok
}

// EnterUser enter user
func (h *IoHost) EnterUser(uid UID, i interface{}) {
	h.Users[uid] = i
}

// FindOrEnterUser find or add user
func (h *IoHost) FindOrEnterUser(uid UID, i interface{}) bool {
	if _, ok := h.FindUser(uid); ok {
		return false
	}

	h.EnterUser(uid, i)
	return true
}

// LeaveUser leave user
func (h *IoHost) LeaveUser(uid UID) {
	delete(h.Users, uid)
}

// CountUser count user
func (h *IoHost) CountUser() int {
	return len(h.Users)
}

// ListUser list user
func (h *IoHost) ListUser() (out string) {
	for k := range h.Users {
		out += fmt.Sprintf("%s, ", k.ToString())
	}
	return out
}

// FindRoom find room
func (h *IoHost) FindRoom(key string) interface{} {
	if r, ok := h.Rooms[key]; ok {
		return r
	}

	return nil
}

// AddRoom add room
func (h *IoHost) AddRoom(key string, room interface{}) {
	h.Rooms[key] = room
}

// DelRoom delete room
func (h *IoHost) DelRoom(key string) bool {
	r := h.FindRoom(key)
	if nil == r {
		return false
	}

	delete(h.Rooms, key)
	return true
}

// GetDivision get division
func (h *IoHost) GetDivision() *DivMap {
	return h.Divs
}

// SetDivision set division
func (h *IoHost) SetDivision(k uint64, v interface{}) (d uint32, ok bool) {
	d, _ = h.Divs.Set(k, v)
	if InvalidDivIndex == d {
		ok = false
	}

	return d, true
}

// DelDivision del division
func (h *IoHost) DelDivision(d uint32, k uint64) (ok bool) {
	_, ok = h.Divs.Del(d, k)
	return ok
}

// GetEventQueue get event queue
func (h *IoHost) GetEventQueue() *EventQueue {
	return h.EventQueue
}

// Dispatching dispatching
func (h *IoHost) Dispatching(t time.Time, limit int) (int, int) {
	return h.EventQueue.Dispatching(t, limit)
}

// Register register
func (h *IoHost) Register(p unsafe.Pointer, v interface{}) (d uint32, ok bool) {
	var fnv Fnv64
	fnv.FromPointer(p)
	d, ok = h.SetDivision(fnv.ToUint64(), v)

	return d, ok
}

// Unregister unregister
func (h *IoHost) Unregister(p unsafe.Pointer, d uint32) (ok bool) {
	var fnv Fnv64
	fnv.FromPointer(p)
	ok = h.DelDivision(d, fnv.ToUint64())

	return ok
}

// PostToQueue post to event queue
func (h *IoHost) PostToQueue(i interface{}) {
	h.EventQueue.Inject(i)
}

// RunDivision run division
func (h *IoHost) RunDivision(t time.Time, callback DivCallback) {
	m, err := h.Divs.Get(h.DivChecked)
	if nil == err {
		//var erased []eraseNode

		for _, v := range m.Map {
			_ = callback(t, v)
			//err := callback(t, v)
			//if nil != err {
			//	erased = append(erased, eraseNode{
			//		k: k,
			//		v: v,
			//	})
			//}
		}

		// for _, e := range erased {
		// 	//division := c.DivisionIndex.Load()
		// 	ok := l.DelDivision(l.DivisionChecked, e.k)
		// 	if false == ok {
		// 		AppError().Str(elio.LogObject, l.String()).Err(err).
		// 			Msgf("failed to unregister context:%v:%d", e.k, l.DivisionChecked)

		// 	} else {
		// 		AppDebug().Str(elio.LogObject, l.String()).
		// 			Msgf("succeed to unregister context:%v:%d", e.k, l.DivisionChecked)
		// 	}
		// }
	}

	h.DivChecked++
	count := h.Divs.Count()
	h.DivChecked = h.DivChecked % uint32(count)
}

// PostAllToQueue post all to queue
func (h *IoHost) PostAllToQueue(l []interface{}) {
	h.GetEventQueue().AppendAll(l...)
}

const (
	defaultDivCount = 50
)

// Host host
func Host(config ConfigIo, service Service) (host *IoHost, err error) {
	host = NewIoHost(defaultDivCount)
	host.Ios = make([]*Io, config.InCount)

	for i := 0; i < config.InCount; i++ {
		io := ProvideIo(host, config, service)
		host.Ios[i] = io

		var addr *net.TCPAddr
		if addr, err = net.ResolveTCPAddr("tcp", io.Config.InURL); nil == err {
			io.Run(addr)
			host.Url = addr
		}
	}

	return host, err
}
