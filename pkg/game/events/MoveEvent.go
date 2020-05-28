package events

type MoveState struct {
	SourceNumTroops int
	TargetNumTroops int
}

type MoveEvent struct {
	MoveAction    MoveAction
	PreMoveState  MoveState
	PostMoveState MoveState
}
