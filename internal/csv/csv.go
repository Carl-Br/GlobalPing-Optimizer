package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type MeasurementData struct {
	CreatedAt string `json:"createdAt"`
	ID        string `json:"id"`
	Limit     int    `json:"limit"`
	Locations []struct {
		Magic string `json:"magic"`
	} `json:"locations"`
	MeasurementOptions struct {
		Request struct {
			Method string `json:"method"`
		} `json:"request"`
	} `json:"measurementOptions"`
	ProbesCount int `json:"probesCount"`
	Results     []struct {
		Probe struct {
			ASN       int      `json:"asn"`
			City      string   `json:"city"`
			Continent string   `json:"continent"`
			Country   string   `json:"country"`
			Latitude  float64  `json:"latitude"`
			Longitude float64  `json:"longitude"`
			Network   string   `json:"network"`
			Region    string   `json:"region"`
			Resolvers []string `json:"resolvers"`
			State     *string  `json:"state"`
			Tags      []string `json:"tags"`
		} `json:"probe"`
		Result struct {
			Headers         map[string]any `json:"headers"`
			RawBody         string         `json:"rawBody"`
			RawHeaders      string         `json:"rawHeaders"`
			RawOutput       string         `json:"rawOutput"`
			ResolvedAddress string         `json:"resolvedAddress"`
			Status          string         `json:"status"`
			StatusCode      int            `json:"statusCode"`
			StatusCodeName  string         `json:"statusCodeName"`
			Timings         struct {
				DNS       int `json:"dns"`
				Download  int `json:"download"`
				FirstByte int `json:"firstByte"`
				TCP       int `json:"tcp"`
				TLS       int `json:"tls"`
				Total     int `json:"total"`
			} `json:"timings"`
			TLS struct {
				Authorized     bool   `json:"authorized"`
				CipherName     string `json:"cipherName"`
				CreatedAt      string `json:"createdAt"`
				ExpiresAt      string `json:"expiresAt"`
				Fingerprint256 string `json:"fingerprint256"`
				Issuer         struct {
					C  string `json:"C"`
					CN string `json:"CN"`
					O  string `json:"O"`
				} `json:"issuer"`
				KeyBits      int    `json:"keyBits"`
				KeyType      string `json:"keyType"`
				Protocol     string `json:"protocol"`
				PublicKey    string `json:"publicKey"`
				SerialNumber string `json:"serialNumber"`
				Subject      struct {
					CN  string `json:"CN"`
					Alt string `json:"alt"`
				} `json:"subject"`
			} `json:"tls"`
			Truncated bool `json:"truncated"`
		} `json:"result"`
	} `json:"results"`
}

// networkStats represents connection statistics for a network
type networkStats struct {
	Address         string
	City            string
	Country         string
	Network         string
	ConnectionTimes []int
	MinTime         int
	MaxTime         int
	MedianTime      float64
	MeanTime        float64
}

// networkMedian represents a network with its median connection time
type networkMedian struct {
	Address     string
	Location    string
	Network     string
	MedianTime  float64
	MinTime     int
	MaxTime     int
	MeanTime    float64
	SampleCount int
}

func Calculate(resultsFilePath string) {
	measurements := readResults(resultsFilePath)
	networks := groupByCityAndNetwork(measurements)
	statistics := generateNetworkStatistics(networks)

	if len(statistics) > 0 {
		// Save results to CSV file
		filename := strings.Replace(resultsFilePath, ".jsonl", "_stats.csv", 1)
		// Fallback for .json extension (backwards compatibility)
		if filename == resultsFilePath {
			filename = strings.Replace(resultsFilePath, ".json", "_stats.csv", 1)
		}
		saveToCSV(statistics, filename)

		fmt.Printf("\nDetailed statistics saved to: %s\n", filename)
	} else {
		fmt.Println("\nNo valid network data found.")
	}
}

// saveToCSV saves network statistics to a CSV file
func saveToCSV(networks []networkMedian, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"Location",
		"Network",
		"Median (ms)",
		"Mean (ms)",
		"Min (ms)",
		"Max (ms)",
		"Sample Count",
	}

	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data rows
	for _, network := range networks {
		row := []string{
			network.Location,
			network.Network,
			fmt.Sprintf("%.2f", network.MedianTime),
			fmt.Sprintf("%.2f", network.MeanTime),
			strconv.Itoa(network.MinTime),
			strconv.Itoa(network.MaxTime),
			strconv.Itoa(network.SampleCount),
		}

		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
