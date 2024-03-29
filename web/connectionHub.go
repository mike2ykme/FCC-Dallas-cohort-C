package web

import (
	"errors"
	"fmt"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"strconv"
	"strings"
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

	closeConnectionWithMessage := func(message string) {
		connection.Logger.Println(message)
		_ = connection.Connection.WriteMessage(websocket.CloseMessage, []byte{})
		err := connection.Connection.Close()
		if err != nil {
			connection.Logger.Printf("there was an error closing out a connection: %#v\n", connection)
		}
	}

	if !keyExists {
		closeConnectionWithMessage("there is no room with key: " + connection.RoomId.String())
		return
	}

	if !entry.Joinable {
		closeConnectionWithMessage("A game is already in progress in this room.")
		return
	}

	entry.Connections[connection.Connection] = models.Client{}
	entry.ConnectedUsers[connection.UserId] = connection.Username

	rooms[connection.RoomId] = entry

	// The user is only an admin, if they're the first person there.
	// And if we already have a channel, then they're not the first person
	connectMessage := models.InitialConnection{
		MessageType: "initial-connection",
		Action: "REGISTERED",
		Admin:  connection.UserId == entry.AdminId, //!keyExists,
	}
	if err := connection.Connection.WriteJSON(connectMessage); err != nil {
		Configs.Logger.Println(err)
	}

	// broadcast usernames so frontend can show connected users in waiting room
	joinedMessage := models.UserConnectedMessage{
		MessageType: "user-joined",
		Contents: entry.GetConnectedList(),
	}
	if errMap, count := entry.WriteJsonToAllConnections(joinedMessage); count > 0 {
		for conn, err := range errMap {
			Configs.Logger.Println(err.Error())

			_ = conn.WriteMessage(websocket.CloseMessage, []byte{})
			closeErr := connection.Connection.Close()

			if closeErr != nil {
				Configs.Logger.Printf("there was an error closing out a connection: %#v\n", conn)
			}
		}
	}

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
	message.Logger.Println("broadcast received:", message)
	room := rooms[message.RoomId]
	switch strings.ToUpper(message.UserMessage.Action) {
	case LOAD:
		room.Joinable = false
		// since room is a copy, we have to assign it back to the map
		rooms[message.RoomId] = room
		if message.UserId == room.AdminId {
			err := handleAdminLoad(message.RoomId, message.UserMessage.DeckId, message.Conn)
			if err != nil {
				Configs.Logger.Printf("there was an error: %#v", err)
				_ = message.Conn.WriteJSON(models.NewErrorResponse(err.Error()))
			}
		} else {
			Configs.Logger.Printf("%d is trying to access admin ", message.UserId)
			_ = message.Conn.WriteJSON(models.NewErrorResponse("non admin user"))
		}

	case TEXT:
		handleTextMessage(message)
	case SUBMIT:
		if err := handleAnswerSubmissions(message); err != nil {
			Configs.Logger.Printf("there was an error: %#v", err)
			_ = message.Conn.WriteJSON(models.NewErrorResponse(err.Error()))
		}
	case GETRESULTS:
		if err := returnAllResults(message); err != nil {
			Configs.Logger.Printf("there was an error: %#v", err)
			_ = message.Conn.WriteJSON(models.NewErrorResponse(err.Error()))
		}
	}

}

func returnAllResults(message models.UserResponse) error {
	roomID := message.RoomId
	room := rooms[roomID]
	//results := rooms[roomID].Results
	results := make(models.UserResults, 0)
	for userId, userScore := range room.Results {
		results[room.ConnectedUsers[userId]] = userScore
	}

	return message.Conn.WriteJSON(models.Results{
		RoomId:      message.RoomId,
		UserResults: results,
	})
}

func handleAnswerSubmissions(message models.UserResponse) error {
	room := rooms[message.RoomId]

	result, err := strconv.Atoi(message.UserMessage.Message)
	if err != nil {
		return errors.New(fmt.Sprintf("invalid score submitted userID: %d \n", message.UserId))
	}

	if room.TotalQuestions < result {
		return errors.New("cannot have a higher score than available number of questions")
	}

	room.Results[message.UserId] = uint(result)

	for conn := range room.Connections {
		_ = conn.WriteJSON(models.NewResult(message.UserId, room.ConnectedUsers[message.UserId], result, room.TotalQuestions))
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
		_ = conn.WriteJSON(models.DefaultErrorResponse())
		return err
	}

	room.Deck = newDeck
	room.TotalQuestions = len(newDeck.FlashCards)
	rooms[roomId] = room
	//	TODO SHUFFLE function

	if errMap, count := room.WriteJsonToAllConnections(models.NewLoadDeck(newDeck.FlashCards)); count > 0 {
		for c, e := range errMap {
			Configs.Logger.Printf("there was an error writing to connection: %#v -> err: %s \n", c, e.Error())
		}
		return errors.New("there was an error writing deck to all connections")

	}

	return nil
}

func handleTextMessage(message models.UserResponse) {
	for connection := range rooms[message.ChannelId].Connections {
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
	entry.ConnectedUsers = make(map[uint]string, 0)
	entry.Connections = make(map[*websocket.Conn]models.Client)
    entry.Joinable = true
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
