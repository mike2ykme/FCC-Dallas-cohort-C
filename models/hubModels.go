package models

import (
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"log"
)

type RoomCreation struct {
	AdminId   uint
	NewRoomID uuid.UUID
	Logger    *log.Logger
}
type UserResponse struct {
	UserMessage
	Conn      *websocket.Conn `json:"-"`
	ChannelId uuid.UUID       `json:"-"`
	Logger    *log.Logger     `json:"-"`
	UserId    uint
	RoomId    uuid.UUID
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
}

func (r *Room) WriteJsonToAllConnections(i interface{}) (map[*websocket.Conn]error, int) {
	connectedErrors := make(map[*websocket.Conn]error, 0)
	errCount := 0
	for conn, _ := range r.Connections {
		err := conn.WriteJSON(i)
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
    MessageType string `json:"message_type"`
    Contents []ConnectedUser `json:"contents"`
}

func (r *Room) GetConnectedList() []ConnectedUser {
	list := make([]ConnectedUser, len(r.ConnectedUsers))
	idx := 0
	for userId, userName := range r.ConnectedUsers {
		list[idx] = ConnectedUser{
			ID:       userId,
			Username: userName,
		}
		idx++
	}

	return list
}

type UserConnection struct {
	Connection *websocket.Conn
	RoomId     uuid.UUID
	UserId     uint
	Logger     *log.Logger `json:"-"`
	Username   string
}

// I don't like having a separate type for the contents but don't know any
// other way to do this --Joe
type InitialConnectionContents struct {
    Action   string         `json:"action"`
    Admin    bool           `json:"admin"`
    Question string         `json:"question"`
    Answers  []AnswerChoice `json:"answers"`
}
type InitialConnection struct {
    MessageType string `json:"message_type"`
    Contents InitialConnectionContents `json:"contents"`
}
type AnswerChoice struct{}
type Client struct{} // Add more data to this type if needed

type LoadDeck struct {
	Task  string
	Deck  []FlashCard
	Count int
}

func NewLoadDeck(d []FlashCard) LoadDeck {
	return LoadDeck{
		Task:  "QUESTIONS",
		Deck:  d,
		Count: len(d),
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
