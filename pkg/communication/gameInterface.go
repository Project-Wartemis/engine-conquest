package communication

import "github.com/Project-Wartemis/pw-engine/pkg/communication/message"

type GameInterface interface {
	StartGame(message.StartMessage)
	CreateNewGame(message.InviteMessage)
	HandleActionMessage(message.ActionMessage)
	GetRoomId() int
}
