package config

import (
	"fmt"
	"globalping/internal/util"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Target_url                   string        `yaml:"target_url"`
	Number_measurements          int           `yaml:"number_measurements"`
	Seconds_between_measurements time.Duration `yaml:"seconds_between_measurements"`
	Globalping_token             string
	LimitPerMeasurement          int      `yaml:"limit_per_measurement"`
	Locations                    []string `yaml:"locations"`
}

func LoadConfig(filename string) (*Config, error) {

	// read
	data, err := os.ReadFile(filename)
	if err != nil {
		slog.Error("Failed to read config file", "error", err)
		return nil, err
	}

	// unmarshal
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		slog.Error("Failed to unmarshal config file", "error", err)
		return nil, err
	}

	// token
	if err := godotenv.Load(".env"); err != nil {
		slog.Error("Error loading .env file", "error", err)
		return nil, err
	}
	config.Globalping_token = "Bearer " + os.Getenv("GLOBALPING_TOKEN")

	return &config, nil
}

func (c *Config) Print() {
	fmt.Println(util.BoldTerminalText("Config:"))
	fmt.Println(util.BoldTerminalText("TargetUrl:"), c.Target_url)
	fmt.Println(util.BoldTerminalText("Number_measurements:"), c.Number_measurements)
	fmt.Println(util.BoldTerminalText("Seconds_between_measurements:"), c.Seconds_between_measurements)

	token := "set"
	if c.Globalping_token == "" ||
		c.Globalping_token == "Put_Your_Global_ping_token_here" {
		token = c.Globalping_token
	}
	fmt.Println(util.BoldTerminalText("Globalping_token:"), token)
	fmt.Println(util.BoldTerminalText("LimitPerMeasurement:"), c.LimitPerMeasurement)
	fmt.Println(util.BoldTerminalText("Locations:"), c.Locations)
}
