package web

import (
	"errors"
	"flag"
	"os"
	"strings"
	"teamC/models"
)

func LoadConfiguration(cfg *models.Configuration) error {
	var production string
	if cfg.Port = os.Getenv(models.OS_PORT); cfg.Port == models.EMPTY_STRING {
		flag.StringVar(&cfg.Port, models.FLAG_PORT, models.DEFAULT_PORT, models.PORT_USAGE)
	}
	if cfg.Production = strings.ToLower(os.Getenv(models.OS_ENVIRONMENT)) == models.PRODUCTION; cfg.Production == false {
		flag.StringVar(&production, models.FLAG_ENVIRONMENT, models.DEFAULT_ENVIRONMENT, models.ENVIRONMENT_USAGE)
	}

	if cfg.JwtSecret = os.Getenv(models.OS_JWT); cfg.JwtSecret == models.EMPTY_STRING {
		flag.StringVar(&cfg.JwtSecret, models.FLAG_JWT, models.NO_VALID_DEFAULT, "value for the signing keys")
	}

	if cfg.GoogleSecretKey = os.Getenv(models.OS_GOOGLE_SECRET_KEY); cfg.GoogleSecretKey == models.EMPTY_STRING {
		flag.StringVar(&cfg.GoogleSecretKey, models.FLAG_GOOGLE_SECRET_KEY, models.NO_VALID_PRODUCTION_DEFAULT, models.GOOGLE_SECRET_KEY_USAGE)
	}

	if cfg.GoogleAuthKey = os.Getenv(models.OS_GOOGLE_AUTH_KEY); cfg.GoogleAuthKey == models.EMPTY_STRING {
		flag.StringVar(&cfg.GoogleAuthKey, models.FLAG_GOOGLE_AUTH_KEY, models.NO_VALID_PRODUCTION_DEFAULT, models.GOOGLE_AUTH_KEY_USAGE)
	}

	flag.Parse()

	// port required to be prefixed with colon
	if !strings.HasPrefix(cfg.Port, ":") {
		cfg.Port = ":" + cfg.Port
	}

	if strings.ToLower(production) == models.PRODUCTION {
		cfg.Production = true
	} else {
		cfg.Production = false
	}

	if cfg.JwtSecret == "" {
		return errors.New("application cannot have a blank jwt secret")
	}

	if cfg.Production {
		if cfg.GoogleAuthKey == models.NO_VALID_PRODUCTION_DEFAULT ||
			cfg.GoogleSecretKey == models.NO_VALID_PRODUCTION_DEFAULT {
			return errors.New("application is missing production configuration data")
		}
	}

	return nil
}
