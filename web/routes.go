package web

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"log"
	"teamC/models"
)

func WebsocketRoom() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		// When the function returns, unregister the client and close the connection
		fmt.Println("output", c.Params("id", "?"))

		//channelId, err := strconv.ParseUint(c.Params("id", "0"), 10, 64)
		channelId, err := uuid.Parse(c.Params("id", ""))
		if err != nil {
			fmt.Printf("the error is %#v", err)
			fmt.Println(uuid.New())
			return
		}

		userConn := &models.UserConnection{
			Connection: c,
			ChannelId:  channelId,
		}
		// when we exit the function we'll remove these from continuing the broadcasted to
		defer func() {
			unregister <- userConn
			c.Close()
		}()

		register <- userConn

		response := models.UserResponse{
			Conn:      c,
			ChannelId: channelId,
		}
		for {
			if err := c.ReadJSON(&response); err == nil {
				broadcast <- response

			} else {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("read error:", err)
				}
				return // Calls the deferred function, i.e. closes the connection on error
			}
		}
	})
}
