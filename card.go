package main

// a single card

type card struct {
	status int // status: 0-> ready, 1-10 -> position on the pile, 11 -> deleted
	num    int // card number
	suit   int // card suit. 1-spade  2-clubs 3-heart 4-diamond
}

// collection of undistributed cards
type cardDeck struct {
	cards []card
}

func (cDecker *cardDeck) cardLeft() int {
	count := 0
	for _, c := range cDecker.cards {
		if c.status != -1 {
			count++
		}
	}
	return count
}

func (cDecker *cardDeck) addCard(c card) {
	cDecker.cards = append(cDecker.cards, c)
}

func (cDecker *cardDeck) removeCard(c card) {
	for i, t := range cDecker.cards {
		if c == t {
			cDecker.removeIndex(i)
		}
	}
}

func (cDecker *cardDeck) removeIndex(index int) {
	cDecker.cards[index] = card{
		status: -1,
		num:    -1,
		suit:   -1,
	}
}

func (cDecker *cardDeck) removeById(id int) {
	for i, t := range cDecker.cards {
		if t.status == id {
			cDecker.removeIndex(i)
		}
	}
}

func (cDecker *cardDeck) initAllCards() {
	count := 0
	for i := 1; i <= 13; i++ {
		for j := 0; j <= 3; j++ {
			cDecker.cards[count] = card{
				status: 0,
				num:    i,
				suit:   j,
			}
			count++
		}
	}
}
