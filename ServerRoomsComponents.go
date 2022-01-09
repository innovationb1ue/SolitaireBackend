package main

func (s ServerStatusManager) DestroyRoom(room *Room) {
	delete(s.Rooms, room.RoomUUID)
	room = nil
}
