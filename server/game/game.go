package game

import (
	"cardgame/card"
	"cardgame/deck"
	"cardgame/util/slices"
	"math/rand"
	"time"

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
	case ClientKick:
		r.HandleKick(m)
	case ClientStart:
		r.HandleStart(m)
	case ClientDraw:
		r.HandleDraw(m)
	case ClientSend:
		r.HandleSend(m)
	case ClientChat:
		r.HandleChat(m)
	default:
		fmt.Printf("[error] unhandled message type %T\n", m)
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

	if r.IsFull() {
		p.outbound <- &ServerError{"Room is full"}
		return
	}

	r.Players = append(r.Players, p)

	if len(r.Players) == 1 {
		// first player becomes owner
		r.OwnerId = p.Id
	}

	p.outbound <- &ServerAck{}
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

	if p.Id != r.OwnerId {
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

func (r *Room) HandleKick(message ClientKick) {
	p := message.Player

	if p.Id != r.OwnerId {
		log.Println("[error] player is not owner")
		p.outbound <- &ServerError{"player is not owner"}
		return
	}

	if message.Id == p.Id {
		log.Println("[error] player cannot kick themselves")
		p.outbound <- &ServerError{"player cannot kick themselves"}
		return
	}

	for _, player := range r.Players {
		if player.Id == message.Id {
			player.outbound <- &ServerKick{}
			return
		}
	}
}

func (r *Room) HandleStart(message ClientStart) {
	p := message.Player

	if p.Id != r.OwnerId {
		log.Println("[error] player is not owner")
		p.outbound <- &ServerError{"player is not owner"}
		return
	}

	if r.GamePhase != GamePhaseLobby {
		log.Println("[error] game has already started")
		p.outbound <- &ServerError{"game has already started"}
		return
	}

	r.GamePhase = GamePhasePlaying
	// pick random player to start
	r.CurrentTurn = rand.Intn(len(r.Players))
	r.outbound <- &serverPayload{
		message: &ServerStart{
			CurrentTurn: r.CurrentTurn,
		},
	}
}

func (r *Room) HandleDraw(message ClientDraw) {
	p := message.Player

	if r.GamePhase != GamePhasePlaying {
		log.Println("[error] game is not in playing phase")
		p.outbound <- &ServerError{"game is not in playing phase"}
		return
	}

	if p.Id != r.Players[r.CurrentTurn].Id {
		log.Println("[error] player is not current turn")
		p.outbound <- &ServerError{"player is not current turn"}
		return
	}

	c, err := r.drawCard()
	if err != nil {
		// error occurs when there are no cards left
		// should never happen as we replenish the deck after each draw
		log.Println("[error]", err)
		p.outbound <- &ServerError{err.Error()}
		return
	}

	if wild, ok := c.(*card.WildCard); ok {
		r.usedWildCards = append(r.usedWildCards, r.ActiveWildCard)
		r.ActiveWildCard = wild
		r.outbound <- &serverPayload{
			message: &ServerWildCard{
				PlayerId: p.Id,
				Card:     wild,
			},
		}
	} else {
		r.outbound <- &serverPayload{
			message: &ServerDraw{
				PlayerId: p.Id,
				Card:     c.(*card.Card),
			},
		}

		r.CurrentTurn = (r.CurrentTurn + 1) % len(r.Players)
		r.resync()
		r.outbound <- &serverPayload{
			message: &ServerTurn{
				PlayerId: r.Players[r.CurrentTurn].Id,
			},
		}
	}

	if r.DrawPileSize == 0 {
		// reshuffle
		r.recreateDrawPile()
		for _, player := range r.Players {
			player.outbound <- &ServerReshuffle{
				Player: player,
			}
		}
	}
}

func (r *Room) HandleSend(message ClientSend) {
	p := message.Player

	if r.GamePhase != GamePhasePlaying {
		log.Println("[error] game is not in playing phase")
		p.outbound <- &ServerError{"game is not in playing phase"}
	}

	target := r.getPlayer(message.RecipientId)
	if target == nil {
		log.Println("[error] target player not found")
		p.outbound <- &ServerError{"target player not found"}
	}

	senderTop := p.Hand.top()
	targetTop := target.Hand.top()

	if !senderTop.CompatibleWith(targetTop, r.ActiveWildCard) {
		log.Println("[error] cards are not compatible")
		p.outbound <- &ServerError{"cards are not compatible"}
	}

	p.Hand = p.Hand.tail()
	target.Hand = append(target.Hand, senderTop)
	target.Score++

	p.outbound <- &ServerSend{
		SenderId:    p.Id,
		RecipientId: target.Id,
		Card:        senderTop,
	}

	target.outbound <- &ServerSend{
		SenderId:    p.Id,
		RecipientId: target.Id,
		Card:        senderTop,
	}

	r.resync()
}

func (r *Room) HandleChat(message ClientChat) {
	if message.Message == "" {
		return
	}

	if message.RecipientId != nil {
		recipient := r.getPlayer(*message.RecipientId)
		if recipient == nil {
			message.Player.outbound <- &ServerError{"player not found"}
			return
		}
		recipient.outbound <- &ServerChat{
			Timestamp: fmt.Sprint(time.Now().UnixMilli()),
			PlayerId:  message.Player.Id,
			Message:   message.Message,
			Private:   true,
		}
		return
	}

	r.outbound <- &serverPayload{
		message: &ServerChat{
			Timestamp: fmt.Sprint(time.Now().UnixMilli()),
			PlayerId:  message.Player.Id,
			Message:   message.Message,
			Private:   false,
		},
	}
}
