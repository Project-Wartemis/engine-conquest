package events

type DeployAction struct {
	Player    int `json:"player"`
	TileId    int `json:"tileId"`
	NumTroops int `json:"numTroops"`
}
