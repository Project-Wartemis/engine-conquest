package communication

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type BackendConnection struct {
	lock       sync.Mutex
	addr       string
	endpoint   string
	clientId   int
	clientName string
	RoomId     int
	socket     websocket.Conn
	Master     GameInterface
}

func (bc *BackendConnection) GetAddr() string {
	return bc.addr
}

func NewBackendConnection(addr string, endpoint string, clientName string) *BackendConnection {
	return &BackendConnection{
		addr:       addr,
		endpoint:   endpoint,
		clientName: clientName,
	}
}

func (bc *BackendConnection) ConnectToBackend() (*websocket.Conn, error) {

	backendEndpoint := url.URL{Scheme: "ws", Host: bc.addr, Path: bc.endpoint}
	logrus.Infof("Attempting to connect to backend... [%s]", backendEndpoint.String())

	con, resp, err := websocket.DefaultDialer.Dial(backendEndpoint.String(), nil)
	if err != nil {
		respStatusCode := "no response"
		if resp != nil {
			respStatusCode = fmt.Sprintf("%d", resp.StatusCode)
		}
		logrus.Warnf("Failed to connect to backend: [%s] - [%s]", respStatusCode, err)
		return nil, err
	}

	logrus.Infof("Connection with Backend established. Listening for messages...")
	bc.socket = *con
	for {
		err = bc.ListenForMessages()
		if err != nil {
			logrus.Errorf("Listinng for messages failed: [%s]", err)
			time.Sleep(200 * time.Millisecond)
		}
	}
}

// UTILS

func (bc *BackendConnection) SendMessage(message interface{}) error {
	// logrus.Debugf("Sending message: [%s]", message)
	text, err := json.Marshal(message)
	if err != nil {
		logrus.Errorf("Unexpected error while marshalling message to json: [%s]", message)
		return err
	}
	bc.lock.Lock()
	defer bc.lock.Unlock()
	err = bc.socket.WriteMessage(websocket.TextMessage, text)
	if err != nil {
		logrus.Errorf("Unexpected error while sending message: [%s]", message)
		return err
	}
	return nil
}

func (bc *BackendConnection) readMessage() ([]byte, error) {
	_, message, err := bc.socket.ReadMessage()
	if err != nil {
		logrus.Errorf("Unable to read message: [%s]", err)
		return nil, err
	}
	return message, nil
}
