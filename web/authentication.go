package web

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"teamC/Global"
	"time"
    "net/http"
    "io"
    "net/url"
    "encoding/json"
)

func GetJwtMiddleware(cfg *Global.Configuration) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(cfg.JwtSecret),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Println("jwt error handler called, returning 404 to user-- ", err)
			return c.SendStatus(fiber.StatusNotFound)
		},
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
        //fmt.Println(emailRespData)
        // contains email, email_verified, family_name, given_name, name, picture, sub, and locale.
        // we probably care about email and maybe sub. Conveniently, this means I don't have
        // to actually parse the jwt token to get the sub id, yay!

        return c.JSON(fiber.Map{"token": "placeholder"})

	}
}

const hoursInWeek = 168

func SimulatedLoginHandler(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create the Claims
		claims := jwt.MapClaims{
			"username":  "John Doe",
			"firstName": "John",
			"lastName":  "Doe",
			"id":        uint(1),
			"admin":     true,
			"exp":       time.Now().Add(time.Hour * hoursInWeek).Unix(),
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
