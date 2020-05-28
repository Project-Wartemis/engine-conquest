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
	b := NewEmptyBoard(5, []Connection{
		{0, 1},
		{0, 2},
		{0, 4},
		{1, 2},
		{2, 3},
		{3, 4},
	})
	for i, p := range players {
		b.Tiles[i].Owner = p
	}
	return b
}
