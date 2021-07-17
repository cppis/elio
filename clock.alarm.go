package elio

import (
	"time"
)

// Alarm alarm interface
type Alarm interface {
	Reset()
	Check(t time.Time, d time.Duration) bool
	Ring(name string, t time.Time, c *Clock) bool // returns continue flag
}
