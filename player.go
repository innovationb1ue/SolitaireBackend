package main

import (
	"encoding/json"
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
	ExpireTimer    *time.Timer
	// channel of outbound messages
	send chan map[string]interface{}
}

func NewPlayer(Name string, Id string) *Player {
	timer := time.NewTimer(120 * time.Second)
	timer.Stop()
	return &Player{
		Name:           Name,
		Id:             Id,
		Score:          0,
		Desk:           nil,
		Deck:           cardDeck{},
		Conn:           &websocket.Conn{},
		room:           nil,
		isConnected:    false,
		ConnExpireTime: 10 * time.Second,
		send:           make(chan map[string]interface{}, 10),
		ExpireTimer:    timer,
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
		// if receive heartbeat => reset the timer
		if messageJson["action"] == "Heartbeat" {
			log.Print("Reset heartbeat timer")
			p.ExpireTimer.Stop()
			p.ExpireTimer.Reset(p.ConnExpireTime)
			return
		}
		messageJson["sender"] = p
		log.Printf("Player readPump Receive message: %s", message)
		p.room.broadcast <- messageJson
	}
}

// emit message to client side
func (p *Player) writePump() {
	ticker := time.NewTicker(5 * time.Second)
	expireTimer := time.NewTimer(p.ConnExpireTime)
	for {
		select {
		case msg, ok := <-p.send:
			log.Printf("p.send receive %s", msg)
			msg = map[string]interface{}{"card_left": msg["card_left"]}
			//if the chan already closed
			if !ok {
				_ = p.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Print("Chan already closed")
				continue
			}
			// write message
			MsgPack, _ := json.Marshal(msg)
			_ = p.Conn.WriteMessage(websocket.TextMessage, MsgPack)
		// require client for heartbeat
		case _ = <-ticker.C:
			{
				HeartbeatMessage := map[string]interface{}{"action": "Heartbeat"}
				//create a writer
				// write message
				MsgPack, err := json.Marshal(HeartbeatMessage)
				if err != nil {
					log.Print("can not marshal json in Heartbeat message")
				}
				_ = p.Conn.WriteMessage(websocket.TextMessage, MsgPack)
			}
		// on expire
		case _ = <-expireTimer.C:
			{
				log.Print("connection expired")
				p.isConnected = false
				delete(p.room.Players, p.Id)
				p.room = nil
				_ = p.Conn.WriteMessage(websocket.CloseMessage, nil)
				_ = p.Conn.Close()
				p.Conn = nil
				return
			}
		}

	}
}
