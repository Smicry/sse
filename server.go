package sse

import (
	"errors"
	"net/http"
)

// Upgrader specifies parameters for upgrading an HTTP connection to a
// SSE connection.
//
// It is safe to call Upgrader's methods concurrently.
type Upgrader struct {
}

func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request) (http.Flusher, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, errors.New("streaming unsupported")
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	return flusher, nil
}
