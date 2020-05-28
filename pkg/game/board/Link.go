package board

type Link struct {
	Id int
	A  *Tile
	B  *Tile
}

func NewLink(id int, A *Tile, B *Tile) Link {
	return Link{
		Id: id,
		A:  A,
		B:  B,
	}
}
