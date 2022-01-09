package main

import "time"

type HeartBeat struct {
	target   interface{}
	duration time.Duration
	timer    *time.Timer
}

func NewHeartBeat(target interface{}, timeout time.Duration, TimeoutHandlerFunc func()) *HeartBeat {
	timer := time.NewTimer(timeout)
	timer.Stop()
	return &HeartBeat{
		target:   target,
		duration: timeout,
		timer:    timer,
	}
}

func (heartbeat *HeartBeat) Start() {
	heartbeat.timer.Reset(heartbeat.duration)
}
