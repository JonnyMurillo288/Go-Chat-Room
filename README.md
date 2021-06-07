# Go-Chat-Room

### Server
- Server object holds all the rooms, clients, for connected to it for easier implementation of 
clients requesting to leave the room and going back to the home page

- When the server is running contains channels (<- chan) that register and unregister clients into the room websocket

- This loop also contains the brodcast channel to broadcast into the websockets 

### Client
- Client object holds the name, room, http.Client

##### http.Client v websocket.Conn
- The reason for a client object holding both http.Client and websocket.Conn is to ensure that the client 
can connect to the http connection at port 8329 

- websocket.Conn allows for the client to open up a websocket when joining a room to get fast receival and 
sending of messages in a single port on their machine

- send channel sends the message to the server where brodcast to be sent to the appropriate websocket of the other users in the room

### Rooms
- Room object doesnt hold much, but the name and the members in the room
- Holding these object will allow for easier control over the serer and allow for more functionality in the future,
currently not too useful

