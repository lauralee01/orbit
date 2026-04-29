package actions

// Phase C overview (comment-only implementation guide).
//
// You already have:
// - Evaluator package with trigger context
// - Scheduler package that can dispatch scheduled evaluations
// - Webhook sender and payload
//
// Remaining Phase C goal:
// Make the "glue" predictable and production-friendly.
//
// Suggested learning order:
// 1) Facts provider contract + one real provider.
// 2) Trigger contract consistency (manual + schedule share one path).
// 3) Scheduler lifecycle (start/stop with app shutdown).
// 4) Webhook delivery policy (retry/no retry, observability fields).
// 5) End-to-end tests proving trigger -> evaluate -> action.
//
// Keep this phase incremental:
// - Ship tiny slices.
// - Keep backward compatibility.
// - Add tests whenever behavior changes.
