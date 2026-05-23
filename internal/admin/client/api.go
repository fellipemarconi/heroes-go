package client

import (
	"api/internal/module/admin"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func GetStats() (*admin.Stats, error) {
	_ = godotenv.Load()

	baseURL := os.Getenv("ADMIN_API_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("ADMIN_API_URL not set")
	}

	req, err := http.NewRequest(
		http.MethodGet,
		baseURL+"/api/admin/stats",
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if len(body) > 0 {
			return nil, fmt.Errorf(
				"stats request failed with status %d: %s",
				resp.StatusCode,
				strings.TrimSpace(string(body)),
			)
		}
		return nil, fmt.Errorf(
			"stats request failed with status %d",
			resp.StatusCode,
		)
	}

	var stats admin.Stats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, err
	}

	return &stats, nil
}
