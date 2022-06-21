/* Do not change, this code is generated from Golang structs */

export type ClientMessage =
    | ({ type: "change_details" } & ClientChangeDetails)
    | ({ type: "join" } & ClientJoin)
    | ({ type: "leave" } & ClientLeave)
    | ({ type: "kick" } & ClientKick)
    | ({ type: "start" } & ClientStart)
    | ({ type: "draw" } & ClientDraw)
    | ({ type: "send" } & ClientSend)
    | ({ type: "chat" } & ClientChat)

export type ServerMessage =
    | ({ room: Room; type: "join" } & ServerJoin)
    | ({ room: Room; type: "leave" } & ServerLeave)
    | ({ room: Room; type: "draw" } & ServerDraw)
    | ({ room: Room; type: "wild_card" } & ServerWildCard)
    | ({ room: Room; type: "chat" } & ServerChat)
    | ({ room: Room; type: "change_details" } & ServerChangeDetails)
    | ({ room: Room; type: "kick" } & ServerKick)
    | ({ room: Room; type: "start" } & ServerStart)
    | ({ room: Room; type: "reshuffle" } & ServerReshuffle)
    | ({ room: Room; type: "resync" } & ServerResync)
    | ({ room: Room; type: "turn" } & ServerTurn)
    | ({ room: Room; type: "error" } & ServerError)
    | ({ room: Room; type: "ack" } & ServerAck)


export enum GamePhase {
    Lobby = 0,
    Playing = 1,
    End = 2,
}
export enum PlayMode {
    PlayersOnly = 0,
    PlayersAndHub = 1,
    HubOnly = 2,
}
export enum CardType {
    Lines = 0,
    Waves = 1,
    Square = 2,
    Dots = 3,
    Hash = 4,
    Circle = 5,
    Plus = 6,
    Star = 7,
    Invalid = -1,
}
export interface WildCard {
    id: string;
    types: number[];
}
export interface Deck {
    id: string;
    name: string;
    location: string;
    description: string;
    cards: Card[];
    wildCards: WildCard[];
}
export interface Card {
    id: string;
    type: CardType;
    category: string;
}
export interface AvatarConfig {
    eyes: number;
    mouth: number;
    color: number;
}
export interface Player {
    id: string;
    avatar: AvatarConfig;
    name: string;
    score: number;
    cards: Card[];
}
export interface Room {
    id: string;
    timestamp: number;
    name: string;
    description: string;
    maxPlayers: number;
    ownerId: string;
    players: Player[];
    decks: Deck[];
    playMode: PlayMode;
    hubDeviceId: string;
    currentTurn: number;
    gamePhase: GamePhase;
    activeWildCard?: WildCard;
    drawPileSize: number;
}




export interface ClientJoin {
    name: string;
    avatar: AvatarConfig;
}
export interface ClientLeave {

}
export interface ClientKick {
    id: string;
}
export interface ClientStart {

}
export interface ClientDraw {

}
export interface ClientSend {
    recipientId: string;
}
export interface ClientChat {
    message: string;
    recipient?: string;
}
export interface ClientChangeDetails {
    name?: string;
    description?: string;
    maxPlayers?: number;
    password?: string;
    addDecks: string[];
    removeDecks: string[];
    playMode?: PlayMode;
    hubDeviceId?: string;
}
export interface ServerTurn {
    playerId: string;
}
export interface ServerError {
    message: string;
}
export interface ServerAck {

}
export interface ServerKick {

}
export interface ServerStart {
    currentTurn: number;
}
export interface ServerReshuffle {

}
export interface ServerResync {
    topCards: {[key: string]: Card};
}
export interface ServerChat {
    timestamp: string;
    player: string;
    private: boolean;
    message: string;
}
export interface ServerChangeDetails {
    name?: string;
    description?: string;
    maxPlayers?: number;
    decks: string[];
    playMode?: PlayMode;
    hubDeviceId?: string;
}
export interface ServerJoin {
    id: string;
    player: Player;
}
export interface ServerLeave {
    id: string;
}
export interface ServerDraw {
    playerId: string;
    card?: Card;
}
export interface ServerWildCard {
    playerId: string;
    card?: WildCard;
}