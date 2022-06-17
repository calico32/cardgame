package game

import (
	"cardgame/deck"
	"cardgame/util"
	"cardgame/words"
	"fmt"
	"strings"
	"time"
)

var Rooms = make(map[string]*Room)

type PlayMode int

const (
	PlayModePlayersOnly PlayMode = iota
	PlayModePlayersAndHub
	PlayModeHubOnly
)

// Room represents a game room.
type Room struct {
	Id          string       `json:"id"`            // internal room id
	Timstamp    int64        `json:"timestamp"`     // creation timestamp
	Name        string       `json:"name"`          // public-facing room name
	Description string       `json:"description"`   // room description
	MaxPlayers  int          `json:"max_players"`   // maximum number of players
	OwnerId     string       `json:"ownerId"`       // owner's player id
	Players     []*Player    `json:"players"`       // players in the room, including the owner
	Decks       []*deck.Deck `json:"decks"`         // decks in use
	PlayMode    PlayMode     `json:"play_mode"`     // play mode
	HubDeviceId string       `json:"hub_device_id"` // hub device id

	private      bool   // private room?
	passwordHash string // password hash for private rooms

	inbound  chan ClientMessage  // incoming client messages
	outbound chan *serverPayload // outgoing server messages
}

// IsPrivate returns true if the room is private.
func (r *Room) IsPrivate() bool {
	return r.private
}

// TryPassword returns true if the password is correct.
func (r *Room) TryPassword(password string) bool {
	return checkPassword(r.passwordHash, password)
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

func NewRoom(password string) *Room {
	id := util.IdFrom("r", time.Now().String())
	r := Room{
		Id:         id,
		Name:       strings.Join(words.Words(words.English, 4), " "),
		Timstamp:   time.Now().UnixMilli(),
		Players:    []*Player{},
		Decks:      []*deck.Deck{},
		MaxPlayers: 4,
		inbound:    make(chan ClientMessage),
		outbound:   make(chan *serverPayload),
	}
	Rooms[r.Id] = &r
	if password != "" {
		r.SetPassword(password)
	}

	go r.read()
	go r.write()

	return &r
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

		if len(included) > 0 {
			fmt.Println("room->", payload.message)
			for _, p := range included {
				fmt.Println("room->    sending to", p.Id)
				p.outbound <- payload.message
			}
		} else if len(other) > 0 {
			fmt.Println("room->", payload.message)
			for _, p := range other {
				fmt.Println("room->    sending to", p.Id)
				p.outbound <- payload.message
			}
		}
	}
}
