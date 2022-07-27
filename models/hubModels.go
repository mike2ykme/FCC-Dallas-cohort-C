package models

import (
	"fmt"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"log"
)

type BannedPlayer struct {
	ID       uint   `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}

type RoomCreation struct {
	AdminId       uint
	NewRoomID     uuid.UUID
	Logger        *log.Logger
	BannedPlayers []BannedPlayer
}
type UserResponse struct {
	UserMessage
	Conn      *websocket.Conn `json:"-"`
	ChannelId uuid.UUID       `json:"-"`
	Logger    *log.Logger     `json:"-"`
	UserId    uint
	RoomId    uuid.UUID
}

func (u *UserResponse) String() string {
	return fmt.Sprintf("Message: %s \t Action: %s \t RoomID: %s", u.UserMessage.Message, u.UserMessage.Action, u.RoomId.String())
}

type UserMessage struct {
	Action  string `json:"action,omitempty"`
	Message string `json:"message,omitempty"`
	DeckId  int    `json:"deckId,omitempty"`
}

type ServerResponse struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

type Room struct {
	ChannelId      uuid.UUID
	Connections    map[*websocket.Conn]Client
	Admin          *websocket.Conn
	Password       string `json:"-"`
	AdminId        uint
	Deck           Deck
	Results        map[uint]uint
	ConnectedUsers map[uint]string
	TotalQuestions int
	Joinable       bool
	BannedList     []BannedPlayer `json:"-"`
}

func (r *Room) WriteJsonToAllConnections(json interface{}) (map[*websocket.Conn]error, int) {
	connectedErrors := make(map[*websocket.Conn]error, 0)
	errCount := 0
	for conn, _ := range r.Connections {
		err := conn.WriteJSON(json)
		if err != nil {
			connectedErrors[conn] = err
			errCount++
		}
	}
	return connectedErrors, errCount
}

type ConnectedUser struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type UserConnectedMessage struct {
	MessageType string          `json:"message_type"`
	Contents    []ConnectedUser `json:"contents"`
}

func (r *Room) GetConnectedList() []ConnectedUser {
	list := make([]ConnectedUser, len(r.Connections))
	idx := 0
	for _, client := range r.Connections {
		list[idx] = ConnectedUser{
			ID:       client.ID,
			Username: client.Username,
		}
		idx++
	}

	//list := make([]ConnectedUser, len(r.ConnectedUsers))
	//for userId, userName := range r.ConnectedUsers {
	//	list[idx] = ConnectedUser{
	//		ID:       userId,
	//		Username: userName,
	//	}
	//	idx++
	//}

	return list
}

func (r *Room) OnBanList(connection *UserConnection) bool {
	for _, bannedPlayer := range r.BannedList {
		if connection.UserId == bannedPlayer.ID || connection.Username == bannedPlayer.Username {
			return true
		}
	}
	return false
}

type UserConnection struct {
	Connection *websocket.Conn
	RoomId     uuid.UUID
	UserId     uint
	Logger     *log.Logger `json:"-"`
	Username   string
}

func (c *UserConnection) CloseConnectionWithMessage(message string) error {
	c.Logger.Println(message)
	_ = c.Connection.WriteMessage(websocket.CloseMessage, []byte{})
	err := c.Connection.Close()
	return err
}

type InitialConnection struct {
	MessageType string         `json:"message_type"`
	Action      string         `json:"action"`
	Admin       bool           `json:"admin"`
	Question    string         `json:"question"`
	Answers     []AnswerChoice `json:"answers"`
}
type AnswerChoice struct{}
type Client struct {
	ID       uint
	Username string
} // Add more data to this type if needed

type LoadDeck struct {
	MessageType string      `json:"message_type"`
	Deck        []FlashCard `json:"deck"`
	Count       int         `json:"count"`
}

func NewLoadDeck(d []FlashCard) LoadDeck {
	return LoadDeck{
		MessageType: "questions",
		Deck:        d,
		Count:       len(d),
	}
}

type Result struct {
	Task           string `json:"task"`
	UserId         uint   `json:"userId"`
	Username       string `json:"username"`
	Score          int    `json:"score"`
	TotalQuestions int    `json:"totalQuestions"`
}

func NewResult(uId uint, name string, score, count int) Result {
	return Result{
		Task:           "RESULT",
		UserId:         uId,
		Username:       name,
		Score:          score,
		TotalQuestions: count,
	}
}

type ErrorResponse struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

func NewErrorResponse(s string) ErrorResponse {
	return ErrorResponse{
		Action:  "ERROR",
		Message: s,
	}
}
func DefaultErrorResponse() ErrorResponse {
	return NewErrorResponse("invalid message received")
}

type UserResults map[string]uint

type Results struct {
	RoomId uuid.UUID `json:"roomId"`
	UserResults
}
