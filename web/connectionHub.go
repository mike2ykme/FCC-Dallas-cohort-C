package web

import (
	"errors"
	"fmt"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"strconv"
	"teamC/Global"
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

var Configs *Global.Configuration

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

	connection.Logger.Println("Connection registered to new room")
}

const (
	TEXT       = "TEXT"
	SUBMIT     = "SUBMIT"
	LOAD       = "LOAD"
	GETRESULTS = "GETRESULTS"
)

// We need to handle the state of the rooms and the q and a's
func handleBroadcast(message models.UserResponse) {
	message.Logger.Println("message received:", message)
	room := rooms[message.RoomId]
	switch message.UserMessage.Action {
	case LOAD:
		if message.UserId == room.AdminId {
			err := handleAdminLoad(message.RoomId, message.UserMessage.DeckId, message.Conn)
			if err != nil {
				Configs.Logger.Printf("there was an error: %#v", err)
				message.Conn.WriteJSON(models.NewErrorResponse(err.Error()))
			}
		} else {
			message.Conn.WriteJSON(models.NewErrorResponse("non admin user"))
		}

	case TEXT:
		handleTextMessage(message)
	case SUBMIT:
		if err := handleAnswerSubmissions(message); err != nil {
			Configs.Logger.Printf("there was an error: %#v", err)
			message.Conn.WriteJSON(models.NewErrorResponse(err.Error()))
		}
	case GETRESULTS:
		if err := returnAllResults(message); err != nil {
			Configs.Logger.Printf("there was an error: %#v", err)
			message.Conn.WriteJSON(models.NewErrorResponse(err.Error()))
		}
	}

}

func returnAllResults(message models.UserResponse) error {
	roomID := message.RoomId
	results := rooms[roomID].Results

	return message.Conn.WriteJSON(results)

}

func handleAnswerSubmissions(message models.UserResponse) error {
	room := rooms[message.RoomId]
	result, err := strconv.Atoi(message.UserMessage.Message)
	if err == nil {

		room.Results[message.UserId] = uint(result)
	} else {
		return errors.New(fmt.Sprintf("there was an error converting the user's submission result for userID: %s\n", message.UserId))
	}

	for conn, _ := range room.Connections {
		username, err := Configs.UserRepo.GetUsernameById(message.UserId)
		if err != nil {
			return err
		}
		conn.WriteJSON(models.NewResult(message.UserId, username, result))
	}

	return nil
}

func handleAdminLoad(roomId uuid.UUID, deckId int, conn *websocket.Conn) error {
	if deckId <= 0 {
		return errors.New("invalid deck ID")
	}
	room := rooms[roomId]
	newDeck := models.Deck{}
	err := Configs.DeckRepo.GetDeckById(&newDeck, uint(deckId))
	if err != nil {
		Configs.Logger.Printf("there was an error getting the deck %s\n", err.Error())
		conn.WriteJSON(models.DefaultErrorResponse())
		return err
	}

	room.Deck = newDeck
	rooms[roomId] = room
	//	TODO SHUFFLE function

	err = room.WriteJsonToAllConnections(models.NewLoadDeck(newDeck.FlashCards))
	if err != nil {
		Configs.Logger.Printf("there was an error writing deck to all connections, %#v \n", err)
		return err
	}

	return nil
}

func handleTextMessage(message models.UserResponse) {
	for connection, _ := range rooms[message.ChannelId].Connections {
		// We're not catching an error from the user but instead if we have a problem writing to them
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
	Configs.Logger.Println("connection unregistered")
}
func handleNewRoom(roomSetup models.RoomCreation) {
	entry, keyExists := rooms[roomSetup.NewRoomID]
	if keyExists {
		roomSetup.Logger.Println("There was an existing room with key: " + roomSetup.NewRoomID.String())
		return
	}
	entry.AdminId = roomSetup.AdminId
	entry.Results = make(map[uint]uint, 0)

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
