package game

type Player struct {
	Name           int
	ArmiesToDeploy int
	Active         bool
}

func NewPlayer(name int, armiesToDeploy int) *Player {
	return &Player{
		Name:           name,
		ArmiesToDeploy: armiesToDeploy,
		Active:         true,
	}
}
