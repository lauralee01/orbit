package evaluator

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lauralee01/orbit/internal/rules"
	"github.com/lauralee01/orbit/internal/storage"
	"log"
	"time"
)

type TriggerSource string

const (
	TriggerSourceManual   TriggerSource = "manual"
	TriggerSourceSchedule TriggerSource = "schedule"
)

type EvalContext struct {
	TriggerSource TriggerSource
	TriggerAt     time.Time
}

type EvalResult struct {
	OK     bool
	Reason string
}

func EvaluateRuleset(
	ctx context.Context,
	db *sql.DB,
	rs storage.StoredRuleset,
	facts rules.Facts,
	evalCtx EvalContext,
) (EvalResult, error) {

	storedRules, err := storage.ListRulesByRulesetID(ctx, db, rs.ID)
	if err != nil {
		return EvalResult{}, fmt.Errorf("list rules: %w", err)
	}

	ruleSlice := make(rules.Rules, len(storedRules))
	for i, row := range storedRules {
		ruleSlice[i] = rules.Rule{
			Field:    row.Field,
			Operator: row.Operator,
			Value:    row.Value,
		}
	}

	ok, evalErr := rules.Evaluate(facts, ruleSlice)

	result := EvalResult{OK: ok}

	if evalErr != nil {
		result.OK = false
		result.Reason = evalErr.Error()
	}

	log.Printf(
		"evaluate: ruleset_id=%d trigger=%s ok=%v reason=%s",
		rs.ID, evalCtx.TriggerSource, result.OK, result.Reason,
	)

	return result, nil
}
