package game

import (
	"cardgame/card"
	"cardgame/util"
	"cardgame/words"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer. (1MB)
	maxMessageSize = 1 << 20
)

type Player struct {
	Id       string       `json:"id"`
	Avatar   AvatarConfig `json:"avatar"`
	Name     string       `json:"name"`
	Score    int          `json:"score"`
	Hand     PlayerHand   `json:"cards"` // Player's hand, top is at the end
	socket   *websocket.Conn
	room     *Room
	outbound chan ServerMessage // outgoing server messages
}

type PlayerHand []*card.Card

// top returns the top card of the player's hand.
func (h PlayerHand) top() *card.Card {
	if len(h) == 0 {
		return nil
	}
	return h[len(h)-1]
}

// tail returns all cards from the player's hand except the top card.
func (h PlayerHand) tail() PlayerHand {
	if len(h) == 0 {
		return PlayerHand{}
	}
	return h[:len(h)-1]
}

type AvatarConfig struct {
	Eyes  int `json:"eyes"`
	Mouth int `json:"mouth"`
	Color int `json:"color"`
}

// read pumps messages from the websocket connection to the room.
//
// The application runs read in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (p *Player) read() {
	defer func() {
		p.room.inbound <- ClientLeave{p}
		p.socket.Close()
	}()
	p.socket.SetReadLimit(maxMessageSize)
	p.socket.SetReadDeadline(time.Now().Add(pongWait))
	p.socket.SetPongHandler(func(string) error { p.socket.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, mesageData, err := p.socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg, err := p.ClientMessageFromJson(mesageData)
		if err != nil {
			log.Printf("error: %v", err)
			p.outbound <- ServerError{
				Message: err.Error(),
			}
			return
		}

		p.room.inbound <- msg
	}
}

// write pumps messages from the room to the websocket connection.
//
// A goroutine running write is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (p *Player) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		p.socket.Close()
	}()
	for {
		select {
		case message, ok := <-p.outbound:
			p.socket.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// room closed
				p.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			fmt.Println("player->", message)

			w, err := p.socket.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("error: %v\n", err)
				return
			}

			s := structs.New(message)
			s.TagName = "json"
			m := s.Map()
			m["type"] = message.ServerType()
			m["room"] = p.room

			if err := json.NewEncoder(w).Encode(m); err != nil {
				return
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			p.socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func NewPlayer(socket *websocket.Conn, r *Room) *Player {
	p := &Player{
		Id:       util.IdFrom("p", socket.RemoteAddr().String()),
		Name:     strings.Join(words.Words(words.English, 2), " "),
		socket:   socket,
		room:     r,
		Hand:     PlayerHand{},
		outbound: make(chan ServerMessage),
	}

	go p.read()
	go p.write()

	return p
}
