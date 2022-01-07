package main

import "github.com/google/uuid"

type Room struct {
	RoomUUID  string
	Players   map[string]*Player
	broadcast chan []byte
}

func newRoom() *Room {
	return &Room{
		RoomUUID:  uuid.NewString(),
		Players:   make(map[string]*Player),
		broadcast: make(chan []byte),
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
				for _, client := range r.Players {
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
