package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader {
    // Allow connections from any origin
    CheckOrigin: func(r *http.Request) bool { return true },    
}

type hub struct {
	messages    chan string
	connections map[int]*websocket.Conn
	addConn     chan *websocket.Conn
}

func (h *hub) sendAll(message string) {
	expired := []int{}
	for i, conn := range h.connections {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			conn.Close()
			expired = append(expired, i)
		}
	}
	// Prune the obsolete connections
	if len(expired) > 0 {
		for _, connId := range expired {
			log.Println("Closed connection:", connId)
			delete(h.connections, connId)
		}
	}
}

func (h *hub) Run() {
	id := 0
	for {
		select {
		// Client has connected
		case c := <-h.addConn:
			log.Println("New connection:", id)
			h.connections[id] = c
			id++
		// A new message has been received
		case c := <-h.messages:
			h.sendAll(c)
		}
	}
}

// Initialize our main hub struct
var h = &hub{
	messages:    make(chan string),
	connections: make(map[int]*websocket.Conn),
	addConn:     make(chan *websocket.Conn),
}

// Receives messages line by line from stdin and sends them through the messages
// channel.
func receive() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		h.messages <- scanner.Text()
	}
}

// Handle upgrades to websocket
func connectionHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	h.addConn <- c
}

func main() {

	portNo := flag.Int("port", 8080, "the port for the websocket server")
	flag.Parse()

	go receive()

	go h.Run()

	log.Printf("Serving on %d...", *portNo)
	http.HandleFunc("/", connectionHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", *portNo), nil)
}
