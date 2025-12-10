package main

import (
	"globalping/internal/config"
	"log/slog"
	"os"
)

func init() {
	// logging
	logLevel := slog.LevelInfo
	handlerOpts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, handlerOpts))
	slog.SetDefault(logger)
}
func main() {
	config, err := config.LoadConfig("config.yml")
	if err != nil {
		return
	}
	config.Print()
}
