package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"log"
	"net/http"
	"net/url"
	"teamC/Global"
	"time"
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
		logger := cfg.Logger
		authCode := string(c.Body())
		//var accessToken string

		accessToken, err := getAccessToken(cfg, authCode)
		if err != nil {
			logger.Printf("There was an error getting the access token %#V\n", err)
			c.SendStatus(fiber.StatusInternalServerError)
		}
		oauthUser, err := getOauthUser(accessToken)

		_ = oauthUser

		return c.Status(fiber.StatusOK).SendString("NOT YET IMPLEMENTED")

	}
}
func getOauthUser(accessToken string) (map[string]interface{}, error) {
	// If they're changing them, and we can get them from the below URL then we need to request and parse it when we request it
	// get endpoints from here: https://accounts.google.com/.well-known/openid-configuration
	// they're always changing them
	client := http.Client{}
	emailReq, err := http.NewRequest("GET", "https://openidconnect.googleapis.com/v1/userinfo", nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("there was an error doing new GET request: %s\n", err))
	}
	emailReq.Header.Set("Authorization", "Bearer "+accessToken)
	emailResp, emailRespErr := client.Do(emailReq)
	if emailRespErr != nil {
		return nil, errors.New(fmt.Sprintf("there was an error doing request: %s\n", err))
	}
	defer emailResp.Body.Close()

	emailRespBody, ioErr := io.ReadAll(emailResp.Body)

	if ioErr != nil {
		return nil, errors.New(fmt.Sprintf("there was an error reading request body: %s\n", err))
	}

	var emailRespData map[string]interface{}
	err = json.Unmarshal(emailRespBody, &emailRespData)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("there was an error unmarshalling data: %s\n", err))
	}

	return emailRespData, nil
	//fmt.Println(emailRespData)
	// contains email, email_verified, family_name, given_name, name, picture, sub, and locale.
	// we probably care about email and maybe sub. Conveniently, this means I don't have
	// to actually parse the jwt token to get the sub id, yay!

	//return c.JSON(fiber.Map{"token": "placeholder"})
}
func getAccessToken(cfg *Global.Configuration, authCode string) (string, error) {
	tokenRequestResp, err := http.PostForm(cfg.OauthPostURL, url.Values{
		"client_id":     {cfg.ClientId},
		"client_secret": {cfg.GoogleSecretKey},
		"code":          {authCode},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {cfg.RedirectURL},
	})

	if err != nil {
		//	Don't panic, this is something we expect to happen occasionally.
		return "", errors.New("there was a problem trying to get the response from Google")
	}
	reqBody := tokenRequestResp.Body
	defer reqBody.Close()

	var tokenRespData map[string]interface{}

	//tokenRespBody, err := io.ReadAll(reqBody) // not sure how I'd handle this error
	tokenRespBody, err := io.ReadAll(reqBody)
	if err != nil {
		return "", errors.New(("there was a problem reading the request body"))
	}

	err = json.Unmarshal(tokenRespBody, &tokenRespData) // or this one
	if err != nil {
		return "", errors.New("there was a problem unmarshalling the data into an struct")
	}

	if accessToken, ok := tokenRespData["access_token"].(string); ok {
		return accessToken, nil
	}

	return "", errors.New("There was an error trying to cast the access_token into a string")
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
