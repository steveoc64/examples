package solitaire

import (
	"github.com/fyne-io/fyne"
	"log"
)
import "github.com/fyne-io/examples/solitaire/faces"

// Suit encodes one of the four possible suits for a playing card
type Suit int

const (
	// SuitClubs is the "Clubs" playing card suit
	SuitClubs Suit = iota
	// SuitDiamonds is the "Diamonds" playing card suit
	SuitDiamonds
	// SuitHearts is the "Hearts" playing card suit
	SuitHearts
	// SuitSpades is the "Spades" playing card suit
	SuitSpades
)

const (
	// ValueJack is a convenience for the card 1 higher than 10
	ValueJack = 11
	// ValueQueen is the value for a queen face card
	ValueQueen = 12
	// ValueKing is the value for a king face card
	ValueKing = 13
)

// Card is a single playing card, it has a face value and a suit associated with it.
type Card struct {
	Value int
	Suit  Suit
}

// Face returns a resource that can be used to render the associated card
func (c Card) Face() fyne.Resource {
	return faces.ForCard(c.Value, int(c.Suit))
}

// NewCard returns a new card instance with the specified suit and value (1 based for Ace, 2 is 2 and so on).
func NewCard(value int, suit Suit) *Card {
	if value < 1 || value > 13 {
		log.Fatal("Invalid card face value")
	}

	return &Card{value, suit}
}
