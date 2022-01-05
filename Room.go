package main

type Room struct {
	RoomId    int
	Players   map[int]*Player
	broadcast chan []byte
}

func newRoom(RoomId int) *Room {
	return &Room{
		RoomId:    RoomId,
		Players:   make(map[int]*Player),
		broadcast: make(chan []byte),
	}
}

func (r *Room) AddPlayer(p *Player) {
	r.Players[p.Id] = p
}

func (r *Room) run() {
	for {
		select {
		case message := <-r.broadcast:
			{
				for _, client := range r.Players {
					client.send <- message
				}
			}
		}
	}
}
