package game

import (
	"cardgame/deck"
	"cardgame/util"
	"cardgame/words"
	"log"
	"strings"
	"time"
)

type hubMessage struct {
	clientMessage ClientMessage
	player        *Player
}

type Hub struct {
	RegionCode string
	Version    string

	Rooms map[string]*Room // RoomId -> Room

	inbound chan *hubMessage // incoming client messages
}

func (h *Hub) NewRoom(password string) *Room {
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
	h.Rooms[r.Id] = &r
	if password != "" {
		r.SetPassword(password)
	}

	go r.read()
	go r.write()

	return &r
}

func (h *Hub) read() {
	for {
		select {
		case msg, ok := <-h.inbound:
			if !ok {
				return
			}

			switch clientMessage := msg.clientMessage.(type) {
			case ClientJoin:
				h.handleJoin(clientMessage)
			case ClientLeave:
				h.handleLeave(clientMessage)
			default:
				log.Printf("[error] bad message type %T sent to hub\n", clientMessage)
			}
		}
	}
}

func (h *Hub) handleJoin(msg ClientJoin) {
	p := msg.Player
	r, ok := h.Rooms[msg.RoomId]

	if p.room != nil {
		p.outbound <- &ServerError{"You are already in a room"}
		return
	}

	if r == nil || !ok {
		p.outbound <- &ServerError{"Room not found"}
		return
	}

	if r.IsPrivate() && !r.CheckPassword(msg.Password) {
		p.outbound <- &ServerError{"Incorrect password"}
		return
	}

	r.inbound <- msg
}

func (h *Hub) handleLeave(msg ClientLeave) {
	if msg.Player.room != nil {
		msg.Player.room.inbound <- msg
		msg.Player.room = nil
	}
}

var HubMain *Hub

func init() {
	HubMain = &Hub{
		RegionCode: "global",

		Rooms:   make(map[string]*Room),
		inbound: make(chan *hubMessage),
	}
	go HubMain.read()
}
