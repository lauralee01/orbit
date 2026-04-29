package actions

// Phase C.2 — Trigger contract + shared evaluation path (comment-only).
//
// Current state:
// - Manual API uses evaluator.EvaluateRuleset with TriggerSourceManual.
// - Scheduler uses dispatcher with TriggerSourceSchedule.
// - Webhook payload includes trigger_source and optional schedule_snapshot.
//
// Goal:
// - Keep manual and scheduled flows semantically consistent.
//
// Tasks:
// 1) Document canonical trigger context fields:
//    - trigger_source
//    - triggered_at (from EvalContext.TriggerAt)
//    - schedule snapshot only for scheduled runs
//
// 2) Ensure one evaluation core is always used:
//    - Any future trigger (queue, webhook ingest, worker) must call evaluator package.
//    - Avoid duplicated evaluate logic inside handlers/scheduler.
//
// 3) Define expected failure semantics by trigger type:
//    - Manual trigger: API response communicates result immediately.
//    - Scheduled trigger: result is logged + optionally sent to webhook.
//
// 4) Add a short architecture note (README or docs/C.md):
//    - "Orbit core evaluates rules; trigger adapters feed facts + context."
//
// 5) Decide if scheduled "no-op" should still emit webhook:
//    - For now likely no (skip when facts missing).
//    - Document this decision clearly.
//
// Exit criteria:
// - Every trigger path can be reasoned about with one shared contract.
// - Webhook consumers can inspect trigger_source reliably.
