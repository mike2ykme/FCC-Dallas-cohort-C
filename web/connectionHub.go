package web

import (
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"log"
	"teamC/models"
)

// These should only be accessed via the runHub's goroutine so we can ensure there are no data issues.
// this way we also don't have to have any locks on any data
//var rooms = make(map[uint64]models.Room)
var rooms = make(map[uuid.UUID]models.Room)
var register = make(chan *models.UserConnection)
var unregister = make(chan *models.UserConnection)
var broadcast = make(chan models.UserResponse)
var newRoom = make(chan models.RoomCreation)

func handleRegistration(connection *models.UserConnection) {
	//https://stackoverflow.com/questions/42605337/cannot-assign-to-struct-field-in-a-map
	entry, keyExists := rooms[connection.RoomId]

	if !keyExists {
		connection.Logger.Println("there is no room with key: " + connection.RoomId.String())
		_ = connection.Connection.WriteMessage(websocket.CloseMessage, []byte{})
		if err := connection.Connection.Close(); err != nil {
			connection.Logger.Printf("there was an error closing out a connection: %#v\n", connection)
		}
	}

	if entry.Connections == nil {
		entry.Connections = make(map[*websocket.Conn]models.Client)
	}

	entry.Connections[connection.Connection] = models.Client{}
	rooms[connection.RoomId] = entry

	// The user is only an admin, if they're the first person there.
	// And if we already have a channel, then they're not the first person
	connectMessage := models.InitialConnection{
		Action: "REGISTERED",
		Admin:  connection.UserId == entry.AdminId, //!keyExists,
	}
	_ = connection.Connection.WriteJSON(connectMessage)

	log.Println("Connection registered to new room")
}

// We need to handle the state of the rooms and the q and a's
func handleBroadcast(message models.UserResponse) {
	log.Println("message received:", message)
	for connection := range rooms[message.ChannelId].Connections {
		if err := connection.WriteJSON(message); err != nil {
			_ = connection.WriteMessage(websocket.CloseMessage, []byte{})
			_ = connection.Close()
			delete(rooms[message.ChannelId].Connections, connection)
		}
	}
}

func handleUnregister(connection *models.UserConnection) {
	// Remove the client from the hub
	delete(rooms[connection.RoomId].Connections, connection.Connection)
	if len(rooms[connection.RoomId].Connections) < 1 {
		delete(rooms, connection.RoomId)
	}
	log.Println("connection unregistered")
}
func handleNewRoom(roomSetup models.RoomCreation) {
	entry, keyExists := rooms[roomSetup.NewRoomID]
	if keyExists {
		roomSetup.Logger.Println("There was an existing room with key: " + roomSetup.NewRoomID.String())
		return
	}
	entry.AdminId = roomSetup.AdminId
	rooms[roomSetup.NewRoomID] = entry

	roomSetup.Logger.Println("successfully setup a new room")
}
func RunHub() {
	for {
		select {
		case connection := <-register:
			handleRegistration(connection)

		case message := <-broadcast:
			handleBroadcast(message)

		case connection := <-unregister:
			handleUnregister(connection)

		case room := <-newRoom:
			handleNewRoom(room)
		}

	}
}
