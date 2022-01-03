package main

// Player object for a player
type Player struct {
	Name   string
	Id     int
	Score  int
	Desk   [][]card
	Decker cardDecker
}
