package web

import (
	"errors"
	"flag"
	"log"
	"os"
	"strings"
	"teamC/Global"
)

func LoadConfiguration(cfg *Global.Configuration) error {
	cfg.Logger = log.Default()

	var production string
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

	flag.Parse()

	// port required to be prefixed with colon
	if !strings.HasPrefix(cfg.Port, ":") {
		cfg.Port = ":" + cfg.Port
	}

	if strings.ToLower(production) == Global.PRODUCTION {
		cfg.Production = true
	} else {
		cfg.Production = false
	}

	if cfg.JwtSecret == "" {
		return errors.New("application cannot have a blank jwt secret")
	}

	if cfg.Production {
		if cfg.GoogleAuthKey == Global.NO_VALID_PRODUCTION_DEFAULT ||
			cfg.GoogleSecretKey == Global.NO_VALID_PRODUCTION_DEFAULT {
			return errors.New("application is missing production configuration data")
		}
	}

	return nil
}
