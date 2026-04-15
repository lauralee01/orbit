package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/lauralee01/orbit/internal/rules"
	"github.com/lauralee01/orbit/internal/storage"
)

type evaluateRequest struct {
	RulesetID int64          `json:"ruleset_id"`
	Facts     map[string]any `json:"facts"`
}

type evaluateResponse struct {
	OK     bool   `json:"ok"`
	Reason string `json:"reason,omitempty"`
}

func Evaluate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
		defer r.Body.Close()

		var req evaluateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON", Detail: err.Error()})
			return
		}
		if req.RulesetID <= 0 || req.Facts == nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request"})
			return
		}
		storedRules, err := storage.ListRulesByRulesetID(r.Context(), db, req.RulesetID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to list rules", Detail: err.Error()})
			return
		}
		ruleSlice := make(rules.Rules, len(storedRules))
		for i, row := range storedRules {
			ruleSlice[i] = rules.Rule{Field: row.Field, Operator: row.Operator, Value: row.Value}
		}
		ok, err := rules.Evaluate(req.Facts, ruleSlice)
		if err != nil {
			if errors.Is(err, rules.ErrMissingFact) || errors.Is(err, rules.ErrFactValueMismatch) {
				writeJSON(w, http.StatusOK, evaluateResponse{
					OK:     false,
					Reason: err.Error(),
				})
				return
			}

			// real error
			writeJSON(w, http.StatusInternalServerError, errorResponse{
				Error:  "failed to evaluate rules",
				Detail: err.Error(),
			})
			return
		}
		writeJSON(w, http.StatusOK, evaluateResponse{OK: ok})
	}
}
