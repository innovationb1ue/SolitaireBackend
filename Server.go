package main

import (
	"SolitaireBackend/Decoders"
	"SolitaireBackend/HttpHelper"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type ServerStatusManager struct {
	Players          map[string]*Player
	Rooms            map[string]*Room
	UnregisterRoom   chan string
	UnregisterPlayer chan string
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
	http.HandleFunc("/player/create", s.CreatePlayer)
	http.HandleFunc("/player/socket", s.OpenPlayerSocket)
	http.HandleFunc("/room/create", s.CreateNewRoom)
	http.HandleFunc("/player/join_room", s.JoinRoom)
	// register self services
	go s.closeRoomService()
	go s.deletePlayerService()
	// start server
	_ = http.ListenAndServe(port, nil)
}

func (s *ServerStatusManager) closeRoomService() {
	for {
		select {
		case roomUUID := <-s.UnregisterRoom:
			log.Printf("Unregister room %s at server side", roomUUID)
			delete(s.Rooms, roomUUID)
		}
	}
}

func (s *ServerStatusManager) deletePlayerService() {
	for {
		select {
		case playerId := <-s.UnregisterPlayer:
			log.Printf("Unregister player %s at server side", playerId)
			delete(s.Players, playerId)
		}
	}
}

func (s *ServerStatusManager) CreateNewRoom(w http.ResponseWriter, _ *http.Request) {
	room := NewRoom(s.UnregisterRoom)
	s.Rooms[room.RoomUUID] = room
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(
		map[string]interface{}{
			"message":   "ok create room",
			"status":    0,
			"room_uuid": room.RoomUUID,
		})
	go room.Run()
}

func (s *ServerStatusManager) JoinRoom(w http.ResponseWriter, r *http.Request) {
	// decode frontend data
	ReqJson, err := Decoders.Req2Json(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	PlayerId := ReqJson["player_id"].(string)
	RoomId := ReqJson["room_id"].(string)
	if PlayerId == "" || RoomId == "" {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "Empty necessary parameter received"})
		return
	}
	room := s.Rooms[RoomId]
	p := s.Players[PlayerId]
	if room == nil || p == nil {
		_ = HttpHelper.ReturnJson(w, map[string]interface{}{"message": "player or room not found", "status": -1})
		return
	}
	err = room.AddPlayer(p)
	if err != nil {
		_ = HttpHelper.ReturnJson(w, map[string]interface{}{"message": "Too many players", "status": -1})
		return
	}
	p.room = room
	w.Header().Set("Content-Type", "application/json")
	_ = HttpHelper.ReturnJson(w, map[string]interface{}{"message": "ok join room", "status": 0})
}

func (s *ServerStatusManager) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	// check required variable
	ReqBodyJson, err := Decoders.Req2Json(r.Body)
	if err != nil {
		return
	}
	PlayerName := ReqBodyJson["player_name"].(string)
	if PlayerName == "" {
		_ = HttpHelper.ReturnJson(w, map[string]interface{}{"message": "Failed create player, empty player name. "})
		return
	}
	// create player struct
	p := NewPlayer(PlayerName)
	p.unRegisterSelf = s.UnregisterPlayer
	s.Players[p.Id] = p
	_ = HttpHelper.ReturnJson(w, map[string]interface{}{"message": "ok create player", "player_id": p.Id})
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
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		panic(err)
	}
	player.Conn = c
	player.isConnected = true
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
	// init properties
	s.Players = make(map[string]*Player)
	s.Rooms = make(map[string]*Room)
	s.UnregisterRoom = make(chan string, 1)
}
