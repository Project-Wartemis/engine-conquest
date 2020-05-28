package message

// Sent when registering with Backend/Room
type RegisterMessage struct {
	Message
	ClientType string `json:"clientType"`
	Name       string `json:"name"`
}

// Sent when a new gamestate is available
type StateMessage struct {
	Message
	State interface{} `json:"state"`
	Turn  int         `json:"turn"`
}

// Think bout it later
type EndGameMessage struct {
	Message
	// State string `json:"state"`
}

func NewStateMessage(state interface{}) StateMessage {
	message := Message{
		Type: "state",
	}
	return StateMessage{
		Message: message,
		State:   state,
		Turn:    0,
	}
}

func NewRegisterMessage(name string) RegisterMessage {
	message := Message{
		Type: "register",
	}
	return RegisterMessage{
		Message:    message,
		ClientType: "engine",
		Name:       name,
	}
}
