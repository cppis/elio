package elio

import (
	"time"
)

// EventQueue event queue
type EventQueue struct {
	awaitQueue *UnsafeSlice
	eventQueue *SafeSlice
	events     []interface{}
	//partitionedSessions *PartitionedSlice
	clock    *Clock
	prevTime time.Time
}

// QueueCapacity queue default capacity
const QueueCapacity int = 5000

// NewEventQueue new event queue
func NewEventQueue() *EventQueue {
	if q := new(EventQueue); nil != q {
		q.awaitQueue = NewUnsafeSlice(QueueCapacity)
		q.eventQueue = NewSafeSlice(QueueCapacity)
		//q.partitionedSessions = NewPartitionedSlice(20)
		q.clock = NewClock()
		return q
	}

	return nil
}

// Dispatching dispatching
func (q *EventQueue) Dispatching(t time.Time, l int) (int, int) {
	// clock 처리
	q.clock.Update(t)

	var count int
	// event 처리
	events := q.eventQueue.Fetch()
	if 0 < len(events) {
		q.events = append(q.events, events...)
	}

	qLen := len(q.events)
	if 0 < qLen {
		var fetch []interface{}
		//fmt.Printf("begin[%p] limit:%d, q.len:%d\n", q, l, qLen)
		if l < qLen {
			//fmt.Printf("\t[%p] limit:%d, count:%d, q.len:%d\n", q, l, l, qLen)
			count = l
			fetch, q.events = q.events[:count], q.events[count:]

		} else {
			fetch = q.events
			count = len(fetch)
			//fmt.Printf("\t[%p] limit:-, count:%d, q.len:%d\n", q, count, qLen)
			q.events = nil //make([]interface{}, 0, l)
		}

		//handled := 0
		for _, e := range fetch {
			e.(Event).Handle()
			//handled++
		}

		//fmt.Printf("[%p] end handled:%d, q.len:%d\n", q, handled, len(q.events))
	}

	q.prevTime = t

	return q.clock.Count(), count
}

/*/
func (q *EventQueue) Dispatching(t time.Time, l int) (int, int) {
	// clock 처리
	q.clock.Update(t, t.Sub(q.prevTime))

	// event 처리
	events := q.eventQueue.FetchWithLimit(l) //Fetch() //

	for _, e := range events {
		e.(Event).Handle()
	}

	q.prevTime = t

	return q.clock.Count(), len(events)
}
//*/

// Inject inject
func (q *EventQueue) Inject(event interface{}) {
	q.eventQueue.Append(event)
}

// InjectToAwait inject to await
func (q *EventQueue) InjectToAwait(event interface{}) {
	q.awaitQueue.Append(event)
}

// AppendAll append all
func (q *EventQueue) AppendAll(events ...interface{}) {
	q.eventQueue.AppendAll(events...)
}

// Paste paste
func (q *EventQueue) Paste(events []interface{}) {
	q.eventQueue.Paste(events)
}

// GetClock get clock
func (q *EventQueue) GetClock() *Clock {
	return q.clock
}

// Convey convey await to event queue
func (q *EventQueue) Convey() int {
	e := q.awaitQueue.Fetch()
	l := len(e)
	if 0 < l {
		q.ConveySlice(e...)
	}
	return l
}

// ConveySlice convey slice to event queue
func (q *EventQueue) ConveySlice(s ...interface{}) {
	q.eventQueue.AppendAll(s...)
}
