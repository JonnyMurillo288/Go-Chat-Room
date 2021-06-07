package main

import (
	"fmt"
)

type server struct {
	Rooms map[string]*room
	Clients map[*client]bool
	broadcast chan []byte
	register chan *client
	unregister chan *client

}

func newServer() *server {
	fmt.Println("Created a server")
	return &server{
		Rooms: make(map[string]*room),
		Clients: make(map[*client]bool),
		broadcast: make(chan []byte),
		register: make(chan *client),
		unregister: make(chan *client),
	}
}


// runs the server to accept of send messgages and closes out the client socket
func (s *server) run() {
	for {
		select {
		case client := <-s.register:
			s.Clients[client] = true
		case client := <-s.unregister:
			if _,ok := s.Clients[client]; ok {
				delete(s.Clients,client)
				close(client.send)
			}
		case message := <-s.broadcast:
		for client := range s.Clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(s.Clients,client)
			}
		}
			
		}
	}
}
