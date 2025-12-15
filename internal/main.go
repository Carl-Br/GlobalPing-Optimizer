package main

import (
	"fmt"
	"globalping/internal/config"
	"globalping/internal/csv"
	"globalping/internal/globalping"
	"globalping/internal/util"
	"log/slog"
	"os"
)

func main() {
	// logging
	logLevel := slog.LevelInfo
	handlerOpts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, handlerOpts))
	slog.SetDefault(logger)

	config, err := config.LoadConfig("config.yml")
	if err != nil {
		return
	}
	config.Print()
	limits, err := globalping.Limits(config.Globalping_token)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println("\n" + util.BoldTerminalText("Globalping Limits:"))
	fmt.Println(limits)

	totalLimits := limits.Credits.Remaining + limits.RateLimit.Measurements.Create.Remaining
	requiredLimits := config.Number_measurements * config.LimitPerMeasurement

	// check credits
	fmt.Println("\n"+util.BoldTerminalText("Total Limits:"), totalLimits)
	fmt.Println(util.BoldTerminalText("Required Limits/number of requests:"), requiredLimits)
	if totalLimits < requiredLimits {
		fmt.Println("You don't have enough creadits")
		os.Exit(1)
	}

	// start
	fmt.Println("\ns: start, q: quit")

	var input string
	fmt.Print("(s/q): ")
	fmt.Scan(&input)
	if input != "s" {
		os.Exit(1)
	}

	fmt.Println("\n" + util.BoldTerminalText("Making Measurements:"))
	resultFilePath, err := globalping.MakeMultipleMeasurements(*config)
	if err != nil {
		return
	}
	csv.Calculate(resultFilePath)
}
