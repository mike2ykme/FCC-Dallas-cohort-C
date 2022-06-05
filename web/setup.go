package web

import (
	"errors"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"teamC/Global"
)

func LoadConfiguration(cfg *Global.Configuration) error {
	cfg.Logger = log.Default()

	var production string
	var jwtExpirationTemp string

	if cfg.Port = os.Getenv(Global.OS_PORT); cfg.Port == Global.EMPTY_STRING {
		flag.StringVar(&cfg.Port, Global.FLAG_PORT, Global.DEFAULT_PORT, Global.PORT_USAGE)
	}

	if cfg.Production = strings.ToLower(os.Getenv(Global.OS_ENVIRONMENT)) == Global.PRODUCTION; cfg.Production == false {
		flag.StringVar(&production, Global.FLAG_ENVIRONMENT, Global.DEFAULT_ENVIRONMENT, Global.ENVIRONMENT_USAGE)
	}

	if cfg.JwtSecret = os.Getenv(Global.OS_JWT); cfg.JwtSecret == Global.EMPTY_STRING {
		flag.StringVar(&cfg.JwtSecret, Global.FLAG_JWT, Global.NO_VALID_DEFAULT, "value for the signing keys")
	}

	if cfg.GoogleSecretKey = os.Getenv(Global.OS_GOOGLE_SECRET_KEY); cfg.GoogleSecretKey == Global.EMPTY_STRING {
		flag.StringVar(&cfg.GoogleSecretKey, Global.FLAG_GOOGLE_SECRET_KEY, Global.NO_VALID_PRODUCTION_DEFAULT, Global.GOOGLE_SECRET_KEY_USAGE)
	}

	if cfg.GoogleAuthKey = os.Getenv(Global.OS_GOOGLE_AUTH_KEY); cfg.GoogleAuthKey == Global.EMPTY_STRING {
		flag.StringVar(&cfg.GoogleAuthKey, Global.FLAG_GOOGLE_AUTH_KEY, Global.NO_VALID_PRODUCTION_DEFAULT, Global.GOOGLE_AUTH_KEY_USAGE)
	}

	if jwtExpirationTemp = os.Getenv(Global.OS_JWT_EXPIRY); jwtExpirationTemp == Global.EMPTY_STRING {
		flag.StringVar(&jwtExpirationTemp, Global.FLAG_JWT_EXPIRY, Global.JWT_DEFAULT_OF_HOURS_IN_WEEK_STRING, Global.JWT_EXPIRY_USAGE)
	}

	if cfg.RedirectUri = os.Getenv(Global.OS_REDIRECT_URI); cfg.RedirectUri == Global.EMPTY_STRING {
		flag.StringVar(&cfg.RedirectUri, Global.FLAG_REDIRECT_URI, Global.NO_VALID_PRODUCTION_DEFAULT, Global.REDIRECT_URI_USAGE)
	}

	if cfg.DatabaseURL = os.Getenv(Global.OS_DATABASE_URL); cfg.DatabaseURL == Global.EMPTY_STRING {
		flag.StringVar(&cfg.DatabaseURL, Global.FLAG_DATABASE_URL, Global.NO_VALID_PRODUCTION_DEFAULT, Global.DATABASE_URL_USAGE)
	}

	if cfg.AutoMigrate = strings.ToLower(os.Getenv(Global.OS_AUTO_MIGRATE)) == "true"; cfg.AutoMigrate == false {
		flag.BoolVar(&cfg.AutoMigrate, Global.FLAG_AUTO_MIGRATE, Global.AUTO_MIGRATE_DEFAULT, Global.AUTO_MIGRATE_USAGE)
	}

	flag.Parse()

	// port required to be prefixed with colon
	if !strings.HasPrefix(cfg.Port, ":") {
		cfg.Port = ":" + cfg.Port
	}

	if cfg.Production == false {
		if strings.ToLower(production) == Global.PRODUCTION {
			cfg.Production = true
		} else {
			cfg.Production = false
		}
	}

	if val, err := strconv.Atoi(jwtExpirationTemp); err == nil && val > 0 {
		cfg.JWTExpiration = val
	} else {
		cfg.JWTExpiration = Global.JWT_DEFAULT_OF_HOURS_IN_WEEK_INT
		cfg.Logger.Printf("there was an error converting jwtExpiration: %#V \t or val was less than 1 %d \n", err.Error(), val)
	}

	if cfg.JwtSecret == "" {
		return errors.New("application cannot have a blank jwt secret")
	}

	if cfg.Production {
		if cfg.GoogleAuthKey == Global.NO_VALID_PRODUCTION_DEFAULT ||
			cfg.GoogleSecretKey == Global.NO_VALID_PRODUCTION_DEFAULT ||
			cfg.RedirectUri == Global.EMPTY_STRING ||
			cfg.DatabaseURL == Global.EMPTY_STRING {
			return errors.New("application is missing production configuration data")
		}
	} else {
		if cfg.RedirectUri == Global.EMPTY_STRING {
			cfg.RedirectUri = Global.REDIRECT_URI_DEFAULT
		}
	}

	return nil
}
