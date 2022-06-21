package game

import (
	"cardgame/card"
	"cardgame/util/slices"
)

type (
	serverPayload struct {
		include set
		exclude set
		message ServerMessage
	}

	ServerMessage interface{ ServerType() string }

	// ServerChangeDetails is sent to all players when the room details change.
	ServerChangeDetails struct {
		Name        *string   `json:"name"`
		Description *string   `json:"description"`
		MaxPlayers  *int      `json:"maxPlayers"`
		DeckIds     []string  `json:"decks"`
		PlayMode    *PlayMode `json:"playMode"`
		HubDeviceId *string   `json:"hubDeviceId"`
	}
	// ServerJoin is sent to all players when a new player joins the room.
	ServerJoin struct {
		Id     string `json:"id"`
		Player Player `json:"player"`
	}
	// ServerAck is sent to a player when they join the room.
	ServerAck struct {
	}
	// ServerLeave is sent to all players when a player leaves the room.
	ServerLeave struct {
		Id string `json:"id"`
	}
	// ServerKick is sent to a player when they are kicked from the room.
	ServerKick struct {
	}
	// ServerStart is sent to all players when the game starts.
	ServerStart struct {
		CurrentTurn int `json:"currentTurn"`
	}
	// ServerDraw is sent to all players when a player draws a card.
	ServerDraw struct {
		PlayerId string     `json:"playerId"`
		Card     *card.Card `json:"card"`
	}
	// ServerWildCard is sent to all players when a player draws a wild card.
	ServerWildCard struct {
		PlayerId string         `json:"playerId"`
		Card     *card.WildCard `json:"card"`
	}
	// ServerReshuffle is sent to a player when the deck is reshuffled.
	// This event is sent individually to each player to update their own deck.
	ServerReshuffle struct {
		Player *Player
	}
	// ServerSend is sent to the sender and the recipient when a card is sent.
	ServerSend struct {
		SenderId    string     `json:"senderId"`
		RecipientId string     `json:"recipientId"`
		Card        *card.Card `json:"card"`
	}
	// ServerChat is sent to all players when a chat message is sent.
	ServerChat struct {
		Timestamp string `json:"timestamp"`
		PlayerId  string `json:"player"`
		Private   bool   `json:"private"`
		Message   string `json:"message"`
	}
	// ServerResync is sent to all players occasionally to resync the top cards of the hands.
	ServerResync struct {
		TopCards map[string]*card.Card `json:"topCards"` // playerId -> card
	}
	// ServerTurn is sent to all players when a player's turn begins.
	ServerTurn struct {
		PlayerId string `json:"playerId"`
	}
	// ServerError is sent to a player when an error occurs.
	ServerError struct {
		Message string `json:"message"`
	}
)

func (s ServerChangeDetails) ServerType() string { return "change_details" }
func (s ServerJoin) ServerType() string          { return "join" }
func (s ServerAck) ServerType() string           { return "ack" }
func (s ServerLeave) ServerType() string         { return "leave" }
func (s ServerKick) ServerType() string          { return "kick" }
func (s ServerStart) ServerType() string         { return "start" }
func (s ServerDraw) ServerType() string          { return "draw" }
func (s ServerWildCard) ServerType() string      { return "wild_card" }
func (s ServerReshuffle) ServerType() string     { return "reshuffle" }
func (s ServerSend) ServerType() string          { return "send" }
func (s ServerChat) ServerType() string          { return "chat" }
func (s ServerResync) ServerType() string        { return "resync" }
func (s ServerTurn) ServerType() string          { return "turn" }
func (s ServerError) ServerType() string         { return "error" }

var ServerMessageTypes = slices.AssociateReverseBy([]ServerMessage{
	ServerChangeDetails{},
	ServerJoin{},
	ServerAck{},
	ServerLeave{},
	ServerKick{},
	ServerStart{},
	ServerDraw{},
	ServerWildCard{},
	ServerReshuffle{},
	ServerChat{},
	ServerResync{},
	ServerTurn{},
	ServerError{},
}, func(t ServerMessage) string { return t.ServerType() })
