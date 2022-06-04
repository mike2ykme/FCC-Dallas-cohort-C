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
	"teamC/db"
	"teamC/models"
	"time"
)

func GetJwtMiddleware(cfg *Global.Configuration) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(cfg.JwtSecret),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			cfg.Logger.Println("jwt error handler called, returning 404 to user-- ", err)
			return c.SendStatus(fiber.StatusNotFound)
		},
	})
}

func ProductionLoginHandler(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authCode := string(c.Body())

		if accessToken, err := getAccessToken(authCode, cfg); err == nil {

			if googleResp, err := getGoogleResponse(accessToken); err == nil {

				if user, err := getUser(googleResp, cfg.UserRepo, cfg.Logger); err == nil {

					if signedJWT, err := mapUserToSignedJWT(&user, cfg); err == nil {
						return c.JSON(fiber.Map{"token": signedJWT})
					} else {
						cfg.Logger.Println(fmt.Sprintf("there was an error trying to get the user, %#v", err))
					}
				} else { // getUser
					cfg.Logger.Println(fmt.Sprintf("there was an error trying to get the user, %#v", err))
				}
			} else { // getGoogleResponse
				cfg.Logger.Println(fmt.Sprintf("there was an error, %#v", err))
			}
		} else { // getAccessToken
			cfg.Logger.Println(fmt.Sprintf("there was an error, %#v", err))
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}
}

func mapUserToSignedJWT(user *models.User, cfg *Global.Configuration) (string, error) {
	claims := jwt.MapClaims{
		"username":  user.Username,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"id":        user.ID,
		"admin":     false,
		"exp":       time.Now().Add(time.Hour * time.Duration(cfg.JWTExpiration)).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(cfg.JwtSecret))
	if err != nil {
		return "", err
		//return c.SendStatus(fiber.StatusInternalServerError)
	}
	return t, nil
}

func getValFromInterfaceMap(googleResponse map[string]interface{}, str string) string {
	if val, ok := googleResponse[str]; ok {
		if converted, ok := val.(string); ok {
			return converted
		}
	}
	return ""
}

const (
	SUB         = "sub"
	EMAIL       = "email"
	GIVEN_NAME  = "given_name"
	FAMILY_NAME = "family_name"
)

func getUser(googleResponse map[string]interface{}, userRepo db.UserRepository, logger *log.Logger) (models.User, error) {
	subId := getValFromInterfaceMap(googleResponse, SUB)
	if subId == "" {
		return models.User{}, errors.New("there was a problem casting the sub ID to string")
	}

	user := models.User{}
	err := userRepo.GetUserBySubId(&user, subId)
	if err != nil {
		return models.User{}, errors.New("there was an error getting the user by sub ID")
	}

	if user.ID == 0 {
		user = models.User{
			Username:  getValFromInterfaceMap(googleResponse, EMAIL),
			SubId:     subId,
			FirstName: getValFromInterfaceMap(googleResponse, GIVEN_NAME),
			LastName:  getValFromInterfaceMap(googleResponse, FAMILY_NAME),
		}
		userRepo.SaveUser(&user)
		logger.Println("We've saved a user")
	} else {
		logger.Println("We already had this user")
	}

	return user, nil
}

func getGoogleResponse(accessToken string) (map[string]interface{}, error) {
	// get endpoints from here: https://accounts.google.com/.well-known/openid-configuration
	// they're always changing them
	client := http.Client{}
	emailReq, err := http.NewRequest("GET", "https://openidconnect.googleapis.com/v1/userinfo", nil)
	if err != nil {
		return nil, err
	}
	emailReq.Header.Set("Authorization", "Bearer "+accessToken)
	emailResp, emailRespErr := client.Do(emailReq)

	if emailRespErr != nil {
		return nil, emailRespErr
	}
	defer emailResp.Body.Close()

	emailRespBody, ioErr := io.ReadAll(emailResp.Body)
	if ioErr != nil {
		return nil, ioErr
	}

	var emailRespData map[string]interface{}
	err = json.Unmarshal(emailRespBody, &emailRespData)
	//fmt.Println(emailRespData)
	return emailRespData, err
	//googleResponse = emailRespData

}

const (
	CLIENT_ID  = "849784468632-n9upp7q0umm82uecp5h3pfdervht7sjg.apps.googleusercontent.com"
	GRANT_TYPE = "authorization_code"
)

func getAccessToken(authCode string, cfg *Global.Configuration) (string, error) {
	tokenRequestResp, err := http.PostForm("https://oauth2.googleapis.com/token", url.Values{
		"client_id":     {CLIENT_ID},
		"client_secret": {cfg.GoogleSecretKey},
		"code":          {authCode},
		"grant_type":    {GRANT_TYPE},
		"redirect_uri":  {cfg.RedirectUri},
	})
	if err != nil {
		return "", err
	}
	defer tokenRequestResp.Body.Close()
	var tokenRespData map[string]interface{}
	var tokenRespBody []byte
	tokenRespBody, err = io.ReadAll(tokenRequestResp.Body) // not sure how I'd handle this error
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(tokenRespBody, &tokenRespData) // or this one
	if err != nil {
		return "", err
	}
	token, ok := tokenRespData["access_token"].(string)
	if !ok {
		return "", errors.New("there was an error casting the access_token to a string")
	}
	return token, nil

}

//const hoursInWeek = 168

func SimulatedLoginHandler(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create the Claims
		claims := jwt.MapClaims{
			"username":  "John Doe",
			"firstName": "John",
			"lastName":  "Doe",
			"id":        uint(1),
			"admin":     true,
			"exp":       time.Now().Add(time.Hour * time.Duration(cfg.JWTExpiration)).Unix(),
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
