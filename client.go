package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)


const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}


type client struct {
	Name string
	Room *room
	Cli *http.Client
	// added for the websocket implementation
	Server *server //
	conn *websocket.Conn //
	send chan []byte //
}

//clients struct hold http.client obj
var clients []struct{}

func (s *server) newClient() *client {
	fmt.Println("Created new client??")
	return &client{
		Cli: &http.Client{
			Timeout: 20 * time.Second,
			Transport: &http.Transport{
				DisableKeepAlives: false,
				Proxy:             http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				MaxIdleConns:          100,
				IdleConnTimeout:       50 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
		},},
		Server: s,
		conn: &websocket.Conn{},
		send: make(chan []byte),
		
	}
}

func makeRequest(c *client, url string) (int,error) {
	resp, err := c.Cli.Get(url)
	if err != nil {
		return 0,err
	}
	defer resp.Body.Close()
	return resp.StatusCode,nil
}

func (c *client) login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := templates.Templates["newClient"]
		t.Execute(w,nil)
	} else {
		r.ParseForm()
		fmt.Println("Client Name:",r.Form["username"])
		c.Name = r.Form["username"][0]
	}
	http.Redirect(w,r,joinURL,http.StatusFound)
}

func (c *client) quitCurrentRoom() {
	c.Room = nil
	c.conn = &websocket.Conn{}
	makeRequest(c,joinURL)
}

//reads messages pumped from the websocket connection to the server
func (c *client) readPump() {
	defer func() {
		c.Server.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { 
		c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil
	})
	for {
		_,message,err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error with closing socket: %s",err.Error())
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message,newline,space,-1))
		c.Server.broadcast <- message
	}
}

func (c *client) writePump() {
	prefix := []byte(c.Name + ":")
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message,ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// if c.send channel is closed, close the socket
				c.conn.WriteMessage(websocket.CloseMessage,[]byte{})
				return
			}
			
			// prepare the writer by calling nextwriter
			w,err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("Error with writing the message: %s",err.Error())
				return 
			}
			message = append(prefix, message...) // write the prfix so we know what user sent what message
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage,nil); err != nil {
				return
			}
		}
	}
}

// handles websocket requests from the peer
// pass in existing client from their login and the room they are joining 
func (c *client) serveWs(s *server, w http.ResponseWriter, r *http.Request) {
	conn,err := upgrader.Upgrade(w,r,nil)
	if err != nil {
		log.Println("error upgrading to websocket:",err.Error())
		return
	}
	c.Server = s
	c.conn = conn
	c.send = make(chan[]byte, 256)

	c.Server.register <- c

	go c.writePump()
	go c.readPump()
}

