package model

import "time"

type Log struct {
	Timestamp time.Time     `json:"timestamp"`
	Latency   time.Duration `json:"latency"`
	Path      string        `json:"path"`
	OS        string        `json:"os"`
}
