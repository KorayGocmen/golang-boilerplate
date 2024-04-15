package duration

import "time"

func Seconds(t int) time.Duration {
	return time.Duration(t) * time.Second
}
