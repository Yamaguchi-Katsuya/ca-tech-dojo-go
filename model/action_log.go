package model

import "time"

type ActionLog struct {
	Timestamp   time.Time     `json:"timestamp"`
	Path        string        `json:"path"`
	OS          string        `json:"os"`
	ElapsedTime time.Duration `json:"elapsed_time"`
}
