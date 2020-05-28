package game

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/Project-Wartemis/pw-engine/pkg/communication"
	"github.com/Project-Wartemis/pw-engine/pkg/communication/message"
	"github.com/Project-Wartemis/pw-engine/pkg/game/board"
	"github.com/Project-Wartemis/pw-engine/pkg/game/events"
	"github.com/sirupsen/logrus"
)

type Game struct {
	Lock                 sync.Mutex
	connection           *communication.BackendConnection
	Gamestate            *Gamestate
	WaitingOnPlayers     map[int]struct{} // Sets don't exist
	pendingDeployActions []events.DeployAction
	pendingMoveActions   []events.MoveAction
}

func NewGame(addr string, endpoint string, engineName string) *Game {
	g := &Game{}
	g.connection = communication.NewBackendConnection(addr, endpoint, engineName)
	g.connection.Master = g
	g.Gamestate = &Gamestate{}
	return g
}

func (g *Game) OpenConnection() {
	for {
		g.connection.ConnectToBackend()
		logrus.Errorf("Game [%d] Lost connection to backend. Reconnecting...", g.connection.RoomId)
		time.Sleep(200 * time.Millisecond)
	}
}

// Implement GameInterface
func (g *Game) StartGame(msg message.StartMessage) {
	g.Gamestate.Players = []Player{}
	for _, p := range msg.Players {
		g.Gamestate.Players = append(g.Gamestate.Players, Player{
			Name:           p,
			ArmiesToDeploy: 5,
		})
	}

	players := []int{}
	for _, p := range g.Gamestate.Players {
		players = append(players, p.Name)
		g.Gamestate.Board = *board.LoadExampleBoard(players)
	}

	g.resetWaitingList()

	logrus.Info("Sending initial Gamestate")
	g.sendInitialGamestate()

}

func (g *Game) sendInitialGamestate() {
	stages := g.Gamestate.ProcessTurn([]events.DeployAction{}, []events.MoveAction{})

	*g.Gamestate = stages.End

	toExport := stages.Export()

	g.connection.SendMessage(message.NewStateMessage(toExport))
}

func (g *Game) HandleActionMessage(msg message.ActionMessage) {
	g.Lock.Lock()
	defer g.Lock.Unlock()
	// Check if this action is timely
	if _, ok := g.WaitingOnPlayers[msg.Player]; !ok {
		logrus.Warningf("Unexpectedly received move from player %d", msg.Player)
	}
	delete(g.WaitingOnPlayers, msg.Player)

	// TODO: Check if action is valid

	// Handle action
	g.ParsePlayerAction(msg)

	// Received all actions?
	if len(g.WaitingOnPlayers) > 0 {
		logrus.Debugf("Still waiting for players")
		// Nope, still waiting
		return
	}

	// Calculate the new gamestate
	stages := g.Gamestate.ProcessTurn(g.pendingDeployActions, g.pendingMoveActions)
	g.resetActions()

	*g.Gamestate = stages.End

	toExport := stages.Export()

	// Send response if necessary
	logrus.Debug("Send new gamestame")
	g.connection.SendMessage(message.NewStateMessage(toExport))

	// Wait for all players
	g.resetWaitingList()
}

func (g *Game) resetActions() {
	g.pendingDeployActions = []events.DeployAction{}
	g.pendingMoveActions = []events.MoveAction{}
}

func (g *Game) resetWaitingList() {
	g.WaitingOnPlayers = map[int]struct{}{}
	for _, p := range g.Gamestate.Players {
		g.WaitingOnPlayers[p.Name] = struct{}{}
	}
}

type PlayerAction struct {
	Deploys []events.DeployAction `json:"deploys"`
	Moves   []events.MoveAction   `json:"moves"`
}

func (g *Game) ParsePlayerAction(msg message.ActionMessage) {
	action := &PlayerAction{}
	err := json.Unmarshal(msg.Action, action)
	if err != nil {
		logrus.Errorf("Failed to parse action message: %s due to [%s]", msg, err)
	}

	for _, deployAction := range action.Deploys {
		g.pendingDeployActions = append(g.pendingDeployActions, deployAction)
	}
	for _, moveAction := range action.Moves {
		g.pendingMoveActions = append(g.pendingMoveActions, moveAction)
	}
}

func (g *Game) CreateNewGame(msg message.InviteMessage) {
	logrus.Warn("Can't create a new game IN a game...")
}
func (g *Game) GetRoomId() int { return g.connection.RoomId }
