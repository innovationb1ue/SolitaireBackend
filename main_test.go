package main

import (
	"encoding/json"
	"log"
	"net/http"
	"testing"
)

func TestRoom_Run(t *testing.T) {
	t.Log("Inside test room")
	r := newRoom()
	go r.Run()

	s := ServerStatusManager{}
	s.Init()
	go s.Start(":9123")
	//c0Conn, _, _ := websocket.DefaultDialer.Dial("ws://127.0.0.1:9123/player/socket/?player_id=0", nil)
	//c1Conn, _, _ := websocket.DefaultDialer.Dial("ws://127.0.0.1:9123/player/socket/?player_id=1", nil)
	resp, err := http.Get("http://127.0.0.1:9123/room/create")
	if err != nil {
		panic(err)
	}
	v := &map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		log.Fatal("Failed to decode json response")
	}
	log.Print("resp json unmarshal = ", v)
	//c0Conn.WriteMessage(websocket.TextMessage, []byte("123"))
	//c1Conn.WriteMessage(websocket.TextMessage, []byte("123"))
}
