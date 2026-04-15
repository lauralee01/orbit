package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

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
		ok, evalErr := rules.Evaluate(req.Facts, ruleSlice)

		var evalOK bool
		var evalReason string

		switch {
		case evalErr != nil && (errors.Is(evalErr, rules.ErrMissingFact) || errors.Is(evalErr, rules.ErrFactValueMismatch)):
			evalOK = false
			evalReason = evalErr.Error()
		case evalErr != nil:
			writeJSON(w, http.StatusInternalServerError, errorResponse{
				Error:  "failed to evaluate rules",
				Detail: evalErr.Error(),
			})
			return
		default:
			evalOK = ok
		}

		ruleset, err := storage.GetRulesetByID(r.Context(), db, req.RulesetID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to get ruleset", Detail: err.Error()})
			return
		}

		if ruleset.WebhookURL != "" {
			payload := map[string]any{
				"ruleset_id":   req.RulesetID,
				"ok":           evalOK,
				"evaluated_at": time.Now().Format(time.RFC3339),
			}
			if evalReason != "" {
				payload["reason"] = evalReason
			}
			jsonData, err := json.Marshal(payload)
			if err != nil {
				log.Printf("webhook: marshal: %v", err)
			} else {
				postCtx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
				defer cancel()
				hreq, err := http.NewRequestWithContext(postCtx, http.MethodPost, ruleset.WebhookURL, bytes.NewReader(jsonData))
				if err != nil {
					log.Printf("webhook: new request: %v", err)
				} else {
					hreq.Header.Set("Content-Type", "application/json")
					resp, err := http.DefaultClient.Do(hreq)
					if err != nil {
						log.Printf("webhook: post: %v", err)
					} else {
						resp.Body.Close()
						if resp.StatusCode < 200 || resp.StatusCode >= 300 {
							log.Printf("webhook: bad status: %s", resp.Status)
						}
					}
				}
			}
		}

		writeJSON(w, http.StatusOK, evaluateResponse{OK: evalOK, Reason: evalReason})
	}
}
