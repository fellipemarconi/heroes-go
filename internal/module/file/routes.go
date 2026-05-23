package file

import "github.com/go-chi/chi/v5"

func RoutesFile(r chi.Router) {
	service := NewFileService()
	handler := NewFileHandler(service)

	r.Get("/api/file", handler.GetFile)
}
