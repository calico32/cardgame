package game

import (
	"cardgame/util/slices"
	"encoding/json"
	"errors"
	"log"
	"reflect"
)

type set map[string]struct{}

type (
	clientPayload struct {
		Type string `json:"type"`
	}

	ClientMessage interface{ ClientType() string }

	// ClientChangeDetails is sent by the room owner to change the room details and add/remove decks
	ClientChangeDetails struct {
		Player *Player `json:"-"`

		Name        *string   `json:"name"`
		Description *string   `json:"description"`
		MaxPlayers  *int      `json:"maxPlayers"`
		Password    *string   `json:"password"`    // new password for private rooms, or "" for public rooms
		AddDecks    []string  `json:"addDecks"`    // IDs of decks to add
		RemoveDecks []string  `json:"removeDecks"` // IDs of decks to remove
		PlayMode    *PlayMode `json:"playMode"`    // new play mode
		HubDeviceId *string   `json:"hubDeviceId"` // ID of the hub device to use
	}
	// ClientJoin is sent by a new player joining the room.
	ClientJoin struct {
		Player *Player `json:"-"`

		Name   string       `json:"name"`
		Avatar AvatarConfig `json:"avatar"`
	}
	// ClientLeave is sent by a player leaving the room.
	ClientLeave struct {
		Player *Player `json:"-"`
	}
	// ClientKick is sent by the room owner to kick a player.
	ClientKick struct {
		Player *Player `json:"-"`

		Id string `json:"id"`
	}
	// ClientStart is sent by the room owner to start the game.
	ClientStart struct {
		Player *Player `json:"-"`
	}
	// ClientDraw is sent by a player to draw a card.
	ClientDraw struct {
		Player *Player `json:"-"`
	}
	// ClientSend is sent by a player to send a card to another player.
	ClientSend struct {
		Player *Player `json:"-"`

		RecipientId string `json:"recipientId"`
	}
	// ClientChat is sent by a player to send a chat message.
	ClientChat struct {
		Player *Player `json:"-"`

		Message     string  `json:"message"`
		RecipientId *string `json:"recipient"` // RecipientId is set if the message is a private message.
	}
)

func (c ClientChangeDetails) ClientType() string { return "change_details" }
func (c ClientJoin) ClientType() string          { return "join" }
func (c ClientLeave) ClientType() string         { return "leave" }
func (c ClientKick) ClientType() string          { return "kick" }
func (c ClientStart) ClientType() string         { return "start" }
func (c ClientDraw) ClientType() string          { return "draw" }
func (c ClientSend) ClientType() string          { return "send" }
func (c ClientChat) ClientType() string          { return "chat" }

var ClientMessageTypes = slices.AssociateReverseBy([]ClientMessage{
	ClientChangeDetails{},
	ClientJoin{},
	ClientLeave{},
	ClientKick{},
	ClientStart{},
	ClientDraw{},
	ClientSend{},
	ClientChat{},
}, func(t ClientMessage) string { return t.ClientType() })

// ClientMessageFromJson converts a byte slice into a ClientMessage.
//go:todo avoid unmarshalling twice?
func (p *Player) ClientMessageFromJson(data []byte) (msg ClientMessage, err error) {
	var payload clientPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	for name, typeVal := range ClientMessageTypes {
		if payload.Type == name {
			c := reflect.New(reflect.TypeOf(typeVal))
			if err := json.Unmarshal(data, c.Interface()); err != nil {
				log.Printf("error: %v", err)
				return nil, err
			}
			c.Elem().FieldByName("Player").Set(reflect.ValueOf(p))
			return c.Elem().Interface().(ClientMessage), nil
		}
	}
	return nil, errors.New("unknown message type")
}
