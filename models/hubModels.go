package models

import "github.com/gofiber/websocket/v2"

type UserResponse struct {
	Action    string          `json:"action"`
	Message   string          `json:"message"`
	Conn      *websocket.Conn `json:"-"`
	ChannelId uint64          `json:"-"`
}
type Room struct {
	ChannelId   uint64
	Connections map[*websocket.Conn]Client
	Admin       *websocket.Conn
	Password    string
}

type UserConnection struct {
	Connection *websocket.Conn
	ChannelId  uint64
}

type InitialConnection struct {
	Action   string         `json:"action"`
	Admin    bool           `json:"admin"`
	Question string         `json:"question"`
	Answers  []AnswerChoice `json:"answers"`
}
type AnswerChoice struct{}
type Client struct{} // Add more data to this type if needed
