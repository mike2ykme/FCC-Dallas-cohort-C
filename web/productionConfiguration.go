package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"io"
	"net/http"
	"teamC/db/rdbms"

	"teamC/Global"
	"time"
)

// get endpoints from here: https://accounts.google.com/.well-known/openid-configuration
// they're always changing them
func getUserInfoEndpoint() (string, error) {
	resp, err := http.Get("https://accounts.google.com/.well-known/openid-configuration")
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var gResp = make(map[string]interface{})
	err = json.Unmarshal(body, &gResp)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%#v", gResp["userinfo_endpoint"]), nil
}

func ProductionConfiguration(cfg *Global.Configuration) error {
	{
		repo, err := rdbms.NewRdbmsRepository(cfg.DatabaseURL, "postgres")
		if cfg.AutoMigrate {
			repo.AutoMigrate()
		}
		if err != nil {
			return err
		}
		cfg.UserRepo = repo
		cfg.DeckRepo = repo
		cfg.FlashcardRepo = repo
		cfg.AnswerRepo = repo
	}
	endpoint, err := getUserInfoEndpoint()
	if err != nil {
		return err
	}
	if endpoint == "" {
		return errors.New("invalid return from getUserInfoEndpoint, return was empty")
	}

	cfg.UserInfoEndpoint = endpoint

	app := cfg.WebApp
	app.Use(cors.New())
	app.Use(limiter.New(limiter.Config{
		Max:        cfg.LimiterConfig.Max,
		Expiration: cfg.LimiterConfig.ExpirationSeconds * time.Second,
	}))

	app.Post("/login", ProductionLoginHandler(cfg))
	return nil
}
