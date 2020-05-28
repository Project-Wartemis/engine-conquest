package events

type Army struct {
	Player      int
	Troops      int
	Destination int
	Origin      int
}

type TileState struct {
	NumTroops int
	Owner     int
}

type SiegeEvent struct {
	Location       int
	Armies         []Army
	PostSiegeState TileState
}
