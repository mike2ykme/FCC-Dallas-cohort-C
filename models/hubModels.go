package models

import (
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

type RoomCreation struct {
	AdminId uint
	NewRoom uuid.UUID
}
type UserResponse struct {
	Action    string          `json:"action"`
	Message   string          `json:"message"`
	Conn      *websocket.Conn `json:"-"`
	ChannelId uuid.UUID       `json:"-"`
}
type Room struct {
	ChannelId   uuid.UUID
	Connections map[*websocket.Conn]Client
	Admin       *websocket.Conn
	Password    string
}

type UserConnection struct {
	Connection *websocket.Conn
	ChannelId  uuid.UUID
}

type InitialConnection struct {
	Action   string         `json:"action"`
	Admin    bool           `json:"admin"`
	Question string         `json:"question"`
	Answers  []AnswerChoice `json:"answers"`
}
type AnswerChoice struct{}
type Client struct{} // Add more data to this type if needed
