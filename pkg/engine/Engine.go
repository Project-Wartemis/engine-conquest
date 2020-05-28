package engine

import (
	"fmt"
	"time"

	"github.com/Project-Wartemis/pw-engine/pkg/communication"
	"github.com/Project-Wartemis/pw-engine/pkg/communication/message"
	"github.com/Project-Wartemis/pw-engine/pkg/game"
	"github.com/sirupsen/logrus"
)

type Engine struct {
	EngineName      string
	lobbyConnection *communication.BackendConnection
	games           []*game.Game
}

func NewConquestEngine() *Engine {
	return &Engine{EngineName: "Conquest"}
}

func (e *Engine) Start(addr string, lobbyEndpoint string) {
	e.lobbyConnection = communication.NewBackendConnection(addr, lobbyEndpoint, e.EngineName)
	e.lobbyConnection.Master = *e
	// Connect to backend
	for {
		e.lobbyConnection.ConnectToBackend()
		logrus.Error("Lost connection to backend. Reconnecting...")
		time.Sleep(200 * time.Millisecond)
	}
}

// Implement GameInterface
func (e Engine) StartGame(msg message.StartMessage) {
	logrus.Warn("The engine can't be started...")
}
func (e Engine) CreateNewGame(msg message.InviteMessage) {
	roomEndpoint := fmt.Sprintf("/socket/%d", msg.Room)
	newGame := game.NewGame(e.lobbyConnection.GetAddr(), roomEndpoint, e.EngineName)
	e.games = append(e.games, newGame)
	newGame.OpenConnection()
}

func (e Engine) HandleActionMessage(msg message.ActionMessage) {
	logrus.Warn("The engine can't handle action messages...")
}

func (e Engine) GetRoomId() int { return e.lobbyConnection.RoomId }
