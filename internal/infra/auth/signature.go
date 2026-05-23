package auth

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"api/internal/pkg/apierror"
)

const (
	AdminSignatureHeader       = "X-Admin-Signature"
	AdminTimestampHeader       = "X-Admin-Timestamp"
	AdminSignatureMaxClockSkew = 5 * time.Minute
)

var (
	adminPrivateKeyOnce sync.Once
	adminPrivateKey     ed25519.PrivateKey
	adminPrivateKeyErr  error

	adminPublicKeyOnce sync.Once
	adminPublicKey     ed25519.PublicKey
	adminPublicKeyErr  error
)

// LoadAdminPrivateKey loads the admin private key from the path specified in ADMIN_PRIVATE_KEY_PATH env var.
// Uses sync.Once to ensure the key is only loaded once.
func LoadAdminPrivateKey() (ed25519.PrivateKey, error) {
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

// LoadAdminPublicKey loads the admin public key from the path specified in ADMIN_PUBLIC_KEY_PATH env var.
// Uses sync.Once to ensure the key is only loaded once.
func LoadAdminPublicKey() (ed25519.PublicKey, error) {
	adminPublicKeyOnce.Do(func() {
		path := os.Getenv("ADMIN_PUBLIC_KEY_PATH")
		if path == "" {
			adminPublicKeyErr = fmt.Errorf("ADMIN_PUBLIC_KEY_PATH not set")
			return
		}

		data, err := os.ReadFile(path)
		if err != nil {
			adminPublicKeyErr = fmt.Errorf("failed to read admin public key: %w", err)
			return
		}

		block, _ := pem.Decode(data)
		if block == nil {
			adminPublicKeyErr = fmt.Errorf("invalid admin public key PEM")
			return
		}

		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			adminPublicKeyErr = fmt.Errorf("failed to parse admin public key: %w", err)
			return
		}

		publicKey, ok := key.(ed25519.PublicKey)
		if !ok {
			adminPublicKeyErr = fmt.Errorf("admin public key is not ed25519")
			return
		}

		adminPublicKey = publicKey
	})

	return adminPublicKey, adminPublicKeyErr
}

// SignAdminRequest creates a timestamp and signature for an admin request.
// Returns the timestamp and base64-encoded signature.
func SignAdminRequest(method string, path string) (string, string, error) {
	privateKey, err := LoadAdminPrivateKey()
	if err != nil {
		return "", "", err
	}

	timestamp := time.Now().Unix()
	payload := fmt.Sprintf("%s\n%s\n%d", method, path, timestamp)
	signature := ed25519.Sign(privateKey, []byte(payload))

	return strconv.FormatInt(timestamp, 10), base64.StdEncoding.EncodeToString(signature), nil
}

// RequireAdminSignature returns a middleware that validates admin request signatures.
// It checks the X-Admin-Signature and X-Admin-Timestamp headers against the provided public key.
func RequireAdminSignature(publicKey ed25519.PublicKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			timestampStr := r.Header.Get(AdminTimestampHeader)
			signature := r.Header.Get(AdminSignatureHeader)
			if timestampStr == "" || signature == "" {
				apierror.Send(w, apierror.New(http.StatusUnauthorized, "admin access denied"))
				return
			}

			timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
			if err != nil {
				apierror.Send(w, apierror.New(http.StatusUnauthorized, "admin access denied"))
				return
			}

			signedAt := time.Unix(timestamp, 0)
			now := time.Now()
			if signedAt.After(now.Add(AdminSignatureMaxClockSkew)) ||
				signedAt.Before(now.Add(-AdminSignatureMaxClockSkew)) {
				apierror.Send(w, apierror.New(http.StatusUnauthorized, "admin access denied"))
				return
			}

			signatureBytes, err := base64.StdEncoding.DecodeString(signature)
			if err != nil {
				apierror.Send(w, apierror.New(http.StatusUnauthorized, "admin access denied"))
				return
			}

			payload := fmt.Sprintf("%s\n%s\n%d", r.Method, r.URL.Path, timestamp)
			if !ed25519.Verify(publicKey, []byte(payload), signatureBytes) {
				apierror.Send(w, apierror.New(http.StatusUnauthorized, "admin access denied"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
