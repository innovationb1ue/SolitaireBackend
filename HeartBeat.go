package main

import (
	"log"
	"time"
)

type HeartBeat struct {
	target         interface{}
	duration       time.Duration
	timer          *time.Timer
	timeOutHandler func()
}

func NewHeartBeat(target interface{}, timeout time.Duration, TimeoutHandlerFunc func()) *HeartBeat {
	timer := time.NewTimer(timeout)
	timer.Stop()
	return &HeartBeat{
		target:         target,
		duration:       timeout,
		timer:          timer,
		timeOutHandler: TimeoutHandlerFunc,
	}
}

func (heartbeat *HeartBeat) Start() {
	heartbeat.timer.Reset(heartbeat.duration)
	select {
	case t := <-heartbeat.timer.C:
		log.Print(t, "Heartbeat timeout triggered")
		heartbeat.timeOutHandler()
	}
}

func (heartbeat HeartBeat) Interrupt() {
	heartbeat.timer.Stop()
	heartbeat.timer.Reset(heartbeat.duration)
}
