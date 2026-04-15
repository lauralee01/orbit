package handlers

import (
	"encoding/json"
	"net/http"
)

// Limit how large a request body can be (protects the server from huge uploads).
const maxBodyBytes = 1 << 20 // 1 MiB

// errorResponse is a small, consistent shape for API errors.
type errorResponse struct {
	Error  string `json:"error"`
	Detail string `json:"detail,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
