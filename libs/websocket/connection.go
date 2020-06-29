package websocket

import (
	"encoding/json"
	"sync"
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
	Conn       *gorillaWebsocket.Conn
	IP         string
	close      chan bool
	readMutex  sync.Mutex
	writeMutex sync.Mutex
}

// CloseError gorilla websocket CloseError alias
type CloseError = gorillaWebsocket.CloseError

// NewConnection create a new connection object
func NewConnection(ws *gorillaWebsocket.Conn, ip string) *Connection {
	return &Connection{ws, ip, make(chan bool, 1), sync.Mutex{}, sync.Mutex{}}
}

// Write response
func (conn *Connection) Write(response *Response) error {
	bytes, err := json.Marshal(response)
	if err != nil {
		return err
	}

	conn.writeMutex.Lock()
	defer conn.writeMutex.Unlock()

	conn.Conn.SetWriteDeadline(time.Now().Add(writeTimeout))

	return conn.Conn.WriteMessage(gorillaWebsocket.TextMessage, bytes)
}

// Read request
func (conn *Connection) Read() (*Request, error) {
	request := Request{}

	conn.readMutex.Lock()
	err := conn.Conn.ReadJSON(&request)
	conn.readMutex.Unlock()

	if err != nil {
		return &request, err
	}

	if err = val.Struct(request); err != nil {
		return nil, err
	}

	request.options, err = json.Marshal(request.Options)
	if err != nil {
		return &request, err
	}

	return &request, err
}

// ping to client
func (conn *Connection) ping() error {
	conn.writeMutex.Lock()
	defer conn.writeMutex.Unlock()

	conn.Conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	return conn.Conn.WriteMessage(gorillaWebsocket.PingMessage, []byte{})
}

// Shutdown gracefully
func (conn *Connection) Shutdown() {
	conn.writeMutex.Lock()
	defer conn.writeMutex.Unlock()

	conn.Conn.WriteMessage(gorillaWebsocket.CloseMessage, []byte{})
	conn.Close()
}

// Close connection
func (conn *Connection) Close() {
	conn.close <- true
	conn.Conn.Close()
}
