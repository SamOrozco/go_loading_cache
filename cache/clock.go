package cache

import "time"

type Clock interface {
	Now() time.Time
}

type LocalClock struct {
}

func (l LocalClock) Now() time.Time {
	return time.Now()
}
