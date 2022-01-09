package main

import (
	"github.com/google/uuid"
	"log"
)

type Room struct {
	RoomUUID  string
	Players   map[string]*Player
	broadcast chan map[string]interface{}
}

func newRoom() *Room {
	return &Room{
		RoomUUID:  uuid.NewString(),
		Players:   make(map[string]*Player),
		broadcast: make(chan map[string]interface{}),
	}
}

func (r *Room) AddPlayer(p *Player) {
	r.Players[p.Id] = p
}

func (r *Room) Run() {
	for {
		select {
		case message := <-r.broadcast:
			{
				log.Print("Room broadcast receive message: ", message)
				for _, client := range r.Players {
					if message["sender"] == client {
						continue
					}
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(r.Players, client.Id)
					}
				}
			}
		}
	}
}
