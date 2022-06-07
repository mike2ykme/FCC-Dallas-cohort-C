package web

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"teamC/Global"
	"teamC/models"
)

func WebsocketRoom(cfg *Global.Configuration) fiber.Handler {
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
		userId := c.Locals(USER_ID).(uint)
		userConn := &models.UserConnection{
			UserId:     userId,
			Connection: c,
			RoomId:     channelId,
			Logger:     cfg.Logger,
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
			Logger:    cfg.Logger,
		}
		errorCount := 0
		for {
			if err := c.ReadJSON(&response); err == nil {
				broadcast <- response

			} else {
				errorCount += 1
				cfg.Logger.Printf("there was an error trying to read a json response from %d", userId)

				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					cfg.Logger.Println("read error:", err)
					return // Calls the deferred function, i.e. closes the connection on error
				}

				//c.WriteMessage(websocket.TextMessage, []byte("invalid message"))
				c.WriteJSON(models.UserResponse{
					Action:  "ERROR",
					Message: "invalid message received",
				})

				if errorCount > cfg.MaxWSErrors {
					c.WriteJSON(models.UserResponse{
						Action:  "CLOSING",
						Message: "goodbye",
					})
					cfg.Logger.Println("User surpassed max errors, terminating")
					return
				}
			}
		}
	})
}
