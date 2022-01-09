package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestRoom_Run(t *testing.T) {
	t.Log("Inside test room")
	r := newRoom()
	go r.Run()

	s := ServerStatusManager{}
	s.Init()
	go s.Start(":9123")

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
	RoomUUID := (*v)["room_uuid"].(string)

	p1Resp, err := http.Get("http://127.0.0.1:9123/player/create?player_name=Jeff")
	p1Json := DecodeJson(p1Resp)
	p1Id := (*p1Json)["player_id"].(string)
	log.Print(p1Id)

	p2Resp, err := http.Get("http://127.0.0.1:9123/player/create?player_name=Nerd")
	p2Json := DecodeJson(p2Resp)
	p2Id := (*p2Json)["player_id"].(string)
	log.Print(p2Id)

	data1 := url.Values{"player_id": {p1Id}, "room_id": {RoomUUID}}
	joinResp1, _ := http.PostForm("http://127.0.0.1:9123/player/join_room", data1)
	joinJson1 := DecodeJson(joinResp1)
	log.Print(joinJson1)

	data2 := url.Values{"player_id": {p2Id}, "room_id": {RoomUUID}}
	joinResp2, _ := http.PostForm("http://127.0.0.1:9123/player/join_room", data2)
	joinJson2 := DecodeJson(joinResp2)
	log.Print(joinJson2)

	// connect websocket
	c0Conn, _, _ := websocket.DefaultDialer.Dial("ws://127.0.0.1:9123/player/socket/?player_id="+p1Id, nil)
	c1Conn, _, _ := websocket.DefaultDialer.Dial("ws://127.0.0.1:9123/player/socket/?player_id="+p2Id, nil)
	_, t1, _ := c0Conn.ReadMessage()
	log.Printf("%s", t1)
	_ = c0Conn.WriteMessage(websocket.TextMessage, []byte("123"))
	_ = c1Conn.WriteMessage(websocket.TextMessage, []byte("123"))
	go PrintConnMsg(c0Conn)
	go PrintConnMsg(c1Conn)
	time.Sleep(5 * time.Second)
}

func DecodeJson(resp *http.Response) *map[string]interface{} {
	v := &map[string]interface{}{}
	json.NewDecoder(resp.Body).Decode(v)
	return v
}

func PrintConnMsg(conn *websocket.Conn) {
	for {
		_, msg, _ := conn.ReadMessage()
		log.Printf("Conn client receive message %s", msg)
	}
}
