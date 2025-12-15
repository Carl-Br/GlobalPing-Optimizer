package globalping

import (
	"bytes"
	"encoding/json"
	"fmt"
	"globalping/internal/config"
	"io"
	"log/slog"
	"net/http"
)

func initiateMeasurement(config config.Config) (map[string]any, error) {
	apiURL := "https://api.globalping.io/v1/measurements"
	locationMap := make([]map[string]string, len(config.Locations))
	for i, loc := range config.Locations {
		locationMap[i] = map[string]string{"magic": loc}
	}

	data := map[string]any{
		"limit":             fmt.Sprintf("%d", config.LimitPerMeasurement),
		"inProgressUpdates": false,
		"locations":         locationMap,
		"measurementOptions": map[string]any{
			"request":   map[string]string{"method": "GET"},
			"ipVersion": 4,
		},
		"target": config.Target_url,
		"type":   "HTTP",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		slog.Error("Failed to marshal JSON", "error", err)
		return nil, err
	}
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		slog.Error("Failed to create HTTP request", "error", err)
		return nil, err
	}
	req.Header.Set("Authorization", config.Globalping_token)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Failed to perform HTTP request", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}
