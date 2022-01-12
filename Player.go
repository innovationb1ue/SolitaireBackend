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
	Desk           [][]card
	Deck           cardDeck
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
		Desk:           nil,
		Deck:           cardDeck{},
		Conn:           &websocket.Conn{},
		room:           nil,
		isConnected:    false,
		ConnExpireTime: 10 * time.Second,
		LastSeen:       time.Now(),
		send:           make(chan map[string]interface{}, 10),
	}
}

// receive player message and pump it to the server room broadcast
func (p *Player) readPump() {
	// read messages
	for {
		_, message, err := p.Conn.ReadMessage()
		if err != nil {
			log.Printf("Conn Receive %v, destroying Player Conn", err)
			p.Destroy()
			break
		}
		messageJson := Decoders.Msg2Map(message)
		// switch to handler
		switch messageJson["action"] {
		case "Heartbeat":
			p.LastSeen = time.Now()
			log.Print("Reset heartbeat timer")
		default:
			messageJson["sender"] = p
			p.room.broadcast <- messageJson
		}
	}
}

// emit message to client side
func (p *Player) writePump() {
	ticker := time.NewTicker(5 * time.Second)
	expireTicker := time.NewTicker(10 * time.Second)
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
				if !p.isConnected {
					break
				}
				// write message
				HeartbeatMessage := map[string]interface{}{"action": "Heartbeat"}
				MsgPack, _ := json.Marshal(HeartbeatMessage)
				_ = p.Conn.WriteMessage(websocket.TextMessage, MsgPack)
			}
		case t := <-expireTicker.C:
			if (time.Now().Sub(p.LastSeen)) > p.ConnExpireTime {
				log.Print(time.Now().Sub(p.LastSeen))
				log.Print(t, "onDestroy")
				p.Destroy()
				break
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
