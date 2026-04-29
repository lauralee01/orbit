package actions

// Phase C.3 — Runtime reliability tasks (comment-only).
//
// Focus:
// - Scheduler lifecycle
// - Webhook delivery policy
// - Operational observability
//
// C3.1 Scheduler lifecycle in cmd/orbit/main.go
// - Replace context.Background() with a cancellable root context.
// - Handle OS signals (SIGINT/SIGTERM) for graceful shutdown.
// - On shutdown:
//   - stop HTTP server gracefully
//   - cancel scheduler context so cron loop exits cleanly
//
// C3.2 Webhook policy decision
// - Current behavior: best-effort single attempt + log.
// - Decide and document one policy:
//   - no retries (simple)
//   - fixed retries with small backoff
//   - queued retries (later)
// - If retries are added, define idempotency key shape early.
//
// C3.3 Logging and traceability
// - Standardize log fields:
//   ruleset_id, trigger_source, ok, reason, webhook_status
// - Avoid duplicate/non-structured noisy lines.
// - Keep logs useful for "why didn't my scheduled job notify?" debugging.
//
// C3.4 Config knobs (optional, minimal)
// - scheduler refresh interval (default 1m)
// - webhook timeout
// - retry enable/disable
//
// Exit criteria:
// - Graceful shutdown works.
// - Delivery behavior is explicit and stable.
// - You can diagnose scheduler and webhook paths from logs.
