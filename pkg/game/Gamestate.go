package game

import (
	"fmt"
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
	// Make copy of board
	newBoard := gs.Board

	// Validate deploy actions

	// Convert deploy actions into events
	deployEvents := ConvertDeployActionsToEvents(dpActions, &gs.Board)
	// Execute deploys
	executeDeploys(deployEvents, &newBoard)

	// Validate move actions with new board

	// Muster armies
	armies := musterArmies(mvActions, &newBoard)

	// Execute friendly movement
	moveEvents := executeMoves(&armies, &newBoard)
	// Execute battles
	battleEvents := executeBattles(&armies, &newBoard)

	stages.Travel = Gamestate{
		Board:   newBoard,
		Players: gs.Players,
	}

	// Execute sieges
	siegeEvents := ExecuteSieges(&armies, &newBoard)

	// Setup moves,battles and sieges
	logrus.Debugf("n moves:   %d", len(moveEvents))
	logrus.Debugf("n battles: %d", len(battleEvents))
	logrus.Debugf("n sieges:  %d", len(siegeEvents))

	// Check which players are still active
	newPlayers := checkActivePlayers(gs.Players, newBoard)
	// Calculate player income
	gs.Players = calculatePlayerIncome(gs.Players, newBoard)

	stages.DeployEvents = deployEvents
	stages.MoveEvents = moveEvents
	stages.BattleEvents = battleEvents
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

func checkActivePlayers(players []Player, brd board.Board) []Player {
	activePlayerSet := make(map[int]struct{})
	for _, tile := range brd.Tiles {
		activePlayerSet[tile.Owner] = struct{}{}
	}
	newPlayers := players
	for i, p := range newPlayers {
		_, ok := activePlayerSet[p.Name]
		newPlayers[i].Active = ok
	}
	return newPlayers
}

func calculatePlayerIncome(players []Player, brd board.Board) []Player {
	playerIdMap := map[int]*Player{}
	for i, p := range players {
		players[i].ArmiesToDeploy = 5
		playerIdMap[p.Name] = &players[i]
	}
	for _, tile := range brd.Tiles {
		if tile.Owner < 0 {
			continue // Owned by neutral player
		}
		playerIdMap[tile.Owner].ArmiesToDeploy += 1
	}
	return players
}

func executeDeploys(dps []events.DeployEvent, brd *board.Board) {
	for _, de := range dps {
		logrus.Debugf("Player [%d] deploys [%d] troops to [%d]", de.Player, de.After.NumTroops-de.Before.NumTroops, de.TileId)
		brd.Tiles[de.TileId].Garrison = de.After.NumTroops
	}
}

func musterArmies(mvActions []events.MoveAction, brd *board.Board) []events.Army {
	armies := []events.Army{}
	for _, mv := range mvActions {
		// Create the army
		army := events.Army{
			Player:      mv.Player,
			Troops:      mv.NumTroops,
			Destination: mv.TargetTileId,
			Origin:      mv.SourceTileId,
		}
		armies = append(armies, army)
		// Remove troops from tile
		brd.Tiles[mv.SourceTileId].Garrison -= mv.NumTroops
		if brd.Tiles[mv.SourceTileId].Garrison < 0 {
			logrus.Fatalf("Number of troops in tile [%d] is negative: [%d]", mv.SourceTileId, brd.Tiles[mv.SourceTileId].Garrison)
		}
	}
	return armies
}

func executeMoves(armies *[]events.Army, brd *board.Board) []events.MoveEvent {
	moveEvents := []events.MoveEvent{}
	for i := len(*armies) - 1; i >= 0; i-- {
		sourceTile := (*armies)[i].Origin
		targetTile := (*armies)[i].Destination
		// Check if owner is the same
		if brd.Tiles[sourceTile].Owner != brd.Tiles[targetTile].Owner {
			// This is a battle or siege
			continue
		}
		// Make the move on the board
		nTroops := (*armies)[i].Troops
		// if brd.Tiles[sourceTile].Garrison < 0 {
		// 	logrus.Fatalf("Number of troops in tile [%d] is negative: [%d]", sourceTile, brd.Tiles[sourceTile].Garrison)
		// }
		// brd.Tiles[sourceTile].Garrison -= nTroops
		brd.Tiles[targetTile].Garrison += nTroops

		logrus.Debugf("Player [%d] moves [%d] troops from [%d] to [%d]", (*armies)[i].Player, nTroops, sourceTile, targetTile)
		// Convert this army to a MoveEvent
		newMoveEvent := events.MoveEvent{(*armies)[i]}
		moveEvents = append(moveEvents, newMoveEvent)
		// Remove this armies from the army list
		*armies = append((*armies)[:i], (*armies)[i+1:]...)
	}
	return moveEvents
}

func executeBattles(armies *[]events.Army, brd *board.Board) []events.BattleEvent {
	// Group armies per link
	group := map[string][]int{}
	for i, army := range *armies {
		a, b := sortInts(army.Origin, army.Destination)
		key := fmt.Sprintf("%d-%d", a, b)
		if _, ok := group[key]; ok {
			group[key] = append(group[key], i)
		} else {
			group[key] = []int{i}
		}
	}
	// Record all indices of armies to delete
	armiesToDelete := []int{}
	battles := []events.BattleEvent{}
	for _, armyIndices := range group {
		if len(armyIndices) == 1 {
			// This will be a siege army
			continue
		}
		// Get participating armies from the list
		participatingArmies := []events.Army{}
		for _, i := range armyIndices {
			participatingArmies = append(participatingArmies, (*armies)[i])
			armiesToDelete = append(armiesToDelete, i)
		}

		battleEvent := GenerateBattleEventFromArmies(participatingArmies, brd)
		battles = append(battles, battleEvent)
		*armies = append(*armies, battleEvent.VictoriousArmy)
	}

	// Sort armies to delete descending
	sort.Slice(armiesToDelete, func(i, j int) bool { return i > j })
	// Remove them
	for _, indexToDelete := range armiesToDelete {
		*armies = append((*armies)[:indexToDelete], (*armies)[indexToDelete+1:]...)
	}

	return battles
}

func GenerateBattleEventFromArmies(armies []events.Army, brd *board.Board) events.BattleEvent {
	// Find tile ID
	tileId := GetLinkId(armies[0].Destination, armies[0].Origin, brd)

	// Sort armies by descending number of troops
	sort.Slice(armies, func(i, j int) bool { return armies[i].Troops > armies[j].Troops })

	// Create victorious army
	victoriousArmy := events.Army{
		Player:      armies[0].Player,
		Troops:      armies[0].Troops - armies[1].Troops,
		Destination: armies[0].Destination,
		Origin:      armies[0].Origin,
	}

	// Create battle event
	battleEvent := events.BattleEvent{
		Location:       tileId,
		Armies:         armies,
		VictoriousArmy: victoriousArmy,
	}

	return battleEvent
}

func ExecuteSieges(armies *[]events.Army, brd *board.Board) []events.SiegeEvent {
	sieges := []events.SiegeEvent{}

	siegedTiles := map[int][]events.Army{}

	for _, army := range *armies {
		if _, ok := siegedTiles[army.Destination]; ok {
			siegedTiles[army.Destination] = append(siegedTiles[army.Destination], army)
		} else {
			siegedTiles[army.Destination] = []events.Army{army}
		}
	}

	for _, participatingArmies := range siegedTiles {
		siegeEvent := ExecuteSiege(participatingArmies, brd)
		sieges = append(sieges, siegeEvent)
	}

	return sieges
}

func ExecuteSiege(armies []events.Army, brd *board.Board) events.SiegeEvent {
	tileId := armies[0].Destination
	besiegedPlayer := brd.Tiles[tileId].Owner
	logrus.Infof("Player [%d] is begin besieged on Tile [%d] with [%d] armies",
		besiegedPlayer, tileId, len(armies))

	// Create army for defenders
	defendingArmy := events.Army{
		Player:      besiegedPlayer,
		Troops:      brd.Tiles[tileId].Garrison,
		Destination: tileId,
		Origin:      tileId,
	}
	brd.Tiles[tileId].Garrison = 0
	armies = append(armies, defendingArmy)

	// Sort armies by DESCENDING num troops
	sort.Slice(armies, func(i, j int) bool { return armies[i].Troops > armies[j].Troops })

	for _, army := range armies {
		logrus.Infof("Player [%d] controls [%d] troops from [%d]", army.Player, army.Troops, army.Origin)
	}

	// Number of troops left for winning army
	numTroopsLeft := armies[0].Troops - armies[1].Troops

	// It's a draw. Does does not change owner
	victoriousPlayer := besiegedPlayer
	if numTroopsLeft > 0 {
		victoriousPlayer = armies[0].Player
	}

	brd.Tiles[tileId].Garrison = numTroopsLeft
	brd.Tiles[tileId].Owner = victoriousPlayer

	// Update siege
	siege := events.SiegeEvent{}
	siege.Armies = armies
	siege.Location = tileId
	siege.PostSiegeState.NumTroops = numTroopsLeft
	siege.PostSiegeState.Owner = victoriousPlayer

	logrus.Infof("Siege won by player [%d] with [%d] left", victoriousPlayer, numTroopsLeft)

	return siege
}
