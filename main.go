package main

import "fmt"

func main() {
	Decker := cardDecker{cards: make([]card, 52)}
	Decker.initAllCards()
	fmt.Println(Decker.cardLeft())
	Decker.removeIndex(2)
	fmt.Println(Decker)
	fmt.Println(Decker.cardLeft())
}
