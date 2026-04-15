package storage

import (
	"context"
	"database/sql"
	"fmt"
)

type StoredRuleset struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	WebhookURL string `json:"webhook_url"`
}

func CreateRuleset(ctx context.Context, db *sql.DB, name string, webhookURL string) (int64, error) {
	var id int64
	err := db.QueryRowContext(ctx, `INSERT INTO rulesets (name, webhook_url) VALUES ($1, $2) RETURNING id`, name, webhookURL).Scan(&id)
	return id, err
}

func ListRulesets(ctx context.Context, db *sql.DB) ([]StoredRuleset, error) {
	rows, err := db.QueryContext(ctx, `SELECT id, name, COALESCE(webhook_url, '') FROM rulesets`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []StoredRuleset
	for rows.Next() {
		var r StoredRuleset
		if err := rows.Scan(&r.ID, &r.Name, &r.WebhookURL); err != nil {
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

func UpdateRuleset(ctx context.Context, db *sql.DB, id int64, name string, webhookURL string) error {
	_, err := db.ExecContext(ctx, `UPDATE rulesets SET name = $1, webhook_url = $2 WHERE id = $3`, name, webhookURL, id)
	return err
}

func GetRulesetByID(ctx context.Context, db *sql.DB, id int64) (StoredRuleset, error) {
	var r StoredRuleset
	err := db.QueryRowContext(ctx, `SELECT id, name, COALESCE(webhook_url, '') FROM rulesets WHERE id = $1`, id).Scan(&r.ID, &r.Name, &r.WebhookURL)
	if err != nil {
		return StoredRuleset{}, fmt.Errorf("failed to get ruleset: %w", err)
	}
	return r, nil
}
