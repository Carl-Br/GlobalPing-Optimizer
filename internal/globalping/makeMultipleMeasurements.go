package globalping

import (
	"encoding/json"
	"fmt"
	"globalping/internal/config"
	"log/slog"
	"os"
	"strings"
	"time"
)

func appendToJSONL(data map[string]any, filename string) error {
	// append or create new
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

// writeFormattedJSON writes all collected results as a single pretty-printed
// JSON array, so the results can be opened/formatted by regular JSON tools
// (which expect one JSON document, unlike line-delimited JSONL).
func writeFormattedJSON(results []map[string]any, jsonlPath string) (string, error) {
	jsonPath := strings.TrimSuffix(jsonlPath, ".jsonl") + ".json"

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(jsonPath, data, 0644); err != nil {
		return "", err
	}

	return jsonPath, nil
}

// returns the path to the results file
func MakeMultipleMeasurements(config config.Config) (string, error) {
	url_simplified := strings.ReplaceAll(config.Target_url, ".", "-")
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	resultFilePath := fmt.Sprintf("./results/%s_%s.jsonl", url_simplified, timestamp)
	if err := os.MkdirAll("./results", 0755); err != nil {
		slog.Error("creating directory", "err", err)
		return "", err
	}
	startTime := time.Now()
	results := make([]map[string]any, 0, config.Number_measurements)
	for i := range config.Number_measurements {
		result, err := makeMeasurement(config)
		if err != nil {
			return "", err
		}
		if err := appendToJSONL(result, resultFilePath); err != nil {
			slog.Error("appending to jsonl", "err", err)
		}
		results = append(results, result)
		// Show progress
		remainingIterations := config.Number_measurements - (i + 1)
		if remainingIterations > 0 {
			// Estimated time per iteration (measurement + delay)
			elapsed := time.Since(startTime)

			fmt.Printf("\r%-100s\rProgress: %d/%d completed | Remaining: %d | Elapsed: %s",
				"", // Leert die Zeile erst
				i+1,
				config.Number_measurements,
				remainingIterations,
				elapsed.Round(time.Second))

			time.Sleep(config.Seconds_between_measurements)
		} else {
			// Last iteration completed
			totalElapsed := time.Since(startTime)
			fmt.Printf("\r%-100s\rMeasurements completed: %d/%d | Total duration: %s\n",
				"", // Leert die Zeile erst
				config.Number_measurements,
				config.Number_measurements,
				totalElapsed.Round(time.Second))
		}
	}

	if jsonPath, err := writeFormattedJSON(results, resultFilePath); err != nil {
		slog.Error("writing formatted json", "err", err)
	} else {
		fmt.Printf("Formatted JSON saved to: %s\n", jsonPath)
	}

	return resultFilePath, nil
}
