package web

import (
	"flag"
	"strings"
	"teamC/models"
)

func LoadConfiguration(cfg *models.Configuration) error {
	flag.StringVar(&cfg.Port, "port", "9000", "Set Port to start server on, defaults to 8080")
	flag.BoolVar(&cfg.Production, "production", false, "Should we load testing pages and configuration")
	flag.StringVar(&cfg.JwtSecret, "jwtSecret", "SECRET", "value for the signing keys")
	flag.StringVar(&cfg.GoogleSecretKey, "googleSecretKey", "", "value used for the Google Secret Key")
	flag.StringVar(&cfg.GoogleAuthKey, "googleAuthKey", "", "value used for the Google Auth Key")

	//if config.GoogleOAuthSecret = os.Getenv(app.GOOGLE_SECRET_KEY_OS); config.GoogleOAuthSecret == "" {
	//	flag.StringVar(&config.GoogleOAuthSecret, app.GOOGLE_SECRET_KEY_FLAG, app.GOOGLE_SECRET_DEFAULT, "Used to set the Google OAuth Secret")
	//}

	/*

		GOOGLE_AUTH_KEY_OS   string = "GOOGLE_AUTH_KEY"
		GOOGLE_AUTH_KEY_FLAG string = "googleAuthKey"
		GOOGLE_AUTH_DEFAULT  string = ""

		GOOGLE_SECRET_KEY_OS   string = "GOOGLE_SECRET_KEY"
		GOOGLE_SECRET_KEY_FLAG string = "googleSecretKey"
		GOOGLE_SECRET_DEFAULT  string = ""
	*/
	flag.Parse()

	// port required to be prefixed with colon
	if !strings.HasPrefix(cfg.Port, ":") {
		cfg.Port = ":" + cfg.Port
	}

	// Leaving option open to return errors if we validate config in the future
	return nil
}
