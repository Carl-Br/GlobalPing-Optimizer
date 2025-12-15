package csv

import (
	"bufio"
	"encoding/json"
	"log/slog"
	"os"
	"strings"
)

func readResults(resultsFilePath string) []MeasurementData {
	// Open the JSONL file
	file, err := os.Open(resultsFilePath)
	if err != nil {
		slog.Error("Error opening file", "err", err)
		os.Exit(1)
	}
	defer file.Close()

	// Read JSONL file line by line
	var measurements []MeasurementData
	scanner := bufio.NewScanner(file)

	// Increase the buffer size to handle large lines
	const maxCapacity = 10 * 1024 * 1024 // 10MB per line
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Parse each line as a JSON object
		var measurement MeasurementData
		err := json.Unmarshal([]byte(line), &measurement)
		if err != nil {
			slog.Error("Error parsing JSON line", "err", err, "line", lineNumber, "data", line)
			continue
		}

		measurements = append(measurements, measurement)
	}

	if err := scanner.Err(); err != nil {
		slog.Error("Error reading file", "err", err)
		os.Exit(1)
	}

	if len(measurements) == 0 {
		slog.Error("No valid measurement data found in the file")
		os.Exit(1)
	}
	return measurements
}
