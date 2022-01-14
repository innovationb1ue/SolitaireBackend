package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestInitAllCards(t *testing.T) {
	res := InitAllCards()
	if len(res) != 11 {
		t.FailNow()
	}
	JsonBytes, _ := json.Marshal(res)
	JsonString := string(JsonBytes)
	fmt.Println(JsonString)
}
