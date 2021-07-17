package main

import (
	"fmt"
	"time"

	"github.com/cppis/elio"
)

const (
	// AlarmNameShow alarm show
	AlarmNameShow = "main.show"
)

const (
	// MetricEventsf events metric format
	MetricEventsf = "%s.events.%X"
	// MetricEventCountf event count metric format
	MetricEventCountf = "%s.eventcount.%X"
	// MetricAlarmf alarm metric format
	MetricAlarmf = "%s.alarms.%X"
	// MetricRoomCurrentf current room metric format
	MetricRoomCurrentf = "%s.rooms"
	// MetricUserCurrentf current user metric format
	MetricUserCurrentf = "%s.users"
)

// AlarmDot alarm show
type AlarmDot struct {
	duration time.Duration
	elapsed  time.Duration
}

// NewAlarmDot new alarm dot
func NewAlarmDot(d time.Duration) elio.Alarm {
	a := new(AlarmDot)
	a.init(d)

	return a
}

// String string
func (a *AlarmDot) String() string {
	return fmt.Sprintf("AlarmDot::%p", a)
}

// init
func (a *AlarmDot) init(d time.Duration) {
	a.duration = d
	a.elapsed = 0
}

// Reset reset
func (a *AlarmDot) Reset() {
	a.elapsed = 0
}

// Check check
func (a *AlarmDot) Check(t time.Time, d time.Duration) bool {
	a.elapsed += d
	if a.duration < a.elapsed {
		a.elapsed = 0
		return true
	}

	return false
}

// Ring ring
func (a *AlarmDot) Ring(name string, t time.Time, c *elio.Clock) bool {
	elio.AppDebug().Str(elio.LogObject, a.String()).Msgf("rings alarm")

	// false 를 리턴할 경우, alarm 컨테이너에서 삭제됩니다.

	//fmt.Printf(".")
	return true
}
