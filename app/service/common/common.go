package common

import "time"

func DurationDays(t time.Time) (days int) {
	return int(time.Now().Sub(t).Hours() / 24)
}
