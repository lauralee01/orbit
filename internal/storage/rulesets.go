package storage

import (
	"context"
	"database/sql"
)

type StoredRuleset struct {
	ID   int64 `json:"id"`
	Name string `json:"name"`
}

func CreateRuleset(ctx context.Context, db *sql.DB, name string) (int64, error) {
	var id int64
	err := db.QueryRowContext(ctx, `INSERT INTO rulesets (name) VALUES ($1) RETURNING id`, name).Scan(&id)
	return id, err
}

func ListRulesets(ctx context.Context, db *sql.DB) ([]StoredRuleset, error) {
	rows, err := db.QueryContext(ctx, `SELECT id, name FROM rulesets`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []StoredRuleset
	for rows.Next() {
		var r StoredRuleset
		if err := rows.Scan(&r.ID, &r.Name); err != nil {
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
