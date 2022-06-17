package deck

import (
	"cardgame/card"
	"cardgame/util"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// YamlDeck is the yaml representation of a deck.
type YamlDeck struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Cards       []string `yaml:"cards"`
	WildCards   []string `yaml:"wild_cards"`
}

// Deck represents a collection of cards and wild cards.
type Deck struct {
	Id          string          `json:"id"`
	Name        string          `json:"name"`
	Location    string          `json:"location"`
	Description string          `json:"description"`
	Cards       []card.Card     `json:"cards"`
	WildCards   []card.WildCard `json:"wild_cards"`
}

var decks = make(map[string]*Deck)

func scanDecks(dir string, prefix string) {
	dirListing, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, file := range dirListing {
		name := file.Name()

		if name[0] == '_' {
			continue
		}

		if file.IsDir() {
			if strings.Contains(name, ".") {
				fmt.Fprintf(os.Stderr, "[deck] Skipping directory %s because it contains a dot.\n", name)
				continue
			}
			scanDecks(dir+"/"+name, prefix+name+".")
			continue
		}

		if !strings.HasSuffix(name, ".yml") {
			continue
		}

		cleanName := strings.TrimSuffix(name, ".yml")

		if strings.Contains(cleanName, ".") {
			fmt.Fprintf(os.Stderr, "[deck] Skipping deck %s because it contains a dot.\n", cleanName)
			continue
		}

		yamlDeck := YamlDeck{}
		contents, err := os.ReadFile(dir + "/" + name)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(contents, &yamlDeck)
		if err != nil {
			panic(err)
		}

		deck := &Deck{
			Name:        yamlDeck.Name,
			Description: yamlDeck.Description,
			Location:    prefix + cleanName,
		}

		deck.Id = util.LongIdFrom("d", fmt.Sprintf("%s%s|%s|%s|%s", prefix, cleanName, deck.Location, deck.Name, deck.Description))

		if deck.Name == "" {
			deck.Name = deck.Location
		}

		if len(yamlDeck.Cards) == 0 {
			panic(fmt.Sprintf("deck %s has no cards", name))
		}

		for _, cardString := range yamlDeck.Cards {
			c, err := card.CardFromString(cardString)
			if err != nil {
				panic(err)
			}

			// deck id is based on the cards in the deck for caching purposes
			// if new cards are added to the deck, the deck id will change
			deck.Id = util.LongIdFrom("d", deck.Id+"-"+c.String())
			deck.Cards = append(deck.Cards, c)
		}

		for _, cardString := range yamlDeck.WildCards {
			w, err := card.WildCardFromString(cardString)
			if err != nil {
				panic(err)
			}

			deck.Id = util.LongIdFrom("d", deck.Id+"-"+w.String())
			deck.WildCards = append(deck.WildCards, w)
		}

		decks[deck.Id] = deck
	}
}

// InitDecks looks for decks in the given directory, recursively, and loads them into memory.
func InitDecks(dir string) map[string]*Deck {
	scanDecks(dir, "")
	return decks
}

// InitDecksOnce is like InitDecks, but it only will run if the decks haven't been loaded yet.
func InitDecksOnce(dir string) map[string]*Deck {
	if len(decks) > 0 {
		return decks
	}

	return InitDecks(dir)
}

// Decks returns the map of decks.
func Decks() map[string]*Deck {
	return decks
}
