package facts

import (
    "context"
    "database/sql"
    "encoding/json"
	"github.com/lauralee01/orbit/internal/rules"
)

type Provider interface {
    GetFacts(ctx context.Context, rulesetID int64) (rules.Facts, error)
}

type DBProvider struct {
    DB *sql.DB
}

func NewDBProvider(db *sql.DB) *DBProvider {
    return &DBProvider{DB: db}
}

func (p *DBProvider) GetFacts(ctx context.Context, rulesetID int64) (rules.Facts, error) {
    var raw []byte

    err := p.DB.QueryRowContext(ctx,
        `SELECT facts FROM ruleset_facts WHERE ruleset_id = $1`,
        rulesetID,
    ).Scan(&raw)

    if err == sql.ErrNoRows {
        // No facts available → scheduler should skip evaluation
        return nil, nil
    }
    if err != nil {
        return nil, err
    }

    var f rules.Facts
    if err := json.Unmarshal(raw, &f); err != nil {
        return nil, err
    }

    return f, nil
}
