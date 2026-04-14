package handlers

import (
	"encoding/json"
	"net/http"
)

// Limit how large a request body can be (protects the server from huge uploads).
const maxBodyBytes = 1 << 20 // 1 MiB

// echoRequest is the JSON we expect from the client. Field names after `json:"..."`
// are the keys in JSON (like Nest DTOs with @Expose / class-transformer).
type echoRequest struct {
	Message string `json:"message"`
}

// echoResponse is the JSON we send back.
type echoResponse struct {
	Message string `json:"message"`
	OK      bool   `json:"ok"`
}

// errorResponse is a small, consistent shape for API errors.
type errorResponse struct {
	Error  string `json:"error"`
	Detail string `json:"detail,omitempty"`
}

// Echo handles POST /api/echo: read JSON body → struct → write JSON response.
func Echo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
		return
	}

	// Wrap the body so reads cannot exceed maxBodyBytes (then Decode stops with an error).
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	defer r.Body.Close()

	var req echoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON", Detail: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, echoResponse{Message: req.Message, OK: true})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
