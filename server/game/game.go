package game

import (
	"cardgame/deck"
	"cardgame/util/slices"

	"fmt"
	"log"
)

func (r *Room) HandleMessage(message ClientMessage) {
	fmt.Printf("room<- %#v\n", message)

	switch m := message.(type) {
	case ClientJoin:
		r.HandleJoin(m)
	case ClientLeave:
		r.HandleLeave(m)
	case ClientChangeDetails:
		r.HandleChangeDetails(m)
	}
}

func (r *Room) HandleJoin(message ClientJoin) {
	p := message.Player
	if r != p.room {
		log.Println("[error] player is in another room")
		p.outbound <- &ServerError{"player is in another room"}
		return
	}

	for _, player := range r.Players {
		if player.Id == p.Id {
			log.Println("[error] player is already in room")
			p.outbound <- &ServerError{"player is already in room"}
			return
		}
	}

	r.Players = append(r.Players, p)

	if len(r.Players) == 1 {
		// first player becomes owner
		r.OwnerId = p.Id
	}

	p.outbound <- &ServerRoomDetails{p.room}
	r.outbound <- &serverPayload{
		exclude: set{p.Id: {}},
		message: &ServerJoin{
			Id:     p.Id,
			Player: *p,
		},
	}

}

func (r *Room) HandleLeave(message ClientLeave) {
	p := message.Player

	slices.Remove(r.Players, p)

	r.outbound <- &serverPayload{
		message: &ServerLeave{
			Id: p.Id,
		},
	}
}

func (r *Room) HandleChangeDetails(message ClientChangeDetails) {
	p := message.Player

	if r.Id != r.OwnerId {
		log.Println("[error] player is not owner")
		p.outbound <- &ServerError{"player is not owner"}
		return
	}

	if message.Name != nil {
		r.Name = *message.Name
	}
	if message.Password != nil {
		r.SetPassword(*message.Password)
	}
	if message.MaxPlayers != nil {
		r.MaxPlayers = *message.MaxPlayers
	}
	if message.HubDeviceId != nil {
		r.HubDeviceId = *message.HubDeviceId
	}
	if message.PlayMode != nil {
		r.PlayMode = *message.PlayMode
	}
	if len(message.AddDecks) > 0 {
		toAdd := []*deck.Deck{}
		for _, deckId := range message.AddDecks {
			if deck, ok := deck.Decks()[deckId]; ok {
				toAdd = append(toAdd, deck)
			}
		}
		r.Decks = slices.Unique(append(r.Decks, toAdd...))
	}
	if len(message.RemoveDecks) > 0 {
		for _, deckId := range message.RemoveDecks {
			if deck, ok := deck.Decks()[deckId]; ok {
				r.Decks = slices.Remove(r.Decks, deck)
			}
		}
	}

}
