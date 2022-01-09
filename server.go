package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

type ServerStatusManager struct {
	Players map[string]*Player
	Rooms   map[string]*Room
}

var upgrader = websocket.Upgrader{}

func (s *ServerStatusManager) GetPlayer(PlayerId string) (*Player, error) {
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
	http.HandleFunc("/player/create/", s.CreatePlayer)
	http.HandleFunc("/player/socket/", s.OpenPlayerSocket)
	http.HandleFunc("/player/join_room", s.JoinRoom)
	http.HandleFunc("/room/create", s.CreateNewRoom)
	http.HandleFunc("/player/query_all", s.QueryAllPlayer)
	// start server
	_ = http.ListenAndServe(port, nil)
}

func (s *ServerStatusManager) QueryAllPlayer(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "ok create room", "player_ids": s.Players})
}

func (s *ServerStatusManager) CreateNewRoom(w http.ResponseWriter, _ *http.Request) {
	room := newRoom()
	s.Rooms[room.RoomUUID] = room
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "ok create room", "room_uuid": room.RoomUUID})
	go room.Run()
}

func (s *ServerStatusManager) JoinRoom(w http.ResponseWriter, r *http.Request) {
	PlayerId := r.FormValue("player_id")
	RoomId := r.FormValue("room_id")
	room := s.Rooms[RoomId]
	p := s.Players[PlayerId]
	room.AddPlayer(p)
	p.room = room
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "ok join room"})
}

func (s *ServerStatusManager) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	// check required variable
	PlayerName := r.FormValue("player_name")
	if PlayerName == "" {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "Failed create player, empty player name. "})
		return
	}
	// create player struct
	p := NewPlayer(r.FormValue("player_name"), uuid.NewString())
	s.Players[p.Id] = p
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "ok create player", "player_id": p.Id})
}

func (s *ServerStatusManager) OpenPlayerSocket(w http.ResponseWriter, r *http.Request) {
	// get corresponding player
	playerId := r.FormValue("player_id")
	player := s.Players[playerId]
	if player == nil {
		w.Header().Set("ErrorMessage", "Player not found")
		w.WriteHeader(http.StatusConflict)
		return
	}
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
	log.Printf("Player %s Websocket connected", player.Name)
	// call player socket service
	go player.readPump()
	go player.writePump()
}

func (s *ServerStatusManager) Init() {
	// initiate 2 players for test
	s.Players = make(map[string]*Player)
	s.Rooms = make(map[string]*Room)
	for a := 0; a <= 1; a++ {
		p := NewPlayer("Test", strconv.Itoa(a))
		s.Players[p.Id] = p
		for k, v := range s.Players {
			log.Println(k, "===", v)
		}
	}
}
