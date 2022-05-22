package web

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"log"
	"strconv"
	"teamC/models"
)

func WebsocketRoom() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		// When the function returns, unregister the client and close the connection
		fmt.Println("output", c.Params("id", "?"))

		channelId, err := strconv.ParseUint(c.Params("id", "0"), 10, 64)
		if err != nil {
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

			err := c.ReadJSON(&response)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("read error:", err)
				}
				return // Calls the deferred function, i.e. closes the connection on error
			}
			broadcast <- response
		}
	})
}

