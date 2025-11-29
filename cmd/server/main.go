package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/ciliverse/cilikube/internal/app"
)

// just do it ! go!go!go!
func main() {
	configPath := getConfigPath()

	application, err := app.New(configPath)
	if err != nil {
		slog.Error("failed to initialize app", "error", err)
		os.Exit(1)
	}

	slog.Info("starting application", "config", configPath)
	application.Run()
}

// better way to do
func getConfigPath() string {
	config := flag.String("config", "", "config file path")
	flag.Parse()

	if *config != "" {
		return *config
	}

	if env := os.Getenv("CILIKUBE_CONFIG_PATH"); env != "" {
		return env
	}

	return "configs/config.yaml"
}
