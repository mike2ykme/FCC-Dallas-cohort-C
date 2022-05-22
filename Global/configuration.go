package Global

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"teamC/db"
	"time"
)

type Configuration struct {
	Port            string
	Production      bool
	WebApp          *fiber.App
	LimiterConfig   LimiterConfig
	JwtSecret       string
	GoogleSecretKey string
	GoogleAuthKey   string
	UserRepo        db.UserRepository
	DeckRepo        db.DeckRepository
	FlashcardRepo   db.FlashcardRepository
	AnswerRepo      db.AnswerRepository
	Logger          *log.Logger
}

//const logger = log.New(os.Stdout)

type LimiterConfig struct {
	Max               int
	ExpirationSeconds time.Duration
}

const (
	NO_VALID_DEFAULT            = ""
	NO_VALID_PRODUCTION_DEFAULT = ""
	EMPTY_STRING                = ""
	//
	// Webserver setup
	//

	// PORT
	FLAG_PORT    = "port"
	OS_PORT      = "PORT"
	DEFAULT_PORT = ":8080"
	PORT_USAGE   = "Set Port to start server on, defaults to " + DEFAULT_PORT

	// ENVIRONMENT
	PRODUCTION          = "production"
	NONPRODUCTION       = "nonproduction"
	FLAG_ENVIRONMENT    = "environment"
	OS_ENVIRONMENT      = "ENVIRONMENT"
	DEFAULT_ENVIRONMENT = NONPRODUCTION
	ENVIRONMENT_USAGE   = "specifies whether this is a production or nonproduction environment"

	//
	// Authentication
	//

	// JWT
	OS_JWT   = "JWT_SECRET"
	FLAG_JWT = "jwtSecret"

	// Google Secret Key
	FLAG_GOOGLE_SECRET_KEY  = "googleSecretKey"
	OS_GOOGLE_SECRET_KEY    = "GOOGLE_SECRET_KEY"
	GOOGLE_SECRET_KEY_USAGE = "value used for the Google Secret Key"

	// Google Auth Key
	OS_GOOGLE_AUTH_KEY    = "GOOGLE_AUTH_KEY"
	FLAG_GOOGLE_AUTH_KEY  = "googleAuthKey"
	GOOGLE_AUTH_KEY_USAGE = "value used for the Google Auth Key"
)
