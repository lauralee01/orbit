package storage 

import (
	"context"
	"database/sql"
	"fmt"
)

type StoredRule struct {
	ID int64 `json:"id"`
	RulesetID int64 `json:"ruleset_id"`
	Field string `json:"field"`
	Operator string `json:"operator"`
	Value string `json:"value"`
}

func InsertRule(ctx context.Context, db *sql.DB, rulesetID int64, field string, operator string, value string) (int64, error) {
	var id int64
	err := db.QueryRowContext(ctx, `INSERT INTO rules (ruleset_id, field, operator, value) VALUES ($1, $2, $3, $4) RETURNING id`, rulesetID, field, operator, value).Scan(&id)
	return id, err
}

func ListRulesByRulesetID(ctx context.Context, db *sql.DB, rulesetID int64) ([]StoredRule, error) {
	rows, err := db.QueryContext(ctx, `SELECT id, ruleset_id, field, operator, value FROM rules WHERE ruleset_id = $1`, rulesetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rules: %w", err)
	}
	defer rows.Close()
	var out []StoredRule
	for rows.Next() {
		var r StoredRule
		if err := rows.Scan(&r.ID, &r.RulesetID, &r.Field, &r.Operator, &r.Value); err != nil {
			return nil, fmt.Errorf("failed to scan rule: %w", err)
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get rules: %w", err)
	}
	return out, nil 
}

