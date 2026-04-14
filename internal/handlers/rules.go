package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/lauralee01/orbit/internal/storage"
	"net/http"
	"strconv"
)

type createRuleRequest struct {
	RulesetID int64  `json:"ruleset_id"`
	Field     string `json:"field"`
	Operator  string `json:"operator"`
	Value     string `json:"value"`
}

type createRuleResponse struct {
	ID int64 `json:"id"`
}

type listRulesResponse struct {
	Rules []storage.StoredRule `json:"rules"`
}

func CreateRule(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
		defer r.Body.Close()

		var req createRuleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON", Detail: err.Error()})
			return
		}
		if req.RulesetID <= 0 || req.Field == "" || req.Operator == "" || req.Value == "" {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request"})
			return
		}
		id, err := storage.InsertRule(r.Context(), db, req.RulesetID, req.Field, req.Operator, req.Value)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to create rule", Detail: err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, createRuleResponse{ID: id})
	}
}

func ListRules(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
			return
		}
		rulesetID, err := strconv.ParseInt(r.URL.Query().Get("ruleset_id"), 10, 64)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid ruleset_id"})
			return
		}
		rules, err := storage.ListRulesByRulesetID(r.Context(), db, rulesetID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to list rules", Detail: err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, listRulesResponse{Rules: rules})
	}
}
