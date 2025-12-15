package globalping

import (
	"encoding/json"
	"fmt"
	"globalping/internal/config"
	"io"
	"log/slog"
	"net/http"
	"time"
)

func makeMeasurement(config config.Config) (map[string]any, error) {
	measurementResponse, err := initiateMeasurement(
		config)

	if err != nil {
		slog.Error("initiating measurement:", "err", err)
		return nil, err
	}

	measurementID, ok := measurementResponse["id"].(string)
	if !ok {
		slog.Error("Measurement ID not found", "measuremetResponse", measurementResponse)
		return nil, fmt.Errorf("measurement ID not found in response")
	}

	inProgress := true
	var result map[string]any
	for inProgress {
		time.Sleep(1 * time.Second)
		result, err = requestResult(measurementID, config.Globalping_token)
		if err != nil {
			slog.Error("requesting result:", "err", err)
			return nil, err
		}
		inProgress = result["status"] == "in-progress"
	}
	return result, nil
}
func requestResult(measurementID string, token string) (map[string]any, error) {

	apiURL := fmt.Sprintf("https://api.globalping.io/v1/measurements/%s", measurementID)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", token)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return result, nil
}
