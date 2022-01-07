package main

import (
	"github.com/gorilla/websocket"
)

// Player struct
type Player struct {
	Name        string
	Id          string
	Score       int
	Desk        [][]card
	Decker      cardDeck
	Conn        *websocket.Conn
	isConnected bool
	send        chan []byte
}
