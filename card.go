package main

// for a single card
type card struct {
	id   int
	num  int
	suit int
}

// collection of undistributed cards
type cardDecker struct {
	cards []card
}

// count how may cards in card collection
func (cDecker *cardDecker) cardLeft() int {
	count := 0
	for _, c := range cDecker.cards {
		if c.id != -1 {
			count++
		}
	}
	return count
}

func (cDecker *cardDecker) addCard(c card) {
	cDecker.cards = append(cDecker.cards, c)
}

func (cDecker *cardDecker) removeCard(c card) {
	for i, t := range cDecker.cards {
		if c == t {
			cDecker.removeIndex(i)
		}
	}
}

func (cDecker *cardDecker) removeIndex(index int) {
	cDecker.cards[0] = card{
		id:   -1,
		num:  -1,
		suit: -1,
	}
}

func (cDecker *cardDecker) removeById(id int) {
	for i, t := range cDecker.cards {
		if t.id == id {
			cDecker.removeIndex(i)
		}
	}
}

func (cDecker *cardDecker) initAllCards() {
	count := 0
	for i := 1; i <= 13; i++ {
		for j := 0; j <= 3; j++ {
			cDecker.cards[count] = card{
				id:   count,
				num:  i,
				suit: j,
			}
			count++
		}
	}
}
