package message

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type Message struct {
	Type string `json:"type"`
}

// Sent when connected successfully
type ConnectedMessage struct {
	Message
}

// Sent when registration was complete
type RegisterSuccessMessage struct {
	Message
	Id int `json:"id"`
}

// Sent when game starts
type StartMessage struct {
	Message
	Players []int `json:"players"`
}

// Sent by player on move
type ActionMessage struct { // also outgoing
	Message
	Player int             `json:"player"`
	Action json.RawMessage `json:"action"`
}

// Sent when invited to a game
type InviteMessage struct {
	Message
	Client int `json:"client"`
	Room   int `json:"room"`
}

// Sent when error occured
type ErrorMessage struct {
	Message
	Content string `json:"message"`
}

// Parse general message
func ParseMessage(raw []byte) (*Message, error) {
	message := &Message{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		logrus.Warnf("Could not parse Message [%s]", raw)
		return nil, err
	}
	return message, nil
}

func ParseRegisterSuccessMessage(raw []byte) (*RegisterSuccessMessage, error) {
	message := &RegisterSuccessMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		logrus.Warnf("Could not parse RegisterSuccessMessage [%s]", raw)
		return nil, err
	}
	return message, nil
}

func ParseConnectedMessage(raw []byte) (*ConnectedMessage, error) {
	message := &ConnectedMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		logrus.Warnf("Could not parse ConnectedMessage [%s]", raw)
		return nil, err
	}
	return message, nil
}

func ParseInviteMessage(raw []byte) (*InviteMessage, error) {
	message := &InviteMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		logrus.Warnf("Could not parse InviteMessage [%s]", raw)
		return nil, err
	}
	return message, nil
}

func ParseStartMessage(raw []byte) (*StartMessage, error) {
	message := &StartMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		logrus.Warnf("Could not parse StartMessage [%s]", raw)
		return nil, err
	}
	return message, nil
}

func ParseActionMessage(raw []byte) (*ActionMessage, error) {
	message := &ActionMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		logrus.Warnf("Could not parse ActionMessage [%s]", raw)
		return nil, err
	}
	return message, nil
}

func ParseErrorMessage(raw []byte) (*ErrorMessage, error) {
	message := &ErrorMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		logrus.Warnf("Could not parse ErrorMessage [%s]", raw)
		return nil, err
	}
	return message, nil
}
