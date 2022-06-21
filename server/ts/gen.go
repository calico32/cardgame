package main

import (
	"cardgame/card"
	"cardgame/deck"
	"cardgame/game"
	"reflect"
	"strings"

	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

func main() {
	converter := typescriptify.New().WithInterface(true).
		Add(game.Room{}).
		Add(game.Player{}).
		Add(deck.Deck{}).
		Add(card.Card{}).
		Add(card.WildCard{}).
		AddEnum(game.TSAllGamePhases).
		AddEnum(game.TSAllPlayModes).
		AddEnum(card.TSAllCardTypes)

	for _, typ := range game.ClientMessageTypes {
		converter = converter.Add(typ)
	}

	for _, typ := range game.ServerMessageTypes {
		converter = converter.Add(typ)
	}

	var extras strings.Builder
	// export type ClientMessage = { type: "join" } & ClientJoin | { type: "leave" } & ClientLeave;
	extras.WriteString("export type ClientMessage =\n")
	for t, typ := range game.ClientMessageTypes {
		typeName := reflect.TypeOf(typ).Name()
		extras.WriteString("    | ")
		extras.WriteString("({ type: \"")
		extras.WriteString(t)
		extras.WriteString("\" } & ")
		extras.WriteString(typeName)
		extras.WriteString(")\n")
	}
	extras.WriteString("\n")
	extras.WriteString("export type ServerMessage =\n")
	for t, typ := range game.ServerMessageTypes {
		typeName := reflect.TypeOf(typ).Name()
		extras.WriteString("    | ")
		extras.WriteString("({ room: Room; type: \"")
		extras.WriteString(t)
		extras.WriteString("\" } & ")
		extras.WriteString(typeName)
		extras.WriteString(")\n")
	}

	converter.AddImport(extras.String())

	err := converter.ConvertToFile("ts/models.ts")
	if err != nil {
		panic(err.Error())
	}
}
