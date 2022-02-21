package timer

import "time"

type Timer struct {
	startTime time.Time
	timeNow   func() time.Time
}

func NewTimer() *Timer {
	timer := &Timer{
		timeNow: time.Now,
	}
	timer.startTime = timer.timeNow()
	return timer
}

func (t *Timer) GetElapsedTime() time.Duration {
	return t.timeNow().Sub(t.startTime)
}
