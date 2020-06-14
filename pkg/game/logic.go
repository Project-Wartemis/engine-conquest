package game

import (
	"github.com/Project-Wartemis/pw-engine/pkg/game/board"
	"github.com/Project-Wartemis/pw-engine/pkg/game/events"
	"github.com/sirupsen/logrus"
)

func DeployEventFromAction(da events.DeployAction, b *board.Board) events.DeployEvent {
	return events.DeployEvent{
		TileId: da.TileId,
		Player: da.Player,
		Before: events.State{NumTroops: b.Tiles[da.TileId].Garrison},
		After:  events.State{NumTroops: b.Tiles[da.TileId].Garrison + da.NumTroops},
	}
}

func ConvertDeployActionsToEvents(dps []events.DeployAction, brd *board.Board) []events.DeployEvent {
	deployEvents := []events.DeployEvent{}
	for _, da := range dps {
		dpEvent := DeployEventFromAction(da, brd)
		deployEvents = append(deployEvents, dpEvent)
	}
	return deployEvents
}

func GetLinkId(a int, b int, brd *board.Board) int {
	a, b = sortInts(a, b)
	for i, l := range brd.Links {
		if l.A.Id == a && l.B.Id == b {
			return i
		}
	}
	logrus.Fatalf("ABORT: Could not find link (%d, %d)", a, b)
	return -1
}

func sortInts(a int, b int) (int, int) {
	if a > b {
		return b, a
	}
	return a, b
}
