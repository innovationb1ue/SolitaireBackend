package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

// Player struct
type Player struct {
	Name        string
	Id          string
	Score       int
	Desk        [][]card
	Deck        cardDeck
	Conn        *websocket.Conn
	room        *Room
	isConnected bool
	// channel of outbound messages
	send chan []byte
}

func NewPlayer(Name string, Id string) *Player {
	return &Player{
		Name:        Name,
		Id:          Id,
		Score:       0,
		Desk:        nil,
		Deck:        cardDeck{},
		Conn:        &websocket.Conn{},
		room:        nil,
		isConnected: false,
		send:        make(chan []byte, 10),
	}
}

// receive player message and pump it to the server room
func (p *Player) readPump() {
	// defer close connection
	defer func() { _ = p.Conn.Close() }()
	// read messages
	for {
		_, message, err := p.Conn.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
		}
		log.Printf("Player readPump Receive message %s", message)
		p.room.broadcast <- message
	}
}

// emit message to client side
func (p *Player) writePump() {
	defer func() { _ = p.Conn.Close() }()
	ticker := time.NewTicker(60 * time.Second)
	log.Printf("WritePump active")
	for {
		select {
		case msg, ok := <-p.send:
			log.Printf("p.send receive %s", msg)
			//if the chan already closed
			if !ok {
				_ = p.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Print("Chan already closed")
				return
			}
			//create a writer
			w, err := p.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("%v", err)
				panic(err)
			}
			// write message
			_, _ = w.Write(msg)
		//Heartbeat
		case msg := <-ticker.C:
			fmt.Println(msg)
		}
	}
}
