package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/lauralee01/orbit/internal/evaluator"
	"github.com/lauralee01/orbit/internal/storage"
	"log"
	"net/http"
	"time"
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

		// Load ruleset
		ruleset, err := storage.GetRulesetByID(r.Context(), db, req.RulesetID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to get ruleset", Detail: err.Error()})
			return
		}

		// Build evaluation context
		evalCtx := evaluator.EvalContext{
			TriggerSource: evaluator.TriggerSourceManual,
			TriggerAt:     time.Now(),
		}

		// Shared evaluation path
		result, err := evaluator.EvaluateRuleset(r.Context(), db, ruleset, req.Facts, evalCtx)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
			return
		}

		// Webhook (now using result + evalCtx)
		if ruleset.WebhookURL != "" {
			payload := evaluator.WebhookPayload{
				RulesetID:     ruleset.ID,
				OK:             result.OK,
				Reason:         result.Reason,
				EvaluatedAt:   evalCtx.TriggerAt.Format(time.RFC3339),
				TriggerSource: string(evalCtx.TriggerSource),
			}

			jsonData, _ := json.Marshal(payload)

			postCtx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
			defer cancel()

			req, _ := http.NewRequestWithContext(postCtx, http.MethodPost, ruleset.WebhookURL, bytes.NewReader(jsonData))
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Printf("webhook: post: %v", err)
			} else {
				resp.Body.Close()
				if resp.StatusCode < 200 || resp.StatusCode >= 300 {
					log.Printf("webhook: bad status: %s", resp.Status)
				}
			}
		}

		writeJSON(w, http.StatusOK, evaluateResponse{
			OK:     result.OK,
			Reason: result.Reason,
		})
	}
}
