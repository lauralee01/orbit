package actions

// Phase C.4 — Test plan checklist (comment-only).
//
// Add these tests as you implement C:
//
// 1) Evaluator/dispatcher integration tests
// - Manual trigger path includes trigger_source=manual.
// - Scheduled trigger path includes trigger_source=schedule + schedule_snapshot.
// - Webhook payload fields remain backward-compatible.
//
// 2) Scheduler behavior tests
// - registers enabled rulesets
// - removes disabled/deleted rulesets
// - refresh updates changed cron spec
// - invalid cron/timezone rulesets are skipped (no panic)
//
// 3) Facts provider tests
// - valid data returns expected facts
// - no data returns nil facts => run skipped
// - provider error path is logged and does not crash scheduler
//
// 4) End-to-end smoke (local script is fine)
// - create scheduled ruleset + rule
// - run scheduler with short cron
// - confirm evaluation log + webhook output
//
// 5) Regression tests
// - POST /api/evaluate manual path unchanged
// - Existing webhook payload fields still present
//
// Definition of done for C tests:
// - At least one automated test per component (scheduler, dispatcher, provider).
// - One documented manual E2E flow you can run in under 5 minutes.
