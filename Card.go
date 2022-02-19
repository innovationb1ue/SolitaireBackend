package main

import (
	"math/rand"
	"time"
)

var ranks = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
var ranksNum = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13"}
var suits = []string{"heart", "diamond", "spades", "clubs"}

//var symbols = map[string]string{"heart": "♥", "clubs": "♣", "spades": "♠", "diamond": "♦"}

func InitAllCards() [][]map[string]interface{} {
	var initDeck []map[string]interface{}
	for i, r := range ranks {
		for _, s := range suits {
			if s == "clubs" || s == "diamond" {
				continue
			}
			initDeck = append(initDeck, map[string]interface{}{"rank": r, "isDown": true, "suit": s,
				"deck": 1, "rank_num": ranksNum[i]})
			initDeck = append(initDeck, map[string]interface{}{"rank": r, "isDown": true, "suit": s,
				"deck": 2, "rank_num": ranksNum[i]})
		}
	}
	// get a rand
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(initDeck), func(i, j int) {
		initDeck[i], initDeck[j] = initDeck[j], initDeck[i]
	})
	// split into 10 piles and a deck
	chunkedDeck := chunkBy(initDeck, 5, 10)
	for _, d := range chunkedDeck {
		d[len(d)-1]["isDown"] = false
	}
	return chunkedDeck
}

func chunkBy(items []map[string]interface{}, chunkSize int, groupCount int) (chunks [][]map[string]interface{}) {
	for chunkSize < len(items) && groupCount > 0 {
		groupCount--
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}
