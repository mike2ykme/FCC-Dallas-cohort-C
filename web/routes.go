package web

import (
	"fmt"
    "io"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
    "net/http"
    "net/url"
    "encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"strconv"
	"teamC/Global"
	"teamC/models"
	"time"
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

func ProductionLoginHandler(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
        authCode := string(c.Body())

        tokenRequestResp, err := http.PostForm("https://oauth2.googleapis.com/token", url.Values{
            "client_id": {"849784468632-n9upp7q0umm82uecp5h3pfdervht7sjg.apps.googleusercontent.com"},
            "client_secret": {cfg.GoogleSecretKey},
            "code": {authCode},
            "grant_type": {"authorization_code"},
            "redirect_uri": {"http://127.0.0.1:3000/oauth-redirect"},
        })
        if err != nil {
            panic(err)  // not sure what fiber error would actually be best
        }
        defer tokenRequestResp.Body.Close()
        var tokenRespData map[string]interface{}
        var tokenRespBody []byte
        tokenRespBody, _ = io.ReadAll(tokenRequestResp.Body)  // not sure how I'd handle this error
        _ = json.Unmarshal(tokenRespBody, &tokenRespData)  // or this one
        accessToken := tokenRespData["access_token"].(string)

        // get endpoints from here: https://accounts.google.com/.well-known/openid-configuration
        // they're always changing them
        client := http.Client{}
        emailReq, _ := http.NewRequest("GET", "https://openidconnect.googleapis.com/v1/userinfo", nil)
        emailReq.Header.Set("Authorization", "Bearer " + accessToken)
        emailResp, emailRespErr := client.Do(emailReq)
        if emailRespErr != nil {
            panic(err)
        }
        defer emailResp.Body.Close()
        emailRespBody, ioErr := io.ReadAll(emailResp.Body)
        if ioErr != nil {
            panic(err)
        }
        var emailRespData map[string]interface{}
        _ = json.Unmarshal(emailRespBody, &emailRespData)
        fmt.Println(emailRespData)
        // contains email, email_verified, family_name, given_name, name, picture, sub, and locale.
        // we probably care about email and maybe sub. Conveniently, this means I don't have
        // to actually parse the jwt token to get the sub id, yay!

        return c.JSON(fiber.Map{"token": "placeholder"})

	}
}

func SimulatedLoginHandler(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create the Claims
		claims := jwt.MapClaims{
			"name":  "John Doe",
			"admin": true,
			"exp":   time.Now().Add(time.Hour * 72).Unix(),
		}

		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(cfg.JwtSecret))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(fiber.Map{"token": t})
	}
}
