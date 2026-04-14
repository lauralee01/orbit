package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/lauralee01/orbit/internal/storage"
	"log"
	"net/http"
)

type createRulesetRequest struct {
	Name string `json:"name"`
}

type createRulesetResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type listRulesetsResponse struct {
	Rulesets []storage.StoredRuleset `json:"rulesets"`
}

func CreateRuleset(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

		defer r.Body.Close()

		log.Printf("CreateRuleset: %v before decode", r.Body)

		var req createRulesetRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON", Detail: err.Error()})
			return
		}
		log.Printf("CreateRuleset: %v after decode", req)
		if req.Name == "" {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "name is required"})
			return
		}

		id, err := storage.CreateRuleset(r.Context(), db, req.Name)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to create ruleset", Detail: err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, createRulesetResponse{ID: id, Name: req.Name})
	}
}

func ListRulesets(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
			return
		}

		rulesets, err := storage.ListRulesets(r.Context(), db)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to list rulesets", Detail: err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, listRulesetsResponse{Rulesets: rulesets})
	}
}

// # Step H (later)
//
// Add InsertRule / list rules for a ruleset, then POST/GET for rules — same patterns.
