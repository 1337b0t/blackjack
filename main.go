package main

import (
	"fmt"
	"strings"

	"github.com/1337b0t/deck"
)

type Hand []deck.Card

func (h Hand) String() string {
	strs := make([]string, len(h))
	for i := range h {
		strs[i] = h[i].String()
	}
	return strings.Join(strs, ", ")
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

func (h Hand) DealerString() string {
	return h[0].String() + ", **HIDDEN**"
}

func (h Hand) MinScore() int {
	score := 0

	for _, c := range h {
		score += min(int(c.Rank), 10)
	}

	return score
}

func (h Hand) Score() int {
	minScore := h.MinScore()
	if minScore > 11 {
		return minScore
	}
	for _, c := range h {
		if c.Rank == deck.Ace {
			//Ace already 10
			return minScore + 10
		}
	}
	return minScore
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

type GameState struct {
	Deck   []deck.Card
	State  State
	Player Hand
	Dealer Hand
}

type State uint8

const (
	StatePlayerTurn State = iota
	StateDealerTurn
	StateHandOver
)

func (gs *GameState) CurrentPlayer() *Hand {
	switch gs.State {
	case StatePlayerTurn:
		return &gs.Player
	case StateDealerTurn:
		return &gs.Dealer
	default:
		panic("player state does not exist")
	}
}

func clone(gs GameState) GameState {
	ret := GameState{
		Deck:   make([]deck.Card, len(gs.Deck)),
		State:  gs.State,
		Player: make(Hand, len(gs.Player)),
		Dealer: make(Hand, len(gs.Dealer)),
	}
	copy(ret.Deck, gs.Deck)
	copy(ret.Player, gs.Deck)
	copy(ret.Dealer, gs.Dealer)
	return ret
}

func Shuffle(gs GameState) GameState {
	ret := clone(gs)
	ret.Deck = deck.New(deck.Deck(3), deck.Shuffle)
	return ret
}

func Deal(gs GameState) GameState {
	ret := clone(gs)
	ret.Player = make(Hand, 0, 5)
	ret.Dealer = make(Hand, 0, 5)
	var card deck.Card
	for i := 0; i < 2; i++ {
		card, ret.Deck = draw(ret.Deck)
		ret.Player = append(ret.Player, card)
		card, ret.Deck = draw(ret.Deck)
		ret.Dealer = append(ret.Dealer, card)

	}

	ret.State = StatePlayerTurn

	return ret

}

func Hit(gs GameState) GameState {
	ret := clone(gs)
	hand := ret.CurrentPlayer()
	var card deck.Card
	card, ret.Deck = draw(ret.Deck)
	*hand = append(*hand, card)
	if hand.Score() > 21 {
		return Stand(gs)
	}
	return ret
}

func Stand(gs GameState) GameState {
	ret := clone(gs)
	ret.State++
	return ret
}

func EndHand(gs GameState) GameState {
	ret := clone(gs)
	pScore, dScore := gs.Player.Score(), gs.Dealer.Score()
	fmt.Println("==FINAL HANDS==")
	fmt.Println("Player:", gs.Player, "\nScore:", pScore)
	fmt.Println("Dealer", gs.Dealer, "\nScore:", dScore)

	switch {
	case pScore > 21:
		fmt.Println("busted")
	case dScore > 21:
		fmt.Println("dealer busted")
	case pScore > dScore:
		fmt.Println("you win!")
	case dScore > pScore:
		fmt.Println("you lose!")
	case dScore == pScore:
		fmt.Println("draw")
	}
	fmt.Println()
	ret.Player = nil
	ret.Dealer = nil
	return ret
}

func main() {
	var gs GameState
	gs = Shuffle(gs)
	gs = Deal(gs)

	var input string
	for gs.State == StatePlayerTurn {
		fmt.Println("Player: ", gs.Player)
		fmt.Println("Dealer: ", gs.Dealer.DealerString())
		fmt.Println("What will you do? (h)it, (s)stand")
		fmt.Scanf("%s\n", &input)
		switch input {
		case "h":
			gs = Hit(gs)
		case "s":
			gs = Stand(gs)
		default:
			fmt.Println("Invalid option:", input)
		}
	}

	// if dealer score <=16 we hit
	// if dealer has soft 17, then we hit

	for gs.State == StateDealerTurn {
		if gs.Dealer.Score() <= 16 || (gs.Dealer.Score() == 17 && gs.Dealer.MinScore() != 17) {
			gs = Hit(gs)
		} else {
			gs = Stand(gs)
		}

		EndHand(gs)
	}

}
