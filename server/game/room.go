package game

import (
	"cardgame/card"
	"cardgame/deck"
	"cardgame/util/slices"
	"fmt"
)

// Room represents a game room.
type Room struct {
	Id             string           `json:"id"`             // internal room id
	Timstamp       int64            `json:"timestamp"`      // creation timestamp
	Name           string           `json:"name"`           // public-facing room name
	Description    string           `json:"description"`    // room description
	MaxPlayers     int              `json:"maxPlayers"`     // maximum number of players
	OwnerId        string           `json:"ownerId"`        // owner's player id
	Players        []*Player        `json:"players"`        // players in the room, including the owner
	Decks          []*deck.Deck     `json:"decks"`          // decks in use
	PlayMode       PlayMode         `json:"playMode"`       // play mode
	HubDeviceId    string           `json:"hubDeviceId"`    // hub device id
	CurrentTurn    int              `json:"currentTurn"`    // index of the current player
	GamePhase      GamePhase        `json:"gamePhase"`      // game phase
	ActiveWildCard *card.WildCard   `json:"activeWildCard"` // active wild card
	DrawPileSize   int              `json:"drawPileSize"`   // size of the draw pile
	drawPile       []card.BaseCard  // draw pile
	usedWildCards  []*card.WildCard // already used wild cards

	private      bool   // true if the room is private
	passwordHash string // password hash for private rooms

	hub      *Hub                // hub instance
	inbound  chan ClientMessage  // incoming client messages
	outbound chan *serverPayload // outgoing server messages
}

func (r *Room) getPlayer(id string) *Player {
	for _, p := range r.Players {
		if p.Id == id {
			return p
		}
	}
	return nil
}

func (r *Room) createDrawPile() {
	r.drawPile = []card.BaseCard{}
	for _, d := range r.Decks {
		for _, c := range d.Cards {
			r.drawPile = append(r.drawPile, c)
		}
		for _, w := range d.WildCards {
			r.drawPile = append(r.drawPile, w)
		}
	}
	slices.Shuffle(r.drawPile)
	r.DrawPileSize = len(r.drawPile)
}

func (r *Room) drawCard() (card.BaseCard, error) {
	if len(r.drawPile) == 0 {
		return nil, fmt.Errorf("draw pile is empty")
	}
	c := r.drawPile[0]
	r.drawPile = r.drawPile[1:]
	r.DrawPileSize--
	return c, nil
}

func (r *Room) recreateDrawPile() error {
	if r.DrawPileSize != 0 {
		return fmt.Errorf("draw pile is not empty")
	}

	newDrawPile := []card.BaseCard{}
	for _, p := range r.Players {
		// take every card except the top one

		top := p.Hand.top()
		if top == nil {
			// no cards
			continue
		}
		reusable := p.Hand.tail()
		for _, c := range reusable {
			newDrawPile = append(newDrawPile, c)
		}
		p.Hand = make(PlayerHand, 1)
		p.Hand[0] = top
	}
	slices.Shuffle(newDrawPile)
	r.drawPile = newDrawPile

	// choose new wild card
	r.usedWildCards = append(r.usedWildCards, r.ActiveWildCard)
	slices.Shuffle(r.usedWildCards)
	r.ActiveWildCard, r.usedWildCards = r.usedWildCards[0], r.usedWildCards[1:]

	return nil
}

func (r *Room) resync() {
	topCards := make(map[string]*card.Card)

	for _, p := range r.Players {
		topCards[p.Id] = p.Hand.top()
	}

	r.outbound <- &serverPayload{
		message: &ServerResync{
			TopCards: topCards,
		},
	}

}

// IsPrivate returns true if the room is private.
func (r *Room) IsPrivate() bool {
	return r.private
}

// CheckPassword returns true if the password is correct.
func (r *Room) CheckPassword(password string) bool {
	return checkPassword(r.passwordHash, password)
}

// IsFull returns true if the room is full.
func (r *Room) IsFull() bool {
	return len(r.Players) >= r.MaxPlayers
}

// SetPassword sets the password and privacy of the room.
// If the password is empty, the room is made public.
func (r *Room) SetPassword(password string) {
	if password == "" {
		r.private = false
		r.passwordHash = ""
	} else {
		r.private = true
		r.passwordHash = hashPassword(password)
	}
}

func (r *Room) read() {
	for {
		message, ok := <-r.inbound
		if !ok {
			fmt.Println("room<- closed")
			// room closed
			return
		}
		go r.HandleMessage(message)
	}
}

func (r *Room) write() {
	for {
		payload, ok := <-r.outbound
		if !ok {
			fmt.Println("room-> closed")
			// room closed
			return
		}

		included := []*Player{}
		other := []*Player{}

		for _, p := range r.Players {
			if _, ok := payload.include[p.Id]; ok {
				included = append(included, p)
				continue
			}
			if _, ok := payload.exclude[p.Id]; ok {
				continue
			}

			other = append(other, p)
		}

		var toSend []*Player
		if len(included) > 0 {
			toSend = included
		} else if len(other) > 0 {
			toSend = other
		} else {
			continue
		}

		fmt.Println("room->", payload.message)
		for _, p := range toSend {
			fmt.Println("room->    sending to", p.Id)
			p.outbound <- payload.message
		}
	}
}
