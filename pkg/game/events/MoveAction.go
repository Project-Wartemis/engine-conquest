package events

//  What the player wants to do
type MoveAction struct {
	Player       int `json:"player"`
	NumTroops    int `json:"numTroops"`
	SourceTileId int `json:"sourceTileId"`
	TargetTileId int `json:"targetTileId"`
}
