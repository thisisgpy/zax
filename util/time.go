package util

import "time"

func Now() *time.Time {
	now := time.Now()
	return &now
}

func NowString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
