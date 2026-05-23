package admin

import (
	"api/internal/pkg/apierror"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	adminSignatureHeader       = "X-Admin-Signature"
	adminTimestampHeader       = "X-Admin-Timestamp"
	adminSignatureMaxClockSkew = 5 * time.Minute
)

func requireAdminSignature(publicKey ed25519.PublicKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			timestampStr := r.Header.Get(adminTimestampHeader)
			signature := r.Header.Get(adminSignatureHeader)
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
			if signedAt.After(now.Add(adminSignatureMaxClockSkew)) ||
				signedAt.Before(now.Add(-adminSignatureMaxClockSkew)) {
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

func loadAdminPublicKey() (ed25519.PublicKey, error) {
	path := os.Getenv("ADMIN_PUBLIC_KEY_PATH")
	if path == "" {
		return nil, fmt.Errorf("ADMIN_PUBLIC_KEY_PATH not set")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read admin public key: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("invalid admin public key PEM")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse admin public key: %w", err)
	}

	publicKey, ok := key.(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("admin public key is not ed25519")
	}

	return publicKey, nil
}
