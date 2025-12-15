package globalping

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type LimitsResponse struct {
	RateLimit struct {
		Measurements struct {
			Create struct {
				Type      string `json:"type"`
				Limit     int    `json:"limit"`
				Remaining int    `json:"remaining"`
				Reset     int    `json:"reset"`
			} `json:"create"`
		} `json:"measurements"`
	} `json:"rateLimit"`
	Credits struct {
		Remaining int `json:"remaining"`
	} `json:"credits"`
}

func (l *LimitsResponse) String() string {
	return fmt.Sprintf("Measurements Create Limit: %d,\nRemaining: %d,\nReset: %d seconds,\nCredits Remaining: %d",
		l.RateLimit.Measurements.Create.Limit,
		l.RateLimit.Measurements.Create.Remaining,
		l.RateLimit.Measurements.Create.Reset,
		l.Credits.Remaining)
}

func Limits(token string) (*LimitsResponse, error) {
	url := "https://api.globalping.io/v1/limits"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Error("Error creating request:", "err", err)
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Error making request:", "err", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Non-OK HTTP status:", "status", resp.Status)
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response body:", "err", err)
		return nil, err
	}

	var limits LimitsResponse
	err = json.Unmarshal(body, &limits)
	if err != nil {
		slog.Error("Error unmarshaling JSON:", "err", err, " body", string(body))
		return nil, err
	}

	return &limits, nil

}
