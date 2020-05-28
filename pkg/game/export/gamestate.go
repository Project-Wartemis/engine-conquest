package export

type Gamestate struct {
	Nodes   []Node   `json:"nodes"`
	Links   []Link   `json:"links"`
	Players []Player `json:"players"`
	Deploys []Deploy `json:"deploys"`
	Moves   []Move   `json:"moves"`
	Fights  Fights   `json:"fights"`
	Stages  Stages   `json:"stages"`
}

type Player struct {
	Id     int `json:"id"`
	Income int `json:"income"`
}

type Link struct {
	Id int `json:"id"`
	A  int `json:"a"`
	B  int `json:"b"`
}

type Node struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Stages struct {
	Start  []NodeState `json:"start"`
	Travel []NodeState `json:"travel"`
	End    []NodeState `json:"end"`
}

type NodeState struct {
	Id     int `json:"id"`
	Owner  int `json:"owner"`
	Troops int `json:"troops"`
}

type Deploy struct {
	Node   int `json:"node"`
	Troops int `json:"troops"`
}

type Move struct {
	Source int `json:"source"`
	Target int `json:"target"`
	Troops int `json:"troops"`
}

type Fights struct {
	Travel []Fight `json:"travel"`
	Combat []Fight `json:"combat"`
}

type Fight struct {
	Location int    `json:"location"`
	Armies   []Army `json:"armies"`
}

type Army struct {
	Player int `json:"player"`
	Troops int `json:"troops"`
}
