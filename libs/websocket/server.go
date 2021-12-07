package websocket

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	gorillaWebsocket "github.com/gorilla/websocket"
	log "github.com/hop-/golog"
)

const (
	pingPeriod = 30 * time.Second
)

var (
	upgrader = gorillaWebsocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	handlerFuncs = make(map[string]func(*Connection))
)

// SetHandlerFunc set websocket handler function
func SetHandlerFunc(endpoint string, h func(*Connection)) {
	handlerFuncs[endpoint] = h
	http.HandleFunc(endpoint, handleConnection)
}

// Serve serves websocket in given port
func Serve(port int) {
	// TODO add handler for "/"
	wsPort := fmt.Sprintf(":%d", port)
	// Listen and serve and upgrade http connections
	http.ListenAndServe(wsPort, nil)
}

func getIP(r *http.Request) string {
	lastProxy := strings.Split(r.RemoteAddr, ":")[0]
	xForward := r.Header.Get("X-Forwarded-For")
	xReal := r.Header.Get("X-Real-IP")

	if xForward != "" {
		return strings.Split(xForward, ", ")[0]
	}
	if xReal != "" {
		return xReal
	}

	return lastProxy
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("upgrade") == "" {
		log.Info("Upgrade header is missing: responsing as regular request")
		w.Write([]byte("Connected"))
		return
	}

	log.Infof("New websocket request")

	ipAddress := getIP(r)

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Infof("Accepted from %s", ws.RemoteAddr())

	conn := NewConnection(ws, ipAddress)
	defer conn.Close()

	go heartbeat(conn)
	handlerFuncs[r.URL.Path](conn)
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
