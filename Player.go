package main

import (
	"SolitaireBackend/Decoders"
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
	Conn           *websocket.Conn
	room           *Room
	isConnected    bool
	ConnExpireTime time.Duration
	LastSeen       time.Time
	send           chan map[string]interface{} // channel of outbound messages
}

func NewPlayer(Name string, Id string) *Player {
	return &Player{
		Name:           Name,
		Id:             Id,
		Score:          0,
		Conn:           nil,
		room:           nil,
		isConnected:    false,
		ConnExpireTime: 10 * time.Second,
		LastSeen:       time.Now(),
		send:           make(chan map[string]interface{}, 10),
	}
}

// receive player message and pump it to the server room broadcast
func (p *Player) readPump() {
	for {
		_, message, err := p.Conn.ReadMessage()
		if err != nil {
			log.Printf("Conn Receive error %v, destroying Player Conn", err)
			p.Destroy()
			break
		}
		messageJson, err := Decoders.Msg2Map(message)
		if err != nil {
			panic(err)
		}
		// switch to handler
		switch messageJson["action"] {
		case "Heartbeat":
			p.LastSeen = time.Now()
		case "InitDeck":
			if p.room.Deck == nil {
				p.room.NewRoomDeck()
			}
			deckBytes, _ := json.Marshal(p.room.Deck)
			p.send <- map[string]interface{}{"action": "InitDeck", "Deck": string(deckBytes)}
		default:
			messageJson["sender"] = p
			p.room.broadcast <- messageJson
		}
	}
}

// emit message to client side
func (p *Player) writePump() {
	ticker := time.NewTicker(3 * time.Second)
	expireTicker := time.NewTicker(p.ConnExpireTime)
	for {
		select {
		// broadcast message to client
		case msg, ok := <-p.send:
			log.Printf("p.send receive %s", msg["action"])
			delete(msg, "sender")
			// if the websocket already closed
			if !ok {
				log.Print("Chan already closed")
				_ = p.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				p.Destroy()
				return
			}
			// write message
			MsgPack, _ := json.Marshal(msg)
			_ = p.Conn.WriteMessage(websocket.TextMessage, MsgPack)
		// require heartbeat
		case _ = <-ticker.C:
			{
				if !p.isConnected {
					return
				}
				// write message
				HeartbeatMessage := map[string]interface{}{"action": "Heartbeat"}
				MsgPack, _ := json.Marshal(HeartbeatMessage)
				_ = p.Conn.WriteMessage(websocket.TextMessage, MsgPack)
			}
		// check player heartbeat
		case t := <-expireTicker.C:
			if (time.Now().Sub(p.LastSeen)) > p.ConnExpireTime {
				log.Print(p.LastSeen)
				log.Print(time.Now().Sub(p.LastSeen))
				log.Print(t, " onDestroy")
				p.Destroy()
				return
			}
		}
	}
}

func (p *Player) Destroy() {
	log.Print("connection expired")
	p.isConnected = false
	delete(p.room.Players, p.Id)
	p.room = nil
	_ = p.Conn.WriteMessage(websocket.CloseMessage, nil)
	_ = p.Conn.Close()
	p.Conn = nil
}
