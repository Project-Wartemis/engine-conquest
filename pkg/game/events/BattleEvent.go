package events

type BattleEvent struct {
	Location       int
	Armies         []Army
	VictoriousArmy Army
}
