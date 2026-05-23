package file

import (
	"api/internal/infra/storage"
	"api/internal/pkg/apierror"
	"errors"
	"net/http"
	"strings"
)

type Handler struct {
	service *Service
}

func NewFileHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetFile(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(
		r.URL.Query().Get("path"),
		"/",
	)
	if path == "" {
		apierror.Send(w, apierror.ErrInvalidQuery)
		return
	}

	object, info, err := h.service.GetFile(r.Context(), path)
	if err != nil {
		if errors.Is(err, storage.ErrFileNotFound) {
			apierror.Send(w, apierror.ErrFileNotFound)
			return
		}
		apierror.Send(w, err)
		return
	}
	defer func() {
		_ = object.Close()
	}()

	if info.ContentType != "" {
		w.Header().Set("Content-Type", info.ContentType)
	}

	http.ServeContent(w, r, info.Key, info.LastModified, object)
}
