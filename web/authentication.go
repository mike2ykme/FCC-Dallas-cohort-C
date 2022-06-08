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
	"strconv"
	"strings"
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

				// TODO this needs to be further broken down because it's doing 3 things
				// 1. Parse
				// 2. Querying DB for user
				// 3. Saving user if new user
				// finally returning fully realized user
				if user, err := parseResponseAndSaveUser(googleResp, cfg.UserRepo, cfg.Logger); err == nil {

					if signedJWT, err := mapUserToSignedJWT(&user, cfg); err == nil {
						return c.JSON(fiber.Map{"token": signedJWT})
					} else {
						cfg.Logger.Println(fmt.Sprintf("there was an error trying to get the user, %#v", err))
					}
				} else { // parseResponseAndSaveUser
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

func parseResponseAndSaveUser(googleResponse map[string]interface{}, userRepo db.UserRepository, logger *log.Logger) (models.User, error) {
	subId := getValFromInterfaceMap(googleResponse, SUB)
	if subId == "" {
		return models.User{}, errors.New("there was a problem casting the sub ID to string")
	}

	user := models.User{}
	err := userRepo.GetUserBySubId(&user, subId)
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return models.User{}, errors.New(fmt.Sprintf("there was an error getting the user by sub ID, err is: %#v", err))
	}

	if user.ID == 0 {
		user = models.User{
			Username:  getValFromInterfaceMap(googleResponse, EMAIL),
			SubId:     subId,
			FirstName: getValFromInterfaceMap(googleResponse, GIVEN_NAME),
			LastName:  getValFromInterfaceMap(googleResponse, FAMILY_NAME),
		}
		_, saveErr := userRepo.SaveUser(&user)
		if saveErr != nil {
			return models.User{}, saveErr
		}
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

func SimulatedLoginHandler(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id", "1")

		newId, err := strconv.ParseUint(id, 10, 64)

		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		var user models.User

		if err = cfg.UserRepo.GetUserById(&user, uint(newId)); err != nil && strings.Contains(err.Error(), "unable to find") {
			user.Username = fmt.Sprintf("John_Doe%d", newId)
			user.SubId = fmt.Sprintf("SUBID-%d", newId)
			user.FirstName = fmt.Sprintf("John%d", newId)
			user.LastName = fmt.Sprintf("Doe%d", newId)
			cfg.UserRepo.SaveUser(&user)
		}

		t, err := mapUserToSignedJWT(&user, cfg)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(fiber.Map{"token": t})
	}
}
