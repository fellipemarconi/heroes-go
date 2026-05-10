package ctx

import (
	"api/internal/pkg/apierror"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

func GetUserIDCtx(r *http.Request) (string, error) {
	_, claims, _ := jwtauth.FromContext(r.Context())

	id, ok := claims["user_id"].(string)
	if !ok {
		return "", apierror.Internal(apierror.ErrInternalServer.Error())
	}

	return id, nil
}
