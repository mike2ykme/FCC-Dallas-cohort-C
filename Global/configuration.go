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
	JWTExpiration   int
	RedirectUri     string
	DatabaseURL     string
	AutoMigrate     bool
	MaxWSErrors     int
	LoadTestDecks   bool
}

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
	// TESTING
	//

	OS_LOAD_TEST_DECKS      = "LOAD_TEST_DECKS"
	FLAG_LOAD_TEST_DECKS    = "loadTestDecks"
	DEFAULT_LOAD_TEST_DECKS = false
	LOAD_TEST_DECKS_USAGE   = "if environment is nonproduction and this is true then it will load 2 test decks"
	//
	// Database
	//

	OS_DATABASE_URL    = "DATABASE_URL"
	FLAG_DATABASE_URL  = "databaseUrl"
	DATABASE_URL_USAGE = "specifies the database connection string"

	OS_AUTO_MIGRATE      = "AUTO_MIGRATE"
	FLAG_AUTO_MIGRATE    = "autoMigrate"
	AUTO_MIGRATE_DEFAULT = true
	AUTO_MIGRATE_USAGE   = "specifies if we should automatically migrate the DBs"

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

	OS_JWT_EXPIRY                       = "JWT_EXPIRY"
	FLAG_JWT_EXPIRY                     = "jwtExpiry"
	JWT_EXPIRY_USAGE                    = "value used for the token expiry in hours"
	JWT_DEFAULT_OF_HOURS_IN_WEEK_STRING = "168"
	JWT_DEFAULT_OF_HOURS_IN_WEEK_INT    = 168

	OS_REDIRECT_URI      = "REDIRECT_URI"
	FLAG_REDIRECT_URI    = "redirectURI"
	REDIRECT_URI_USAGE   = "value used to redirect the user"
	REDIRECT_URI_DEFAULT = "http://127.0.0.1:3000/oauth-redirect"

	OS_MAX_WS_ERRORS      = "MAX_WS_ERRORS"
	FLAG_MAX_WS_ERRORS    = "maxWsErrors"
	MAX_WS_ERRORS_USAGE   = "value used to determine how many bad messages we're willing to accept from user"
	MAX_WS_ERRORS_DEFAULT = "3" // we will convert to int in setup
)
