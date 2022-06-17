package card

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllCardTypes(t *testing.T) {
	types := AllCardTypes()
	if len(types) != int(numCardTypes) {
		t.Errorf("expected %d card types, got %d", numCardTypes, len(types))
	}
}

func TestCardTypeSymbol(t *testing.T) {
	for _, c := range AllCardTypes() {
		if c.String() == "" {
			t.Errorf("card type %d has no symbol", c)
		}
	}
}

func TestCardTypeFromString(t *testing.T) {
	for _, s := range symbols {
		if c := TypeFromString(s); c == Invalid {
			t.Errorf("type from string %s is invalid", s)
		}
	}
}

func TestInvalidCardTypeSymbol(t *testing.T) {
	invalid := CardType(-9999)
	s := invalid.String()
	assert.Equal(t, "?", s)
}
