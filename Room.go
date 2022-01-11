package main

import (
	"github.com/google/uuid"
	"log"
	"time"
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
	ticker := time.NewTicker(60 * time.Second)
	for {
		select {
		case message := <-r.broadcast:
			{
				//action := message["action"].(string)
				//CardLeft := int(message["card_left"].(float64))
				//log.Print(action, CardLeft)
				log.Print("Room broadcast receive message: ", message)
				for _, client := range r.Players {
					if message["sender"].(*Player) == client {
						continue
					}
					log.Print("Room broadcast message")
					client.send <- message
				}
			}
		case _ = <-ticker.C:
			{
				aliveFlag := true
				// todo: check player status and destroy the room if no player is alive
				if aliveFlag {
					continue
				}
			}
		}

	}
}

func (r *Room) Destroy() {
	r = nil
}
