package logs

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Entry struct {
	Timestamp  time.Time `json:"timestamp"`
	Method     string    `json:"method"`
	Path       string    `json:"path"`
	Status     int       `json:"status"`
	DurationMs int64     `json:"duration_ms"`
	Error      string    `json:"error,omitempty"`
	Location   string    `json:"location,omitempty"`
}

const bufferSize = 500

var store = newStore(bufferSize)

func List(limit int) []Entry {
	return store.List(limit)
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		duration := time.Since(start)
		entry := Entry{
			Timestamp:  time.Now(),
			Method:     r.Method,
			Path:       r.URL.Path,
			Status:     recorder.status,
			DurationMs: duration.Milliseconds(),
			Error:      recorder.errMessage,
			Location:   recorder.errLocation,
		}

		store.Add(entry)
		log.Printf("%s %s -> %d (%dms)", r.Method, r.URL.Path, recorder.status, duration.Milliseconds())
	})
}

type responseRecorder struct {
	http.ResponseWriter
	status      int
	errMessage  string
	errLocation string
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) SetError(message string, location string) {
	r.errMessage = message
	r.errLocation = location
}

type storeBuffer struct {
	mu      sync.Mutex
	entries []Entry
	next    int
	count   int
}

func newStore(size int) *storeBuffer {
	return &storeBuffer{
		entries: make([]Entry, size),
	}
}

func (s *storeBuffer) Add(entry Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.entries[s.next] = entry
	s.next = (s.next + 1) % len(s.entries)
	if s.count < len(s.entries) {
		s.count++
	}
}

func (s *storeBuffer) List(limit int) []Entry {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.count == 0 {
		return []Entry{}
	}

	if limit <= 0 || limit > s.count {
		limit = s.count
	}

	items := make([]Entry, 0, limit)
	for i := 0; i < limit; i++ {
		index := s.next - 1 - i
		if index < 0 {
			index += len(s.entries)
		}
		items = append(items, s.entries[index])
	}

	for i := 0; i < len(items)/2; i++ {
		j := len(items) - 1 - i
		items[i], items[j] = items[j], items[i]
	}

	for i := range items {
		items[i].Path = strings.TrimSpace(items[i].Path)
	}

	return items
}
