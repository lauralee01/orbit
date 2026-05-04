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

    // Phase C3.2: short timeout for webhook delivery
    postCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()

    req, _ := http.NewRequestWithContext(postCtx, http.MethodPost, url, bytes.NewReader(jsonData))
    req.Header.Set("Content-Type", "application/json")

    // Use a client with a timeout instead of http.DefaultClient
    client := &http.Client{
        Timeout: 3 * time.Second,
    }

    resp, err := client.Do(req)
    if err != nil {
        log.Printf(
            "webhook delivery: ruleset_id=%d trigger=%s status=0 error=%v",
            payload.RulesetID,
            payload.TriggerSource,
            err,
        )
        return
    }
    defer resp.Body.Close()

    // Log structured result
    log.Printf(
        "webhook delivery: ruleset_id=%d trigger=%s status=%d error=%v",
        payload.RulesetID,
        payload.TriggerSource,
        resp.StatusCode,
        nil,
    )
}

