package websocket

import (
	"fmt"
	"net/http"
	"time"

	"../log"

	gorillaWebsocket "github.com/gorilla/websocket"
)

const (
	pingPeriod = 30 * time.Second
)

var (
	upgrader = gorillaWebsocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	handlerFunc func(*Connection)
)

func handleConnection(w http.ResponseWriter, r *http.Request) {
	log.Infof("New websocket request")

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Infof("Accepted from %s", ws.RemoteAddr())

	conn := NewConnection(ws)
	defer conn.Close()

	go heartbeat(conn)
	handlerFunc(conn)
}

func heartbeat(conn *Connection) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-conn.close:
			return
		case <-ticker.C:
			if err := conn.ping(); err != nil {
				defer conn.Close()
				log.Error(err.Error())
				return
			}
		}
	}
}

// SetHandlerFunc set websocket handler function
func SetHandlerFunc(h func(*Connection)) {
	handlerFunc = h
}

// Serve serves websocket in given port
func Serve(port int) {
	http.HandleFunc("/", handleConnection)

	wsPort := fmt.Sprintf(":%d", port)
	// Listen and serve and upgrade http connections
	http.ListenAndServe(wsPort, nil)
}
