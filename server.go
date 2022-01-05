package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

type ServerStatusManager struct {
	Server  interface{}
	Players map[int]*Player
	Rooms   map[int]*Room
}

var upgrader = websocket.Upgrader{}

func (s *ServerStatusManager) GetPlayer(PlayerId int) (*Player, error) {
	for _, p := range s.Players {
		if p.Id == PlayerId {
			return p, nil
		}
	}
	panic("Player not found")
}

func (s *ServerStatusManager) Start(port string) {
	// register handle func
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		p := map[string]interface{}{"message": "Main page"}
		_ = json.NewEncoder(writer).Encode(p)
	})

	http.HandleFunc("/player/socket/", s.OpenPlayerSocket)
	http.HandleFunc("/player/create/", s.CreatePlayer)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}

func (s *ServerStatusManager) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	p := &Player{
		Name:   r.FormValue("player_name"),
		Id:     0,
		Score:  0,
		Desk:   nil,
		Decker: cardDeck{},
		Conn:   &websocket.Conn{},
	}
	// iter through possible ids
	for i := 1; i < 9999; i++ {
		if s.Players[i] == nil {
			p.Id = i
			s.Players[i] = p
			break
		}
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "ok create player", "player_id": p.Id})
}

func (s *ServerStatusManager) OpenPlayerSocket(w http.ResponseWriter, r *http.Request) {
	// get corresponding player
	playerId, _ := strconv.Atoi(r.FormValue("player_id"))
	player := s.Players[playerId]
	// check for multiple connect
	if player.isConnected {
		log.Print("Player already connected")
		w.Header().Set("ErrorMessage", "Multiple connect detected")
		w.WriteHeader(http.StatusConflict)
		return
	}
	// upgrade to websocket
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		panic(err)
	}
	player.Conn = c
	player.isConnected = true
	player.Conn.SetCloseHandler(func(code int, text string) error {
		player.isConnected = false
		_ = player.Conn.WriteJSON(map[string]string{"message": "Successfully disconnected"})
		return nil
	})
	// response
	err = player.Conn.WriteJSON(map[string]string{"message": "ok socket established"})
	if err != nil {
		log.Print("Error returning ok json. ", err)
		return
	}
	// handle inbound message
	// todo: handle all player inbound messages
	// go testHandler(player.Conn)
}

func testHandler(c *websocket.Conn) {
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func (s *ServerStatusManager) Init() {
	s.Players = make(map[int]*Player)
	p := &Player{
		Name:   "Test",
		Id:     0,
		Score:  0,
		Desk:   nil,
		Decker: cardDeck{},
		Conn:   &websocket.Conn{},
	}
	s.Players[0] = p
}
