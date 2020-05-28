package communication

import (
	"github.com/Project-Wartemis/pw-engine/pkg/communication/message"
	"github.com/sirupsen/logrus"
)

func (bc *BackendConnection) ListenForMessages() error {
	for {
		message, err := bc.readMessage()
		if err != nil {
			return err
		}
		go bc.HandleIncommingMessage(message)
	}
}

func (bc *BackendConnection) HandleIncommingMessage(raw []byte) {
	logrus.Infof("Received message: [%s]", raw)
	msg, err := message.ParseMessage(raw)
	if err != nil {
		logrus.Errorf("Could not parse message: [%s]", raw)
		return
	}
	handler := bc.handleDefault
	switch msg.Type {
	case "connected":
		handler = bc.handleConnectedMessage
	case "register":
		handler = bc.handleRegisterSuccessMessage
	case "invite":
		handler = bc.handleInviteMessage
	case "start":
		handler = bc.handleStartMessage
	case "action":
		handler = bc.handleActionMessage
	}
	handler(raw)
}

func (bc *BackendConnection) handleDefault(raw []byte) {}
func (bc *BackendConnection) handleActionMessage(raw []byte) {
	msg, err := message.ParseActionMessage(raw)
	if err != nil {
		return
	}
	bc.Master.HandleActionMessage(*msg)
}

func (bc *BackendConnection) handleStartMessage(raw []byte) {
	msg, err := message.ParseStartMessage(raw)
	if err != nil {
		return
	}
	logrus.Infof("Start game for Room %d with players [%s]", bc.RoomId, msg.Players)
	bc.Master.StartGame(*msg)
}

func (bc *BackendConnection) handleConnectedMessage(raw []byte) {
	logrus.Info("Connected! Sending register request")
	regMsg := message.NewRegisterMessage(bc.clientName)
	bc.SendMessage(regMsg)
}

func (bc *BackendConnection) handleRegisterSuccessMessage(raw []byte) {
	msg, err := message.ParseRegisterSuccessMessage(raw)
	if err != nil {
		return
	}
	logrus.Debugf("Registered! Got ID %d", msg.Id)
	if bc.clientId > 0 {
		logrus.Warningf("Already registered with ID %d. Ignoring new ID %d", bc.clientId, msg.Id)
		return
	}
	bc.clientId = msg.Id
}

func (bc *BackendConnection) handleInviteMessage(raw []byte) {
	msg, err := message.ParseInviteMessage(raw)
	if err != nil {
		logrus.Warnf("Could not parse ConnectedMessage: [%s]", raw)
	}

	bc.Master.CreateNewGame(*msg)
}
