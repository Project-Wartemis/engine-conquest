package game

import (
	"sort"

	"github.com/Project-Wartemis/pw-engine/pkg/game/board"
	"github.com/Project-Wartemis/pw-engine/pkg/game/events"
	"github.com/sirupsen/logrus"
)

type Gamestate struct {
	Board   board.Board
	Players []Player
}

func (gs *Gamestate) ProcessTurn(dpActions []events.DeployAction, mvActions []events.MoveAction) TurnStages {
	logrus.Debug("New Turn!")
	stages := TurnStages{}
	stages.Start = *gs
	// Validate deploy actions
	// Validate move actions

	// Convert deploy actions into events
	deployEvents := ConvertDeployActionsToEvents(dpActions, &gs.Board)
	// Make copy of board
	newBoard := gs.Board

	// Execute deploys
	ExecuteDeploys(deployEvents, &newBoard)

	// Validate move actions with new board

	// Setup and execute moves
	moveEvents := CreateMoveEventsFromActions(&mvActions, &newBoard)
	logrus.Debugf("n moves: %d", len(moveEvents))
	ExecuteMoveEvents(moveEvents, &newBoard)

	// Setup and execute battles
	battleEvents := CreateBattleEventsFromActions(&mvActions, &newBoard)
	logrus.Debugf("n battles: %d", len(battleEvents))
	executedBattleEvents, nonBattles := ExecuteBattles(battleEvents, &newBoard)

	// Save Travel gamestate
	stages.Travel = Gamestate{
		Board:   newBoard,
		Players: gs.Players,
	}

	// Setup and execute sieges
	siegeEvents := CreateSiegeEventsFromActions(&nonBattles, &newBoard)
	logrus.Debugf("n sieges: %d", len(siegeEvents))
	siegeEvents = ExecuteSieges(siegeEvents, &newBoard)

	// Check which players are still active
	activePlayerSet := make(map[int]struct{})
	for _, tile := range newBoard.Tiles {
		activePlayerSet[tile.Owner] = struct{}{}
	}
	newPlayers := gs.Players
	for i, p := range newPlayers {
		_, ok := activePlayerSet[p.Name]
		newPlayers[i].Active = ok
	}

	stages.DeployEvents = deployEvents
	stages.MoveEvents = moveEvents
	stages.BattleEvents = executedBattleEvents
	stages.SiegeEvents = siegeEvents
	stages.End = Gamestate{
		Board:   newBoard,
		Players: newPlayers,
	}

	logrus.Debugf("End of turn:")
	for _, t := range stages.End.Board.Tiles {
		logrus.Debugf("Tile [%d] owned by player [%d] with [%d] troops", t.Id, t.Owner, t.Garrison)
	}

	return stages
}

func ExecuteDeploys(dps []events.DeployEvent, brd *board.Board) {
	for _, de := range dps {
		logrus.Debugf("Player [%d] deploys [%d] troops to [%d]", de.Player, de.After.NumTroops-de.Before.NumTroops, de.TileId)
		brd.Tiles[de.TileId].Garrison = de.After.NumTroops
	}
}

func ExecuteMoveEvents(mvs []events.MoveEvent, brd *board.Board) {
	for _, mv := range mvs {
		sourceTile := mv.MoveAction.SourceTileId
		targetTile := mv.MoveAction.TargetTileId
		numTroops := mv.MoveAction.NumTroops
		logrus.Debugf("Player [%d] moves [%d] troops from [%d] to [%d]", mv.MoveAction.Player, numTroops, sourceTile, targetTile)
		brd.Tiles[sourceTile].Garrison -= numTroops
		brd.Tiles[targetTile].Garrison += numTroops
	}
}

func ExecuteBattles(battles []events.BattleEvent, brd *board.Board) ([]events.BattleEvent, []events.BattleEvent) {
	for i, battle := range battles {
		battles[i] = ExecuteBattle(battle, brd)
	}
	return battles, battles // UPDATE THIS
}
func ExecuteSieges(sieges []events.SiegeEvent, brd *board.Board) []events.SiegeEvent {
	for i, siege := range sieges {
		sieges[i] = ExecuteSiege(siege, brd)
	}
	return sieges
}

func ExecuteBattle(battle events.BattleEvent, brd *board.Board) events.BattleEvent {
	sort.Slice(battle.Armies, func(i, j int) bool { return battle.Armies[i].Troops < battle.Armies[j].Troops })
	numTroopsLeft := battle.Armies[len(battle.Armies)-1].Troops - battle.Armies[len(battle.Armies)-2].Troops

	battle.VictoriousArmy = battle.Armies[len(battle.Armies)-1]
	battle.VictoriousArmy.Troops = numTroopsLeft

	// Update the board
	logrus.Infof("Battle at link [%d] between [%d] armies", battle.Location, len(battle.Armies))
	for _, army := range battle.Armies {
		brd.Tiles[army.Origin].Garrison -= army.Troops
		logrus.Infof("Player [%d] controls [%d] troops from [%d] with destination [%d]", army.Player, army.Troops, army.Origin, army.Destination)
	}
	logrus.Infof("Battle won by player [%d]. [%d] troops left, moving to [%d]",
		battle.VictoriousArmy.Player,
		battle.VictoriousArmy.Troops,
		battle.VictoriousArmy.Destination)

	return battle
}

func ExecuteSiege(siege events.SiegeEvent, brd *board.Board) events.SiegeEvent {
	// Sort armies by DESCENDING num troops
	sort.Slice(siege.Armies, func(i, j int) bool { return siege.Armies[i].Troops > siege.Armies[j].Troops })
	logrus.Infof("Siege at tile [%d] with [%d] armies", siege.Location, len(siege.Armies))
	// Remove troops from origins
	for _, army := range siege.Armies {
		brd.Tiles[army.Origin].Garrison -= army.Troops
		logrus.Infof("Player [%d] controls [%d] troops from [%d]", army.Player, army.Troops, army.Origin)
	}

	// Number of troops left for winning army
	numTroopsLeft := siege.Armies[0].Troops - siege.Armies[1].Troops
	victoriousPlayer := siege.Armies[0].Player

	brd.Tiles[siege.Location].Garrison = numTroopsLeft
	brd.Tiles[siege.Location].Owner = victoriousPlayer

	// Update siege
	siege.PostSiegeState.NumTroops = numTroopsLeft
	siege.PostSiegeState.Owner = victoriousPlayer

	logrus.Infof("Siege won by player [%d] with [%d] left from tile [%d]", victoriousPlayer, numTroopsLeft, siege.Armies[0].Origin)
	logrus.Infof("Owner of siege location: [%d]", brd.Tiles[siege.Location].Owner)
	logrus.Infof("Number of troops in siege location: [%d]", brd.Tiles[siege.Location].Garrison)

	return siege
}
