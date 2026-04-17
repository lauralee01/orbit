package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type StoredRuleset struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	WebhookURL      string `json:"webhook_url,omitempty"`
	ScheduleCron    string `json:"schedule_cron,omitempty"`
	ScheduleTZ      string `json:"schedule_tz,omitempty"`
	ScheduleEnabled bool   `json:"schedule_enabled,omitempty"`
}

func CreateRuleset(ctx context.Context, db *sql.DB, name string, webhookURL string, scheduleCron string, scheduleTZ string, scheduleEnabled bool) (int64, error) {
	var id int64
	scheduleCron = strings.TrimSpace(scheduleCron)
	scheduleTZ = strings.TrimSpace(scheduleTZ)
	if scheduleCron == "" {
		scheduleEnabled = false
	}
	if scheduleTZ == "" {
		scheduleTZ = "UTC"
	}
	err := db.QueryRowContext(ctx, `INSERT INTO rulesets (name, webhook_url, schedule_cron, schedule_tz, schedule_enabled) VALUES ($1, $2, $3, $4, $5) RETURNING id`, name, webhookURL, scheduleCron, scheduleTZ, scheduleEnabled).Scan(&id)
	return id, err
}

func ListRulesets(ctx context.Context, db *sql.DB) ([]StoredRuleset, error) {
	rows, err := db.QueryContext(ctx, `SELECT id, name, COALESCE(webhook_url, ''), COALESCE(schedule_cron, ''), COALESCE(schedule_tz, 'UTC'), COALESCE(schedule_enabled, FALSE) FROM rulesets`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []StoredRuleset
	for rows.Next() {
		var r StoredRuleset
		if err := rows.Scan(&r.ID, &r.Name, &r.WebhookURL, &r.ScheduleCron, &r.ScheduleTZ, &r.ScheduleEnabled); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	// (1) rows.Err() after the loop: rows.Next() can return false when the result set
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func UpdateRuleset(ctx context.Context, db *sql.DB, id int64, name string, webhookURL string, scheduleCron string, scheduleTZ string, scheduleEnabled bool) error {
	_, err := db.ExecContext(ctx, `UPDATE rulesets SET name = $1, webhook_url = $2, schedule_cron = $3, schedule_tz = $4, schedule_enabled = $5 WHERE id = $6`, name, webhookURL, scheduleCron, scheduleTZ, scheduleEnabled, id)
	return err
}

func GetRulesetByID(ctx context.Context, db *sql.DB, id int64) (StoredRuleset, error) {
	var r StoredRuleset
	err := db.QueryRowContext(ctx, `SELECT id, name, COALESCE(webhook_url, ''), COALESCE(schedule_cron, ''), COALESCE(schedule_tz, 'UTC'), COALESCE(schedule_enabled, FALSE) FROM rulesets WHERE id = $1`, id).Scan(&r.ID, &r.Name, &r.WebhookURL, &r.ScheduleCron, &r.ScheduleTZ, &r.ScheduleEnabled)
	if err != nil {
		return StoredRuleset{}, fmt.Errorf("failed to get ruleset: %w", err)
	}
	return r, nil
}
