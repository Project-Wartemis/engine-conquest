package events

type State struct {
	NumTroops int
}

type DeployEvent struct {
	TileId int
	Player int
	Before State
	After  State
}
