package evaluator

import (
	"context"
	"database/sql"

	"github.com/lauralee01/orbit/internal/rules"
	"github.com/lauralee01/orbit/internal/storage"
	"log"
	"time"
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
	// Build scheduled snapshots only for scheduled runs
	var snapshot any
	if evalCtx.TriggerSource == TriggerSourceSchedule {
		snapshot = map[string]any{
			"cron":     rs.ScheduleCron,
			"timezone": rs.ScheduleTZ,
		}
	}

	// Evaluate rules
	result, err := EvaluateRuleset(ctx, d.DB, rs, facts, evalCtx)
	if err != nil {
		log.Printf("dispatcher: evaluation error: %v", err)
		return
	}

	// Send webhook
	if rs.WebhookURL != "" {
		payload := WebhookPayload{
			RulesetID:        rs.ID,
			OK:               result.OK,
			Reason:           result.Reason,
			EvaluatedAt:      evalCtx.TriggerAt.Format(time.RFC3339),
			TriggerSource:    string(evalCtx.TriggerSource),
			ScheduleSnapshot: snapshot,
		}

		SendWebhook(ctx, rs.WebhookURL, payload)
	}

}
