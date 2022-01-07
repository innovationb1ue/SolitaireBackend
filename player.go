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

// receive player message and pump it to the server room
func (p *Player) readPump() {
	defer func() { _ = p.Conn.Close() }()
	for {
		_, message, err := p.Conn.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
		}
		p.room.broadcast <- message
	}
}

// emit message to client side
func (p *Player) writePump() {
	defer func() { _ = p.Conn.Close() }()
	ticker := time.NewTicker(60 * time.Second)
	for {
		select {
		case msg, ok := <-p.send:
			// if the chan already closed
			if !ok {
				_ = p.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// create a writer
			w, err := p.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("%v", err)
			}
			// write message
			fmt.Println("Write msg to client: ", msg)
			_, _ = w.Write(msg)
		// Heartbeat
		case msg := <-ticker.C:
			fmt.Println(msg)
		}
	}
}
