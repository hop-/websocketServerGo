package websocket

import (
	"encoding/json"
	"time"

	gorillaWebsocket "github.com/gorilla/websocket"
	"gopkg.in/go-playground/validator.v9"
)

const (
	writeTimeout = 60 * time.Second
)

var (
	val = validator.New()
)

// Connection structure
type Connection struct {
	Conn  *gorillaWebsocket.Conn
	close chan bool
}

// CloseError gorilla websocket CloseError alias
type CloseError = gorillaWebsocket.CloseError

// NewConnection create a new connection object
func NewConnection(ws *gorillaWebsocket.Conn) *Connection {
	return &Connection{ws, make(chan bool, 1)}
}

// Write response
func (conn *Connection) Write(response Response) error {
	bytes, err := json.Marshal(response)
	if err != nil {
		return err
	}

	conn.Conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	return conn.Conn.WriteMessage(gorillaWebsocket.TextMessage, bytes)
}

// Read request
func (conn *Connection) Read() (*Request, error) {
	request := Request{}
	err := conn.Conn.ReadJSON(&request)
	if err != nil {
		return &request, err
	}

	if err = val.Struct(request); err != nil {
		return nil, err
	}

	request.payload, err = json.Marshal(request.Payload)
	if err != nil {
		return &request, err
	}

	return &request, err
}

// ping to client
func (conn *Connection) ping() error {
	conn.Conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	return conn.Conn.WriteMessage(gorillaWebsocket.PingMessage, []byte{})
}

// Shutdown gracefully
func (conn *Connection) Shutdown() {
	conn.Conn.WriteMessage(gorillaWebsocket.CloseMessage, []byte{})
	conn.Close()
}

// Close connection
func (conn *Connection) Close() {
	conn.close <- true
	conn.Conn.Close()
}
