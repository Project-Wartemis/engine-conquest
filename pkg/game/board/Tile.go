package board

type Tile struct {
	Id        int
	Owner     int
	Garrison  int
	Neighbors []*Tile
}

func NewTile(id int) *Tile {
	return &Tile{
		Id:        id,
		Owner:     -1,
		Garrison:  0,
		Neighbors: []*Tile{},
	}
}

func (t *Tile) AddNeighbor(n *Tile) {
	t.Neighbors = append(t.Neighbors, n)
}
