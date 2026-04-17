package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/lauralee01/orbit/internal/storage"
	"github.com/robfig/cron/v3"
	"net/http"
	"strings"
	"time"
)

type createRulesetRequest struct {
	Name            string `json:"name"`
	WebhookURL      string `json:"webhook_url,omitempty"`
	ScheduleCron    string `json:"schedule_cron,omitempty"`
	ScheduleTZ      string `json:"schedule_tz,omitempty"`
	ScheduleEnabled bool   `json:"schedule_enabled,omitempty"`
}

type createRulesetResponse struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	WebhookURL      string `json:"webhook_url"`
	ScheduleCron    string `json:"schedule_cron"`
	ScheduleTZ      string `json:"schedule_tz"`
	ScheduleEnabled bool   `json:"schedule_enabled"`
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

		var req createRulesetRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON", Detail: err.Error()})
			return
		}

		if strings.TrimSpace(req.Name) == "" {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "name is required"})
			return
		}
		req.Name = strings.TrimSpace(req.Name)
		req.WebhookURL = strings.TrimSpace(req.WebhookURL)
		req.ScheduleCron = strings.TrimSpace(req.ScheduleCron)
		req.ScheduleTZ = strings.TrimSpace(req.ScheduleTZ)

		if req.ScheduleEnabled && req.ScheduleCron == "" {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "schedule cron is required when schedule enabled is true"})
			return
		}
		if req.ScheduleEnabled && req.ScheduleTZ == "" {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "schedule timezone is required when schedule enabled is true"})
			return
		}

		if req.ScheduleEnabled {
			if _, err := cron.ParseStandard(req.ScheduleCron); err != nil {
				writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid schedule cron", Detail: err.Error()})
				return
			}
			if _, err := time.LoadLocation(req.ScheduleTZ); err != nil {
				writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid schedule timezone", Detail: err.Error()})
				return
			}
		}

		id, err := storage.CreateRuleset(r.Context(), db, req.Name, req.WebhookURL, req.ScheduleCron, req.ScheduleTZ, req.ScheduleEnabled)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to create ruleset", Detail: err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, createRulesetResponse{ID: id, Name: req.Name, WebhookURL: req.WebhookURL, ScheduleCron: req.ScheduleCron, ScheduleTZ: req.ScheduleTZ, ScheduleEnabled: req.ScheduleEnabled})
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
