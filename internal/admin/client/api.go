package client

import (
	"api/internal/infra/auth"
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

const (
	adminStatsPath      = "/api/admin/stats"
	adminUsersPath      = "/api/admin/users"
	adminContainersPath = "/api/admin/containers"
	adminRestartPath    = "/api/admin/containers/restart"
	adminLogsPath       = "/api/admin/logs"
)

func GetStats() (*admin.Stats, error) {
	_ = godotenv.Load()

	baseURL := strings.TrimRight(os.Getenv("ADMIN_API_URL"), "/")
	if baseURL == "" {
		return nil, fmt.Errorf("ADMIN_API_URL not set")
	}

	timestamp, signature, err := signAdminRequest(http.MethodGet, adminStatsPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodGet,
		baseURL+adminStatsPath,
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set(auth.AdminTimestampHeader, timestamp)
	req.Header.Set(auth.AdminSignatureHeader, signature)

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

func GetUsers() ([]admin.UserSummary, error) {
	_ = godotenv.Load()

	baseURL := strings.TrimRight(os.Getenv("ADMIN_API_URL"), "/")
	if baseURL == "" {
		return nil, fmt.Errorf("ADMIN_API_URL not set")
	}

	timestamp, signature, err := signAdminRequest(http.MethodGet, adminUsersPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodGet,
		baseURL+adminUsersPath,
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set(auth.AdminTimestampHeader, timestamp)
	req.Header.Set(auth.AdminSignatureHeader, signature)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if len(body) > 0 {
			return nil, fmt.Errorf(
				"users request failed with status %d: %s",
				resp.StatusCode,
				strings.TrimSpace(string(body)),
			)
		}
		return nil, fmt.Errorf(
			"users request failed with status %d",
			resp.StatusCode,
		)
	}

	var users []admin.UserSummary
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	return users, nil
}

func DeleteUser(userID string) error {
	_ = godotenv.Load()

	baseURL := strings.TrimRight(os.Getenv("ADMIN_API_URL"), "/")
	if baseURL == "" {
		return fmt.Errorf("ADMIN_API_URL not set")
	}

	path := fmt.Sprintf("%s/%s", adminUsersPath, userID)
	timestamp, signature, err := signAdminRequest(http.MethodDelete, path)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodDelete,
		baseURL+path,
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set(auth.AdminTimestampHeader, timestamp)
	req.Header.Set(auth.AdminSignatureHeader, signature)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		if len(body) > 0 {
			return fmt.Errorf(
				"delete user failed with status %d: %s",
				resp.StatusCode,
				strings.TrimSpace(string(body)),
			)
		}
		return fmt.Errorf(
			"delete user failed with status %d",
			resp.StatusCode,
		)
	}

	return nil
}

func GetContainers() ([]admin.ContainerSummary, error) {
	_ = godotenv.Load()

	baseURL := strings.TrimRight(os.Getenv("ADMIN_API_URL"), "/")
	if baseURL == "" {
		return nil, fmt.Errorf("ADMIN_API_URL not set")
	}

	timestamp, signature, err := signAdminRequest(http.MethodGet, adminContainersPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodGet,
		baseURL+adminContainersPath,
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set(auth.AdminTimestampHeader, timestamp)
	req.Header.Set(auth.AdminSignatureHeader, signature)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if len(body) > 0 {
			return nil, fmt.Errorf(
				"containers request failed with status %d: %s",
				resp.StatusCode,
				strings.TrimSpace(string(body)),
			)
		}
		return nil, fmt.Errorf(
			"containers request failed with status %d",
			resp.StatusCode,
		)
	}

	var containers []admin.ContainerSummary
	if err := json.NewDecoder(resp.Body).Decode(&containers); err != nil {
		return nil, err
	}

	return containers, nil
}

func RestartContainers() error {
	_ = godotenv.Load()

	baseURL := strings.TrimRight(os.Getenv("ADMIN_API_URL"), "/")
	if baseURL == "" {
		return fmt.Errorf("ADMIN_API_URL not set")
	}

	timestamp, signature, err := signAdminRequest(http.MethodPost, adminRestartPath)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		baseURL+adminRestartPath,
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set(auth.AdminTimestampHeader, timestamp)
	req.Header.Set(auth.AdminSignatureHeader, signature)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		if len(body) > 0 {
			return fmt.Errorf(
				"restart containers failed with status %d: %s",
				resp.StatusCode,
				strings.TrimSpace(string(body)),
			)
		}
		return fmt.Errorf(
			"restart containers failed with status %d",
			resp.StatusCode,
		)
	}

	return nil
}

func GetLogs(limit int) ([]admin.LogEntry, error) {
	_ = godotenv.Load()

	baseURL := strings.TrimRight(os.Getenv("ADMIN_API_URL"), "/")
	if baseURL == "" {
		return nil, fmt.Errorf("ADMIN_API_URL not set")
	}

	path := adminLogsPath
	if limit > 0 {
		path = fmt.Sprintf("%s?limit=%d", adminLogsPath, limit)
	}

	timestamp, signature, err := signAdminRequest(http.MethodGet, adminLogsPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodGet,
		baseURL+path,
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set(auth.AdminTimestampHeader, timestamp)
	req.Header.Set(auth.AdminSignatureHeader, signature)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if len(body) > 0 {
			return nil, fmt.Errorf(
				"logs request failed with status %d: %s",
				resp.StatusCode,
				strings.TrimSpace(string(body)),
			)
		}
		return nil, fmt.Errorf(
			"logs request failed with status %d",
			resp.StatusCode,
		)
	}

	var logs []admin.LogEntry
	if err := json.NewDecoder(resp.Body).Decode(&logs); err != nil {
		return nil, err
	}

	return logs, nil
}

func signAdminRequest(method string, path string) (string, string, error) {
	return auth.SignAdminRequest(method, path)
}
