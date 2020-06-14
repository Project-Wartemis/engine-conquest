package board

import (
	"sort"

	"github.com/sirupsen/logrus"
)

type Board struct {
	Tiles       []*Tile
	Links       []Link
	Connections [][]bool
}

type Connection struct {
	A int
	B int
}

func NewEmptyBoard(numTiles int, connections []Connection) *Board {
	sort.Slice(connections, func(a, b int) bool {
		return connections[a].A < connections[b].A || (connections[a].A == connections[b].A && connections[a].B < connections[b].B)
	})
	board := Board{}

	// Build tiles
	board.Tiles = make([]*Tile, numTiles)
	for i := 0; i < numTiles; i++ {
		board.Tiles[i] = NewTile(i)
	}

	// Initialise all connections to false
	board.Connections = [][]bool{}
	for i := 0; i < numTiles; i++ {
		board.Connections = append(board.Connections, make([]bool, numTiles))
	}

	// Make connection
	for id, pair := range connections {
		// Set connection to true
		board.Connections[pair.A][pair.B] = true
		board.Connections[pair.A][pair.B] = true
		// Add neighbors to Tiles
		board.Tiles[pair.A].AddNeighbor(board.Tiles[pair.B])
		board.Tiles[pair.B].AddNeighbor(board.Tiles[pair.A])
		// Create Links
		board.Links = append(board.Links, NewLink(id, board.Tiles[pair.A], board.Tiles[pair.B]))
	}

	logrus.Debugf("Created empty board with %s tiles", numTiles)

	return &board
}

func LoadExampleBoard(players []int) *Board {
	b := NewEmptyBoard(6, []Connection{
		{0, 1},
		{1, 2},
		{2, 3},
		{3, 4},
		{4, 5},
	})
	if len(players) == 1 {
		b.Tiles[2].Owner = players[0]
		return b
	}
	if len(players) == 2 {
		b.Tiles[0].Owner = players[0]
		b.Tiles[4].Owner = players[1]
		return b
	}

	for i, p := range players {
		b.Tiles[i].Owner = p
	}
	return b
}
