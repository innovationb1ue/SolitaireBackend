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
	r := newRoom(make(chan string))
	go r.Run()

	s := ServerStatusManager{}
	s.Init()
	go s.Start(":9123")

	// create a room
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

	// create 2 players
	p1Resp, err := http.PostForm("http://127.0.0.1:9123/player/create",
		url.Values{
			"player_name": {
				"Jeff",
			},
		})
	p1Json := DecodeJson(p1Resp)
	p1Id := (*p1Json)["player_id"].(string)
	log.Print(p1Id)
	p2Resp, err := http.PostForm("http://127.0.0.1:9123/player/create",
		url.Values{
			"player_name": {
				"Nerd",
			},
		})
	p2Json := DecodeJson(p2Resp)
	p2Id := (*p2Json)["player_id"].(string)
	log.Print(p2Id)
	// join the same room
	data1 := url.Values{"player_id": {p1Id}, "room_id": {RoomUUID}}
	joinResp1, _ := http.PostForm("http://127.0.0.1:9123/playerjoin_room", data1)
	joinJson1 := DecodeJson(joinResp1)
	log.Print(joinJson1)
	data2 := url.Values{"player_id": {p2Id}, "room_id": {RoomUUID}}
	joinResp2, _ := http.PostForm("http://127.0.0.1:9123/playerjoin_room", data2)
	joinJson2 := DecodeJson(joinResp2)
	log.Print(joinJson2)
	// connect websocket
	c0Conn, _, _ := websocket.DefaultDialer.Dial("ws://127.0.0.1:9123/player/socket/?player_id="+p1Id, nil)
	c1Conn, _, _ := websocket.DefaultDialer.Dial("ws://127.0.0.1:9123/player/socket/?player_id="+p2Id, nil)
	_, t1, _ := c0Conn.ReadMessage()
	log.Printf("%s", t1)
	// write some message
	_ = c0Conn.WriteJSON(map[string]interface{}{"message": "123"})
	_ = c1Conn.WriteJSON(map[string]interface{}{"message": "123"})
	// receive broadcast
	go PrintConnMsg(c0Conn)
	go PrintConnMsg(c1Conn)
	// wait for goroutines
	time.Sleep(5 * time.Second)
}

func DecodeJson(resp *http.Response) *map[string]interface{} {
	v := &map[string]interface{}{}
	_ = json.NewDecoder(resp.Body).Decode(v)
	return v
}

func PrintConnMsg(conn *websocket.Conn) {
	for {
		_, msg, _ := conn.ReadMessage()
		log.Printf("Conn client receive message %s", msg)
	}
}
