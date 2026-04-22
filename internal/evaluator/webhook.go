package evaluator

type WebhookPayload struct {
    RulesetID       int64       `json:"ruleset_id"`
    OK              bool        `json:"ok"`
    Reason          string      `json:"reason,omitempty"`
    EvaluatedAt     string      `json:"evaluated_at"`
    TriggerSource   string      `json:"trigger_source"`
    ScheduleSnapshot any        `json:"schedule_snapshot,omitempty"` // optional 
}
