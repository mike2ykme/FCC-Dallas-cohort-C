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
	Action    string          `json:"action"`
	Message   string          `json:"message"`
	Conn      *websocket.Conn `json:"-"`
	ChannelId uuid.UUID       `json:"-"`
	Logger    *log.Logger     `json:"-"`
}
type Room struct {
	ChannelId   uuid.UUID
	Connections map[*websocket.Conn]Client
	Admin       *websocket.Conn
	Password    string `json:"-"`
	AdminId     uint
}

type UserConnection struct {
	Connection *websocket.Conn
	RoomId     uuid.UUID
	UserId     uint
	Logger     *log.Logger `json:"-"`
}

type InitialConnection struct {
	Action   string         `json:"action"`
	Admin    bool           `json:"admin"`
	Question string         `json:"question"`
	Answers  []AnswerChoice `json:"answers"`
}
type AnswerChoice struct{}
type Client struct{} // Add more data to this type if needed
