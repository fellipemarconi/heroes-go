package file

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewFileHandler(t *testing.T) {
	service := &Service{}
	handler := NewFileHandler(service)

	if handler == nil {
		t.Fatal("NewFileHandler() returned nil")
	}

	if handler.service != service {
		t.Errorf("NewFileHandler() service = %p, want %p", handler.service, service)
	}
}

func TestHandler_GetFile(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodGet, "/api/file", nil)
	rec := httptest.NewRecorder()

	h.GetFile(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("GetFile() status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
