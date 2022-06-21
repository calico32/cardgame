package game

type PlayMode int

const (
	PlayModePlayersOnly PlayMode = iota
	PlayModePlayersAndHub
	PlayModeHubOnly
)

type GamePhase int

const (
	GamePhaseLobby GamePhase = iota
	GamePhasePlaying
	GamePhaseEnd
)
