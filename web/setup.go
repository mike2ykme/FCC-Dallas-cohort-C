package web

import (
	"flag"
	"strings"
	"teamC/models"
)

func LoadConfiguration(cfg *models.Configuration) error {
	flag.StringVar(&cfg.Port, "port", "9000", "Set Port to start server on, defaults to 8080")
	flag.BoolVar(&cfg.Production, "production", false, "Should we load testing pages and configuration")

	flag.Parse()

	// port required to be prefixed with colon
	if !strings.HasPrefix(cfg.Port, ":") {
		cfg.Port = ":" + cfg.Port
	}

	// Leaving option open to return errors if we validate config in the future
	return nil
}
