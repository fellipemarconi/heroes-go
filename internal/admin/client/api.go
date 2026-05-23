package client

import (
	"api/internal/module/admin"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

const (
	adminSignatureHeader = "X-Admin-Signature"
	adminTimestampHeader = "X-Admin-Timestamp"
	adminStatsPath       = "/api/admin/stats"
	adminUsersPath       = "/api/admin/users"
)

var (
	adminPrivateKeyOnce sync.Once
	adminPrivateKey     ed25519.PrivateKey
	adminPrivateKeyErr  error
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
	req.Header.Set(adminTimestampHeader, timestamp)
	req.Header.Set(adminSignatureHeader, signature)

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
	req.Header.Set(adminTimestampHeader, timestamp)
	req.Header.Set(adminSignatureHeader, signature)

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

	req.Header.Set(adminTimestampHeader, timestamp)
	req.Header.Set(adminSignatureHeader, signature)

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

func signAdminRequest(method string, path string) (string, string, error) {
	privateKey, err := loadAdminPrivateKey()
	if err != nil {
		return "", "", err
	}

	timestamp := time.Now().Unix()
	payload := fmt.Sprintf("%s\n%s\n%d", method, path, timestamp)
	signature := ed25519.Sign(privateKey, []byte(payload))

	return strconv.FormatInt(timestamp, 10), base64.StdEncoding.EncodeToString(signature), nil
}

func loadAdminPrivateKey() (ed25519.PrivateKey, error) {
	adminPrivateKeyOnce.Do(func() {
		path := os.Getenv("ADMIN_PRIVATE_KEY_PATH")
		if path == "" {
			adminPrivateKeyErr = fmt.Errorf("ADMIN_PRIVATE_KEY_PATH not set")
			return
		}

		data, err := os.ReadFile(path)
		if err != nil {
			adminPrivateKeyErr = fmt.Errorf("failed to read admin private key: %w", err)
			return
		}

		block, _ := pem.Decode(data)
		if block == nil {
			adminPrivateKeyErr = fmt.Errorf("invalid admin private key PEM")
			return
		}

		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			adminPrivateKeyErr = fmt.Errorf("failed to parse admin private key: %w", err)
			return
		}

		privateKey, ok := key.(ed25519.PrivateKey)
		if !ok {
			adminPrivateKeyErr = fmt.Errorf("admin private key is not ed25519")
			return
		}

		adminPrivateKey = privateKey
	})

	return adminPrivateKey, adminPrivateKeyErr
}
