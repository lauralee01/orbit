package evaluator

import (
    "context"
    "database/sql"

    "github.com/lauralee01/orbit/internal/rules"
    "github.com/lauralee01/orbit/internal/storage"
)

type EvaluationDispatcher interface {
    DispatchEvaluation(ctx context.Context, rs storage.StoredRuleset, facts rules.Facts, evalCtx EvalContext)
}

type DirectDispatcher struct {
    DB *sql.DB
}

func (d *DirectDispatcher) DispatchEvaluation(
    ctx context.Context,
    rs storage.StoredRuleset,
    facts rules.Facts,
    evalCtx EvalContext,
) {
    EvaluateRuleset(ctx, d.DB, rs, facts, evalCtx)
}
