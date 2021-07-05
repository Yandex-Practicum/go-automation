package grader

import "time"

type TimeMilli int64

func (t TimeMilli) Duration() time.Duration {
	return time.Duration(t) * time.Millisecond
}
