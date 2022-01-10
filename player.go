package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

// Player struct
type Player struct {
	Name           string
	Id             string
	Score          int
	Desk           [][]card
	Deck           cardDeck
	Conn           *websocket.Conn
	room           *Room
	isConnected    bool
	ConnExpireTime time.Duration
	ExpireTimer    time.Timer
	// channel of outbound messages
	send chan map[string]interface{}
}

func NewPlayer(Name string, Id string) *Player {
	return &Player{
		Name:           Name,
		Id:             Id,
		Score:          0,
		Desk:           nil,
		Deck:           cardDeck{},
		Conn:           &websocket.Conn{},
		room:           nil,
		isConnected:    false,
		ConnExpireTime: 120 * time.Second,
		send:           make(chan map[string]interface{}, 10),
	}
}

// receive player message and pump it to the server room broadcast
func (p *Player) readPump() {
	// defer close connection
	defer func() { _ = p.Conn.Close() }()
	// read messages
	for {
		_, message, err := p.Conn.ReadMessage()
		// decode bytes stream to Json
		messageJson := map[string]interface{}{}
		_ = json.Unmarshal(message, &messageJson)
		if err != nil {
			log.Printf("error: %v in readPump", err)
			return
		}
		log.Printf("Player readPump Receive message: %s", message)
		p.room.broadcast <- messageJson
	}
}

// emit message to client side
func (p *Player) writePump() {
	defer func() { _ = p.Conn.Close() }()
	ticker := time.NewTicker(60 * time.Second)
	expireTimer := time.NewTimer(p.ConnExpireTime)
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
			MsgPack, err := json.Marshal(msg)
			_, _ = w.Write(MsgPack)
		//Heartbeat
		case msg := <-ticker.C:
			{
				// todo: do something to check player alive. Otherwise do not reset the timer
				fmt.Println(msg)
				expireTimer.Stop()
				expireTimer.Reset(p.ConnExpireTime)
			}
		// on expire
		case _ = <-expireTimer.C:
			{
				p.isConnected = false
				delete(p.room.Players, p.Id)
				p.room = nil
				p.Conn = nil
			}
		}

	}
}
