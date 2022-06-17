package deck

import (
	"encoding/json"
	"os"
	"testing"
)

func TestDecks(t *testing.T) {
	decks := InitDecks("../data/decks")
	j, err := json.Marshal(decks)
	if err != nil {
		t.Error(err)
		return
	}
	os.WriteFile("./.decks.json", j, 0644)
}
