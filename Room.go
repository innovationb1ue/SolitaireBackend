package main

import (
	"errors"
	"github.com/google/uuid"
	"log"
	"time"
)

type Room struct {
	RoomUUID   string
	Players    map[string]*Player
	broadcast  chan map[string]interface{}
	Deck       [][]map[string]interface{}
	unRegister chan<- string
}

func newRoom(unRegister chan<- string) *Room {
	return &Room{
		RoomUUID:   uuid.NewString(),
		Players:    make(map[string]*Player),
		broadcast:  make(chan map[string]interface{}),
		Deck:       nil,
		unRegister: unRegister,
	}
}

func (r *Room) AddPlayer(p *Player) error {
	if len(r.Players) >= 2 {
		log.Println("more than 2 player in the same room not supported")
		return errors.New("too many players")
	}
	r.Players[p.Id] = p
	return nil
}

func (r *Room) NewRoomDeck() {
	r.Deck = InitAllCards()
}

func (r *Room) Run() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case message := <-r.broadcast:
			{
				// broadcast the original message to all others players
				//log.Print("Room broadcast receive message: ", message)
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
					log.Println("Destroy Room since no active player")
					r.unRegister <- r.RoomUUID
					return
				}
			}
		}
	}
}
