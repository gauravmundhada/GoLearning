package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

type Server struct {
	mu        sync.Mutex
	conns     map[*websocket.Conn]bool
	connModes map[*websocket.Conn]string
}

func NewServer() *Server {
	return &Server{
		conns:     make(map[*websocket.Conn]bool),
		connModes: make(map[*websocket.Conn]string),
	}
}

// startDataFeed starts a goroutine that sends data feed messages to all clients in "datafeed" mode every 3 seconds
func (s *Server) startDataFeed() {
	go func() {
		for {
			payload := fmt.Sprintf("Data feed - %s", time.Now().Format(time.RFC1123))

			s.mu.Lock()
			for ws, mode := range s.connModes {
				if mode == "datafeed" {
					go func(ws *websocket.Conn) {
						if _, err := ws.Write([]byte(payload)); err != nil {
							fmt.Println("Error sending data feed:", err)
						}
					}(ws)
				}
			}
			s.mu.Unlock()

			time.Sleep(3 * time.Second)
		}
	}()
}

// func (s *Server) handleDataFeed(ws *websocket.Conn) {
// 	fmt.Println("New incoming client connection to the data feed", ws.RemoteAddr())

// 	for {
// 		payload := fmt.Sprintf("Data feed - %d\n", time.Now().UnixNano())

// 		s.mu.Lock()
// 		for sub := range s.datafeedConns {
// 			go func(sub *websocket.Conn) {
// 				if _, err := sub.Write([]byte(payload)); err != nil {
// 					fmt.Println("Error writing to subscriber:", err)
// 				}
// 			}(sub)
// 		}
// 		s.mu.Unlock()

// 		time.Sleep(3 * time.Second)
// 	}
// }

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("New incoming client connection", ws.RemoteAddr())
	s.mu.Lock()
	s.conns[ws] = true // mark the connection as active
	s.mu.Unlock()

	s.readLoop(ws)
}

// readLoop reads messages from the WebSocket connection
func (s *Server) readLoop(ws *websocket.Conn) {
	buff := make([]byte, 1024) // size of bytes
	for {
		n, err := ws.Read(buff)
		if err != nil {
			if err == io.EOF { // client has closed the connection
				break
			}
			fmt.Println("Error in reading from connection", err)
			continue // continue to read next message(if want then we can even drop the connection)
		}
		msg := buff[:n] // dont read the entire buffer (only read the bytes which are read into buff)
		//fmt.Println(string(msg))

		msgStr := string(msg)
		// switch between modes based on client message
		switch msgStr {
		case "mode:chat":
			s.mu.Lock()
			s.connModes[ws] = "chat"
			s.mu.Unlock()
			ws.Write([]byte("Switched to chat mode"))

		case "mode:datafeed":
			s.mu.Lock()
			s.connModes[ws] = "datafeed"
			s.mu.Unlock()
			ws.Write([]byte("Switched to data feed mode"))

		default:
			s.mu.Lock()
			mode := s.connModes[ws]
			s.mu.Unlock()

			if mode == "chat" {
				s.broadcast(msg)
			}
		}

		// we can even reply with something to the client - check the below code
		//ws.Write([]byte("Thank you for your message!"))
		//s.broadcast(msg)
	}
}

func (s *Server) broadcast(msg []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// send the message to all connected clients
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(msg); err != nil {
				fmt.Println("Error in writing to connection", err)
			}
			fmt.Println("Sent message to", ws.RemoteAddr())
		}(ws)
	}
}

// Entry point of the application
func main() {
	server := NewServer()
	server.startDataFeed()

	http.Handle("/ws", websocket.Handler(server.handleWS))
	//http.Handle("/datafeed", websocket.Handler(server.handleDataFeed))
	http.ListenAndServe(":3000", nil)
}
