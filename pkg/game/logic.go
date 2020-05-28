package game

import (
	"fmt"

	"github.com/Project-Wartemis/pw-engine/pkg/game/board"
	"github.com/Project-Wartemis/pw-engine/pkg/game/events"
	"github.com/sirupsen/logrus"
)

func MoveEventFromAction(ma events.MoveAction, b *board.Board) events.MoveEvent {
	moveEvent := events.MoveEvent{}
	moveEvent.MoveAction = ma

	moveEvent.PreMoveState.SourceNumTroops = b.Tiles[ma.SourceTileId].Garrison
	moveEvent.PreMoveState.TargetNumTroops = b.Tiles[ma.TargetTileId].Garrison

	moveEvent.PostMoveState.SourceNumTroops = b.Tiles[ma.SourceTileId].Garrison - ma.NumTroops
	moveEvent.PostMoveState.TargetNumTroops = b.Tiles[ma.TargetTileId].Garrison + ma.NumTroops

	return moveEvent
}

func DeployEventFromAction(da events.DeployAction, b *board.Board) events.DeployEvent {
	return events.DeployEvent{
		TileId: da.TileId,
		Player: da.Player,
		Before: events.State{NumTroops: b.Tiles[da.TileId].Garrison},
		After:  events.State{NumTroops: b.Tiles[da.TileId].Garrison + da.NumTroops},
	}
}

func CreateMoveEventsFromActions(ms *[]events.MoveAction, brd *board.Board) []events.MoveEvent {
	moves := []events.MoveEvent{}
	// Filter out all moves that are friendly
	for i := len(*ms) - 1; i >= 0; i-- {
		m := (*ms)[i]
		if brd.Tiles[(*ms)[i].SourceTileId].Owner == brd.Tiles[(*ms)[i].TargetTileId].Owner {
			moves = append(moves, MoveEventFromAction(m, brd))
			new_ms := append((*ms)[:i], (*ms)[i+1:]...)
			*ms = new_ms
		}
	}
	return moves
}

func CreateBattleEventsFromActions(ms *[]events.MoveAction, brd *board.Board) []events.BattleEvent {
	battles := make(map[string]*events.BattleEvent)

	for i := len(*ms) - 1; i >= 0; i-- {
		a, b := sortInts((*ms)[i].SourceTileId, (*ms)[i].TargetTileId)
		linkId := getLinkId(a, b, brd)
		// Make army for this action
		army := events.Army{
			Player:      (*ms)[i].Player,
			Troops:      (*ms)[i].NumTroops,
			Destination: linkId,
			Origin:      (*ms)[i].SourceTileId,
		}

		// Generate ID for link
		id := fmt.Sprintf("%d-%d", a, b)
		// Check if this battle already exists
		_, ok := battles[id]
		if ok {
			//  the battle already exists, add the army
			battles[id].Armies = append(battles[id].Armies, army)
		} else {
			// Create a new Battle
			battles[id] = &events.BattleEvent{
				Armies:   []events.Army{army},
				Location: linkId,
			}
		}
	}

	// Convert map to list
	battleList := []events.BattleEvent{}
	for _, val := range battles {
		battleList = append(battleList, *val)
	}

	return battleList
}

func CreateSiegeEventsFromActions(battles *[]events.BattleEvent, brd *board.Board) []events.SiegeEvent {
	sieges := map[int]*events.SiegeEvent{}
	for i := len(*battles) - 1; i >= 0; i-- {
		// Check number of participating Armies
		if len((*battles)[i].Armies) > 1 {
			// This is an actual Battle
			continue
		}
		// Single army, convert to Combat
		combatLocation := destinatationFromLinkID((*battles)[i].Armies[0].Destination, (*battles)[i].Armies[0].Origin, brd)
		attackingArmy := (*battles)[i].Armies[0]

		_, ok := sieges[combatLocation]
		if ok {
			sieges[combatLocation].Armies = append(sieges[combatLocation].Armies, attackingArmy)
		} else {
			defenders := armyFromTile(*brd.Tiles[combatLocation])
			newSiege := events.SiegeEvent{
				Location: combatLocation,
				Armies:   []events.Army{attackingArmy, defenders},
			}
			sieges[combatLocation] = &newSiege
		}
		// Remove from travelFights list
		newBattles := append((*battles)[:i], (*battles)[i+1:]...)
		*battles = newBattles
	}
	siegeList := []events.SiegeEvent{}
	for _, val := range sieges {
		siegeList = append(siegeList, *val)
	}

	return siegeList
}

func armyFromTile(tile board.Tile) events.Army {
	return events.Army{
		Player:      tile.Owner,
		Troops:      tile.Garrison,
		Destination: tile.Id,
	}
}

func ConvertActionsToEvents(mvs []events.MoveAction, brd *board.Board) ([]events.MoveEvent, []events.BattleEvent, []events.SiegeEvent) {
	// Extract all friendly moves
	moves := GenerateMoveEvents(&mvs, brd)

	// Convert all other moves to travel fights
	battles := GenerateBattleEvents(&mvs, brd)

	// Convert some travelFights to Combat
	sieges := GenerateSiegeEvents(&battles, brd)
	return moves, battles, sieges
}

func ConvertDeployActionsToEvents(dps []events.DeployAction, brd *board.Board) []events.DeployEvent {
	deployEvents := []events.DeployEvent{}
	for _, da := range dps {
		dpEvent := DeployEventFromAction(da, brd)
		deployEvents = append(deployEvents, dpEvent)
	}
	return deployEvents
}

func destinatationFromLinkID(linkID int, sourceTileId int, brd *board.Board) int {
	for _, link := range brd.Links {
		if link.Id == linkID {
			if sourceTileId == link.A.Id {
				return link.B.Id
			}
			return link.A.Id
		}
	}
	return -1
}

func getLinkId(a int, b int, brd *board.Board) int {
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
