package main

import (
	"testing"
	"time"
)

func TestNewRoom(t *testing.T) {
	unreg := make(chan string)
	_ = NewRoom(unreg)
}

func TestRoom_Run(t *testing.T) {
	unreg := make(chan string)
	room := NewRoom(unreg)
	roomId := room.RoomUUID
	go room.Run()
	time.Sleep(7 * time.Second)
	if roomId != <-unreg {
		t.FailNow()
	}
}

func TestRoom_AddPlayer(t *testing.T) {
	var err error
	unreg := make(chan string)
	room := NewRoom(unreg)
	go room.Run()
	Player1 := NewPlayer("Tester")
	err = room.AddPlayer(Player1)
	if err != nil {
		t.FailNow()
	}
	Player2 := NewPlayer("Tester2")
	err = room.AddPlayer(Player2)
	if err != nil {
		t.FailNow()
	}
	Player3 := NewPlayer("Tester3")
	err = room.AddPlayer(Player3)
	if err == nil {
		t.FailNow()
	}
}
