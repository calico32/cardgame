package card

import "cardgame/util/slices"

// CardType represents a card type, like four dots or circle.
type CardType int

const (
	Lines CardType = iota
	Waves
	Square
	Dots
	Hash
	Circle
	Plus
	Star
	numCardTypes int      = iota
	Invalid      CardType = -1
)

var symbols = [...]string{
	Lines:  "=",
	Waves:  "≈",
	Square: "■",
	Dots:   "⁘",
	Hash:   "♯",
	Circle: "○",
	Plus:   "+",
	Star:   "☆",
}

// String returns symbol for a card type.
func (c CardType) String() (r string) {
	defer func() {
		if recover() != nil {
			r = "?"
		}
	}()
	return symbols[c]
}

// TypeFromString returns a card type from its symbol, or Invalid (-1) if not found.
func TypeFromString(s string) CardType {
	return CardType(slices.IndexOf(symbols[:], s))
}

// CardTypeCount returns the number of card types.
func CardTypeCount() int {
	return numCardTypes
}

// AllCardTypes returns all card types in a slice.
func AllCardTypes() []CardType {
	var types []CardType
	for i := 0; i < int(numCardTypes); i++ {
		types = append(types, CardType(i))
	}
	return types
}
