package main

import (
	"encoding/json"
	"net/http"
)

type ServerStatusManager struct {
	Server  interface{}
	Players map[int]Player
}

func (s *ServerStatusManager) GetPlayer(PlayerId int) (Player, error) {
	for _, p := range s.Players {
		if p.Id == PlayerId {
			return p, nil
		}
	}
	panic("Player not found")
}

func (s ServerStatusManager) Start(port string) {
	// register handle func
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		// get a payload p := Payload{d}
		p := map[string]interface{}{"test": "ok"}
		err := json.NewEncoder(writer).Encode(p)
		if err != nil {
			panic(err)
		}
	})

	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}
