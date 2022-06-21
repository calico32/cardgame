package game

type PlayMode int

const (
	PlayModePlayersOnly PlayMode = iota
	PlayModePlayersAndHub
	PlayModeHubOnly
)

var TSAllPlayModes = []struct {
	Value  PlayMode
	TSName string
}{
	{PlayModePlayersOnly, "PlayersOnly"},
	{PlayModePlayersAndHub, "PlayersAndHub"},
	{PlayModeHubOnly, "HubOnly"},
}

type GamePhase int

const (
	GamePhaseLobby GamePhase = iota
	GamePhasePlaying
	GamePhaseEnd
)

var TSAllGamePhases = []struct {
	Value  GamePhase
	TSName string
}{
	{GamePhaseLobby, "Lobby"},
	{GamePhasePlaying, "Playing"},
	{GamePhaseEnd, "End"},
}
