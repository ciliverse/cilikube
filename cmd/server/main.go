package main

import (
	"log/slog"
	"os"

	"github.com/ciliverse/cilikube/internal/app"
)

// just do it ! go!go!go!
func main() {
	configPath := app.GetConfigPath()
	application, err := app.New(configPath)
	if err != nil {
		slog.Error("failed to initialize app", "error", err)
		os.Exit(1)
	}

	slog.Info("starting application", "config", configPath)
	application.Run()
}
