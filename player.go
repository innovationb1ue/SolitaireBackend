package main

import "github.com/gorilla/websocket"

// Player struct
type Player struct {
	Name        string
	Id          int
	Score       int
	Desk        [][]card
	Decker      cardDeck
	Conn        *websocket.Conn
	isConnected bool
	send        chan []byte
}

func (p Player) JoinRoom(RoomId int) {
	
}
