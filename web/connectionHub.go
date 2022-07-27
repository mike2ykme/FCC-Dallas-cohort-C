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

// These should only be accessed via the runHub's goroutine, so we can ensure there are no data issues.
// this way we also don't have to have any locks on any data
//var rooms = make(map[uint64]models.Room)
var rooms = make(map[uuid.UUID]models.Room)
var register = make(chan *models.UserConnection)
var unregister = make(chan *models.UserConnection)
var broadcast = make(chan models.UserResponse)
var newRoom = make(chan models.RoomCreation)

var Configs *Global.Configuration

// Handles Actions

const (
	ActionRegistered = "REGISTERED"
)

// Message Types

const (
	InitialConnection = "initial-connection"
	UserJoined        = "user-joined"
)

func handleRegistration(connection *models.UserConnection) {
	//https://stackoverflow.com/questions/42605337/cannot-assign-to-struct-field-in-a-map
	entry, keyExists := rooms[connection.RoomId]

	if !keyExists {
		keyExistErr := connection.CloseConnectionWithMessage("there is no room with key: " + connection.RoomId.String())
		if keyExistErr != nil {
			connection.Logger.Printf("there was an error closing out a connection: %#v\n", keyExistErr)
		}
		return
	}

	if !entry.Joinable {
		joinErr := connection.CloseConnectionWithMessage("A game is already in progress in this room.")
		if joinErr != nil {
			connection.Logger.Printf("there was an error closing out a connection: %#v\n", joinErr)
		}
		return
	}

	if entry.OnBanList(connection) {
		onBanErr := connection.CloseConnectionWithMessage("User unable to join game")
		if onBanErr != nil {
			connection.Logger.Printf("there was an error closing out a connection: %#v\n", onBanErr)
		}
	}

	// This is the same kind of data being stored twice.
	// Need to work to combine.
	{
		entry.Connections[connection.Connection] = models.Client{
			ID:       connection.UserId,
			Username: connection.Username,
		}

		entry.ConnectedUsers[connection.UserId] = connection.Username
	}

	rooms[connection.RoomId] = entry

	connectMessage := models.InitialConnection{
		MessageType: InitialConnection,
		Action:      ActionRegistered,
		Admin:       connection.UserId == entry.AdminId,
	}
	if err := connection.Connection.WriteJSON(connectMessage); err != nil {
		Configs.Logger.Println(err)
	}

	// broadcast usernames so frontend can show connected users in waiting room
	joinedMessage := models.UserConnectedMessage{
		MessageType: UserJoined,
		Contents:    entry.GetConnectedList(),
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
	Text       = "TEXT"
	Submit     = "SUBMIT"
	Load       = "LOAD"
	GetResults = "GETRESULTS"
	BanUser    = "BANUSER"
	GetBanList = "GETBANLIST"
)

// We need to handle the state of the rooms and the q and a's
func handleBroadcast(message models.UserResponse) {
	message.Logger.Println("broadcast received:", message.String())
	room := rooms[message.RoomId]

	switch strings.ToUpper(message.UserMessage.Action) {
	case Load:
		room.Joinable = false
		// since room is a copy, we have to assign it back to the map
		rooms[message.RoomId] = room

		if message.UserId == room.AdminId {
			err := handleAdminLoad(message.RoomId, message.UserMessage.DeckId, message.Conn)
			if err != nil {
				Configs.Logger.Printf("there was an error: %#v\n", err)
				_ = message.Conn.WriteJSON(models.NewErrorResponse(err.Error()))
			}
		} else {
			Configs.Logger.Printf("%d is trying to access admin\n", message.UserId)
			_ = message.Conn.WriteJSON(models.NewErrorResponse("non admin user"))
		}
	case BanUser:
		handleAdminBanUser(&room, &message)
		rooms[message.RoomId] = room

	case GetBanList:
		if room.AdminId == message.UserId {
			if banErr := message.Conn.WriteJSON(room.BannedList); banErr != nil {
				Configs.Logger.Printf("there was an error: %#v", banErr)
				_ = message.Conn.WriteJSON(models.NewErrorResponse(banErr.Error()))
			}
		} else {
			Configs.Logger.Printf("%d is trying to access admin\n", message.UserId)
			_ = message.Conn.WriteJSON(models.NewErrorResponse("non admin user"))
		}

	case Text:
		handleTextMessage(message)
	case Submit:
		if err := handleAnswerSubmissions(message); err != nil {
			Configs.Logger.Printf("there was an error: %#v", err)
			_ = message.Conn.WriteJSON(models.NewErrorResponse(err.Error()))
		}
	case GetResults:
		if err := returnAllResults(message); err != nil {
			Configs.Logger.Printf("there was an error: %#v", err)
			_ = message.Conn.WriteJSON(models.NewErrorResponse(err.Error()))
		}
	}

}
func handleAdminBanUser(room *models.Room, message *models.UserResponse) {
	if room.AdminId == message.UserId {
		if bannedUserId, err := strconv.ParseUint(message.UserMessage.Message, 10, 64); err == nil {
			if errMap, count := banUser(uint(bannedUserId), room); count > 0 {
				for _, banErr := range errMap {
					Configs.Logger.Println(banErr.Error())
					Configs.Logger.Printf("unable to ban user %d -- %#v \n ", bannedUserId, banErr)
				}
			}

			if banErr := message.Conn.WriteJSON(room.BannedList); banErr != nil {
				Configs.Logger.Printf("there was an error: %#v", banErr)
				_ = message.Conn.WriteJSON(models.NewErrorResponse(banErr.Error()))
			}
		} else {
			Configs.Logger.Printf("%d sent an invalid userID to be banned\n", message.UserId)
			_ = message.Conn.WriteJSON(models.NewErrorResponse("unable to convert banned ID"))
		}
	} else {
		Configs.Logger.Printf("%d is trying to access admin\n", message.UserId)
		_ = message.Conn.WriteJSON(models.NewErrorResponse("non admin user"))
	}
}
func banUser(userToBanId uint, room *models.Room) (map[*websocket.Conn]error, int) {
	connectedErrors := make(map[*websocket.Conn]error, 0)
	errCount := 0

	room.BannedList = append(room.BannedList, models.BannedPlayer{ID: userToBanId})

	for conn, client := range room.Connections {
		if client.ID == userToBanId {
			_ = conn.WriteMessage(websocket.CloseMessage, []byte{})
			err := conn.Close()
			if err != nil {
				connectedErrors[conn] = err
				errCount++
			}
		}
	}

	return connectedErrors, errCount
}

func returnAllResults(message models.UserResponse) error {
	roomID := message.RoomId
	room := rooms[roomID]

	userResults := make(models.UserResults, 0)
	for userId, userScore := range room.Results {
		userResults[room.ConnectedUsers[userId]] = userScore
	}

	return message.Conn.WriteJSON(models.Results{
		RoomId:      message.RoomId,
		UserResults: userResults,
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

	//for conn := range room.Connections {
	//	_ = conn.WriteJSON(models.NewResult(message.UserId, room.ConnectedUsers[message.UserId], result, room.TotalQuestions))
	//}

	for conn, client := range room.Connections {
		_ = conn.WriteJSON(models.NewResult(message.UserId, client.Username, result, room.TotalQuestions))
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
	newRoom, keyExists := rooms[roomSetup.NewRoomID]

	if keyExists {
		roomSetup.Logger.Println("There was an existing room with key: " + roomSetup.NewRoomID.String())
		return
	}

	newRoom.AdminId = roomSetup.AdminId
	newRoom.Results = make(map[uint]uint, 0)
	newRoom.ConnectedUsers = make(map[uint]string, 0)
	newRoom.Connections = make(map[*websocket.Conn]models.Client)
	newRoom.Joinable = true
	newRoom.BannedList = roomSetup.BannedPlayers

	rooms[roomSetup.NewRoomID] = newRoom
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
