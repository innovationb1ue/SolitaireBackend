package main

import (
	"github.com/google/uuid"
	"log"
	"time"
)

type Room struct {
	RoomUUID       string
	Players        map[string]*Player
	broadcast      chan map[string]interface{}
	UnregisterChan chan string // send self UUID to this chan to unregister this channel
}

func newRoom() *Room {
	return &Room{
		RoomUUID:  uuid.NewString(),
		Players:   make(map[string]*Player),
		broadcast: make(chan map[string]interface{}),
	}
}

func (r *Room) SetUnregisterChan(UnregisterChan chan string) {
	r.UnregisterChan = UnregisterChan
}

func (r *Room) AddPlayer(p *Player) {
	r.Players[p.Id] = p
}

func (r *Room) Run() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case message := <-r.broadcast:
			{
				// broadcast the original message to all other players
				log.Print("Room broadcast receive message: ", message)
				for _, client := range r.Players {
					if message["sender"].(*Player) == client {
						continue
					}
					client.send <- message
				}
			}
		case _ = <-ticker.C:
			{
				PlayerCount := len(r.Players)
				if PlayerCount == 0 {
					r.UnregisterChan <- r.RoomUUID
					return
				}
			}
		}

	}
}
