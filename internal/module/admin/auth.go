package admin

import (
	"api/internal/infra/auth"
	"crypto/ed25519"
	"net/http"
)

func requireAdminSignature(publicKey ed25519.PublicKey) func(http.Handler) http.Handler {
	return auth.RequireAdminSignature(publicKey)
}

func loadAdminPublicKey() (ed25519.PublicKey, error) {
	return auth.LoadAdminPublicKey()
}
