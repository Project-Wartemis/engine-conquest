package game

import (
	"github.com/Project-Wartemis/pw-engine/pkg/game/board"
	"github.com/Project-Wartemis/pw-engine/pkg/game/events"
	"github.com/Project-Wartemis/pw-engine/pkg/game/export"
)

type TurnStages struct {
	Start        Gamestate            `json:"start"`
	Travel       Gamestate            `json:"travel"`
	End          Gamestate            `json:"end"`
	DeployEvents []events.DeployEvent `json:"deployEvents"`
	MoveEvents   []events.MoveEvent   `json:"moveEvent"`
	BattleEvents []events.BattleEvent `json:"battleEvent"`
	SiegeEvents  []events.SiegeEvent  `json:"siegeEvents"`
}

func (stages *TurnStages) Export() export.Gamestate {
	nodes := stages.exportNodes()
	links := stages.exportLinks()
	players := stages.exportPlayers()

	deploys := stages.exportDeploys()
	moves := stages.exportMoves()

	travel := stages.exportBattles()
	combat := stages.exportSieges()
	fights := export.Fights{
		Travel: travel,
		Combat: combat,
	}

	nodeStages := stages.exportStages()

	result := export.Gamestate{}
	result.Nodes = nodes
	result.Links = links
	result.Players = players
	result.Deploys = deploys
	result.Moves = moves
	result.Fights = fights
	result.Stages = nodeStages

	return result
}

func (stages TurnStages) exportNodes() []export.Node {
	nodes := []export.Node{}
	for id := range stages.Start.Board.Tiles {
		n := export.Node{
			Id:   id,
			Name: string(id),
		}
		nodes = append(nodes, n)
	}
	return nodes
}

func (stages TurnStages) exportLinks() []export.Link {
	links := []export.Link{}
	for id, l := range stages.Start.Board.Links {
		L := export.Link{
			Id: id,
			A:  l.A.Id,
			B:  l.B.Id,
		}
		links = append(links, L)
	}
	return links
}

func (stages TurnStages) exportPlayers() []export.Player {
	players := []export.Player{}
	for _, p := range stages.End.Players {
		if p.Active {
			P := export.Player{
				Id:     p.Name,
				Income: p.ArmiesToDeploy,
			}
			players = append(players, P)
		}
	}
	return players
}

func (stages TurnStages) exportDeploys() []export.Deploy {
	deploys := []export.Deploy{}
	for _, d := range stages.DeployEvents {
		D := export.Deploy{
			Node:   d.TileId,
			Troops: d.After.NumTroops - d.Before.NumTroops,
		}
		deploys = append(deploys, D)
	}
	return deploys
}

func (stages TurnStages) exportMoves() []export.Move {
	moves := []export.Move{}
	for _, mv := range stages.MoveEvents {
		MV := export.Move{
			Source: mv.MoveAction.SourceTileId,
			Target: mv.MoveAction.TargetTileId,
			Troops: mv.MoveAction.NumTroops,
		}
		moves = append(moves, MV)
	}
	return moves
}

func (stages TurnStages) exportBattles() []export.Fight {
	fights := []export.Fight{}
	for _, bat := range stages.BattleEvents {
		F := export.Fight{
			Location: bat.Location,
			Armies:   exportArmies(bat.Armies),
		}
		fights = append(fights, F)
	}
	return fights
}

func (stages TurnStages) exportSieges() []export.Fight {
	fights := []export.Fight{}
	for _, bat := range stages.SiegeEvents {
		F := export.Fight{
			Location: bat.Location,
			Armies:   exportArmies(bat.Armies),
		}
		fights = append(fights, F)
	}
	return fights
}

func exportArmies(armies []events.Army) []export.Army {
	result := []export.Army{}
	for _, arm := range armies {
		A := export.Army{
			Player: arm.Player,
			Troops: arm.Troops,
		}
		result = append(result, A)
	}
	return result
}

func (stages TurnStages) exportStages() export.Stages {
	result := export.Stages{}
	result.Start = exportNodeStates(stages.Start.Board)
	result.Travel = exportNodeStates(stages.Travel.Board)
	result.End = exportNodeStates(stages.End.Board)
	return result
}

func exportNodeStates(brd board.Board) []export.NodeState {
	result := []export.NodeState{}
	for _, t := range brd.Tiles {
		ns := export.NodeState{
			Id:     t.Id,
			Owner:  t.Owner,
			Troops: t.Garrison,
		}
		result = append(result, ns)
	}
	return result
}
