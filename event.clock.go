package elio

import "time"

// Clock clock
type Clock struct {
	alarmMap map[string]Alarm
	prev     time.Time
}

// NewClock new clock
func NewClock() *Clock {
	c := new(Clock)
	if nil != c {
		c.alarmMap = make(map[string]Alarm)
	}

	c.prev = time.Now()
	return c
}

// Register regiseter alarm
func (c *Clock) Register(name string, alarm Alarm) bool {
	_, ok := c.alarmMap[name]
	if ok {
		return false
	}

	c.alarmMap[name] = alarm
	return true
}

// RegisterWithTime regiseter alarm with time
func (c *Clock) RegisterWithTime(name string, alarm Alarm, t time.Time) bool {
	_, ok := c.alarmMap[name]
	if ok {
		return false
	}

	c.alarmMap[name] = alarm
	c.prev = t
	return true
}

// Unregister unregister alarm
func (c *Clock) Unregister(name string) bool {
	_, ok := c.alarmMap[name]
	if ok {
		delete(c.alarmMap, name)
	}

	return ok
}

// UnregisterAll unregister all alarm
func (c *Clock) UnregisterAll() {
	for n := range c.alarmMap {
		delete(c.alarmMap, n)
	}
}

// Reset reset
func (c *Clock) Reset(name string) bool {
	a, ok := c.alarmMap[name]
	if ok {
		a.Reset()
	}

	return ok
}

// Update update
func (c *Clock) Update(t time.Time) {
	delta := t.Sub(c.prev)

	for n, a := range c.alarmMap {
		if true == a.Check(t, delta) {
			if false == a.Ring(n, t, c) {
				delete(c.alarmMap, n)
			}
		}
	}

	c.prev = t
}

// Count count
func (c *Clock) Count() int {
	return len(c.alarmMap)
}
