package lib

import (
	"time"
)

// Clock provides time and date functions.
type Clock interface {
	Now() time.Time
	UTC() time.Time
}

// NewClock creates an instance of clock.
func NewClock() Clock {
	return &clock{}
}

type clock struct{}

// Now returns current time.
func (t *clock) Now() time.Time {
	return time.Now()
}

// UTC returns current time in UTC timezone.
func (t *clock) UTC() time.Time {
	return time.Now().UTC()
}
