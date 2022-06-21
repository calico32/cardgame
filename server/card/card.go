package card

import (
	"fmt"
	"strings"
)

type BaseCardType int

const (
	BaseCardTypeNormal BaseCardType = iota
	BaseCardTypeWild
)

type (
	BaseCard interface{ BaseCardType() BaseCardType }

	// Card represents a card in the game. It has a type (star, lines, circle, etc.) and a category (mountain range, cell phone brand, etc.).
	Card struct {
		Id       string   `json:"id"`
		Type     CardType `json:"type"`
		Category string   `json:"category"`
	}

	// WildCard represents a wild card in the game. It has two types.
	WildCard struct {
		Id    string     `json:"id"`
		Types []CardType `json:"types"` // sorted tuple of card types
	}
)

func (c *Card) BaseCardType() BaseCardType     { return BaseCardTypeNormal }
func (w *WildCard) BaseCardType() BaseCardType { return BaseCardTypeWild }

func (c *Card) String() string     { return fmt.Sprintf("%s|%s", c.Type, c.Category) }
func (w *WildCard) String() string { return fmt.Sprintf("%s|%s", w.Types[0], w.Types[1]) }

func (cardA *Card) CompatibleWith(cardB *Card, wildCard *WildCard) bool {
	a := cardA.Type
	b := cardB.Type

	if a == b {
		return true
	}

	if wildCard == nil {
		return false
	}

	return (wildCard.Types[0] == a && wildCard.Types[1] == b) || (wildCard.Types[0] == b && wildCard.Types[1] == a)
}

// CardFromString creates a card from a string representation.
// The string must be of the form "typeSymbol|Category", like "=|Cell Phone Brand".
func CardFromString(s string) (Card, error) {
	c := Card{}

	splitIndex := strings.Index(s, "|")
	if splitIndex == -1 {
		return c, fmt.Errorf("no | found in card string %s", s)
	}

	c.Type = TypeFromString(s[:splitIndex])
	if c.Type == Invalid {
		return c, fmt.Errorf("invalid card type in string %s", s)
	}

	c.Category = s[splitIndex+1:]

	c.Id = NextId("c")
	return c, nil
}

// WildCardFromString creates a wild card from a string representation.
// The string must be of the form "typeA|typeB", like "=|≈".
func WildCardFromString(s string) (WildCard, error) {
	w := WildCard{}

	splitIndex := strings.Index(s, "|")
	if splitIndex == -1 {
		return w, fmt.Errorf("no | found in wildcard string %s", s)
	}

	typeA := TypeFromString(s[:splitIndex])
	typeB := TypeFromString(s[splitIndex+1:])

	if typeA == Invalid || typeB == Invalid {
		return w, fmt.Errorf("invalid card type in string %s", s)
	}

	if typeA == typeB {
		return w, fmt.Errorf("wildcard types must be different in string %s", s)
	}

	// set the types in the wildcard, sorted
	if typeA < typeB {
		w.Types = []CardType{typeA, typeB}
	} else {
		w.Types = []CardType{typeB, typeA}
	}

	w.Id = NextId("w")
	return w, nil
}
