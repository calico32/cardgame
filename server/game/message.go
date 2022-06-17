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
	serverPayload struct {
		include set
		exclude set
		message ServerMessage
	}

	ClientMessage interface{ ClientType() string }
	ServerMessage interface{ ServerType() string }
	// ClientChangeDetails is sent by the room owner to change the room details and add/remove decks
	ClientChangeDetails struct {
		Player *Player `json:"-"`

		Name        *string   `json:"name"`
		Description *string   `json:"description"`
		MaxPlayers  *int      `json:"max_players"`
		Password    *string   `json:"password"`      // new password for private rooms, or "" for public rooms
		AddDecks    []string  `json:"add_decks"`     // IDs of decks to add
		RemoveDecks []string  `json:"remove_decks"`  // IDs of decks to remove
		PlayMode    *PlayMode `json:"play_mode"`     // new play mode
		HubDeviceId *string   `json:"hub_device_id"` // ID of the hub device to use
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

		TargetId string `json:"target_id"`
	}
	// ClientChat is sent by a player to send a chat message.
	ClientChat struct {
		Player *Player `json:"-"`

		Message     string  `json:"message"`
		RecipientId *string `json:"recipient"` // RecipientId is set if the message is a private message.
	}

	// ServerChangeDetails is sent to all players when the room details change.
	ServerChangeDetails struct {
		Name        *string   `json:"name"`
		Description *string   `json:"description"`
		MaxPlayers  *int      `json:"max_players"`
		DeckIds     []string  `json:"decks"`
		PlayMode    *PlayMode `json:"play_mode"`
		HubDeviceId *string   `json:"hub_device_id"`
	}
	// ServerRoomDetails is sent to a player when they join the room.
	ServerRoomDetails struct {
		Room *Room `json:"room"`
	}
	// ServerJoin is sent to all players when a new player joins the room.
	ServerJoin struct {
		Id     string `json:"id"`
		Player Player `json:"player"`
	}
	// ServerLeave is sent to all players when a player leaves the room.
	ServerLeave struct {
		Id string `json:"id"`
	}
	// ServerError is sent to a player when an error occurs.
	ServerError struct {
		Message string `json:"message"`
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

func (s ServerChangeDetails) ServerType() string { return "change_details" }
func (s ServerRoomDetails) ServerType() string   { return "room_details" }
func (s ServerJoin) ServerType() string          { return "join" }
func (s ServerLeave) ServerType() string         { return "leave" }
func (s ServerError) ServerType() string         { return "error" }

var clientMessageTypes = slices.AssociateReverseBy([]ClientMessage{
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
func (p *Player) ClientMessageFromJson(data []byte) (msg ClientMessage, err error) {
	var payload clientPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	for name, typeVal := range clientMessageTypes {
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
