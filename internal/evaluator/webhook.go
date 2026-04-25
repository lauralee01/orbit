package evaluator

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type WebhookPayload struct {
	RulesetID        int64  `json:"ruleset_id"`
	OK               bool   `json:"ok"`
	Reason           string `json:"reason,omitempty"`
	EvaluatedAt      string `json:"evaluated_at"`
	TriggerSource    string `json:"trigger_source"`
	ScheduleSnapshot any    `json:"schedule_snapshot,omitempty"` // optional
}

func SendWebhook(ctx context.Context, url string, payload WebhookPayload) {
	jsonData, _ := json.Marshal(payload)

	postCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(postCtx, http.MethodPost, url, bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf(
			"webhook: ruleset_id=%d trigger=%s error=%v",
			payload.RulesetID,
			payload.TriggerSource,
			err,
		)
		return
	}

	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf(
			"webhook: ruleset_id=%d trigger=%s status=%d",
			payload.RulesetID,
			payload.TriggerSource,
			resp.StatusCode,
		)
	} else {
		resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			log.Printf("webhook: bad status: %s", resp.Status)
		}
	}
}
