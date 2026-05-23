package ctx

import (
	"api/internal/pkg/apierror"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

func GetUserIDCtx(r *http.Request) (string, error) {
	token, claims, err := jwtauth.FromContext(r.Context())
	if err != nil || token == nil {
		return "", apierror.ErrInvalidToken
	}

	id, ok := claims["user_id"].(string)
	if !ok {
		return "", apierror.ErrInvalidToken
	}

	return id, nil
}
