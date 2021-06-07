package main

import (
	"fmt"
	"log"
	"net/http"
)

type room struct {
	Name string
	Members map[string]*client
}


func (c *client) joinRoom(roomName string) {
	fmt.Println("Room Name:",roomName)
	room,ok := c.Server.Rooms[roomName]
	if ok == false {
		log.Printf("No Room, Have to create: %s",roomName)
		c.Server.Rooms[roomName] = c.newRoom(roomName)
	}
	fmt.Println(c.Server.Rooms)
	room = c.Server.Rooms[roomName]
	c.Room = room
	room.Members[c.Name] = c
}

func (c *client) newRoom(roomName string) *room {
	room :=  &room{
		Name: roomName,
		Members: make(map[string]*client),
	}
	room.Members[c.Name] = c
	return room
}


func (s *server) getRooms() map[string]*room {
	return s.Rooms
}

func (c *client) displayRooms(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := templates.Templates["index"]
		rooms := c.Server.getRooms()
		t.Execute(w,rooms)
	} else {
		r.ParseForm()
		c.joinRoom(r.Form["Room"][0])
	}
	// http.Redirect(w,r,wsURL,http.StatusFound)
	http.Redirect(w,r,roomURL,http.StatusFound)
}


