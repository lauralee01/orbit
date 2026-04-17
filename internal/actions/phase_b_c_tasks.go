package actions

// Phase B and C roadmap (comment-only guide).
//
// Purpose:
// - Keep implementation order explicit.
// - Let you work incrementally without guessing scope.
// - Preserve "ship small, then expand" behavior.
//
// ---------------------------------------------------------------------------
// PHASE B: Scheduling literacy
// ---------------------------------------------------------------------------
//
// Goal:
// Add first-class scheduling metadata so rulesets can be time-aware.
// This phase is mainly data model + validation + API ergonomics.
// Do NOT add background workers yet (that can come later).
//
// Why this phase:
// - You need a standard way to express "when this ruleset should run/be valid".
// - You can later use this metadata for cron, polling, or queue triggers.
//
// B1) Schema design (DB migration)
// - Add schedule-related columns on rulesets.
//   Suggested first pass:
//   - schedule_cron TEXT NULL
//   - schedule_tz TEXT NOT NULL DEFAULT 'UTC'
//   - schedule_enabled BOOLEAN NOT NULL DEFAULT FALSE
// - Keep nullable/disabled by default so existing rows remain valid.
// - Create migration file (e.g. migrations/003_schedule.sql).
//
// B2) Storage model updates
// - Extend storage.StoredRuleset with new fields.
// - Update ListRulesets/GetRulesetByID SELECTs to include new columns.
// - Use COALESCE where useful for backward safety in scans.
// - Update CreateRuleset/UpdateRuleset function signatures when ready.
//
// B3) API contract updates
// - Extend create/update request JSON models in handlers:
//   - schedule_cron (optional string)
//   - schedule_tz (optional string, default UTC)
//   - schedule_enabled (optional bool, default false)
// - Return schedule fields in ruleset responses.
//
// B4) Validation rules (important)
// - If schedule_enabled == true, schedule_cron must be non-empty.
// - Validate cron syntax with a library or minimal parser check.
// - Validate timezone with time.LoadLocation(schedule_tz).
// - Return 400 with clear JSON errors for invalid schedule input.
//
// B5) Evaluate behavior gate (minimal)
// - For now: decide one clear behavior and document it:
//   Option A: evaluate always, schedule is metadata only.
//   Option B: reject evaluate when schedule_enabled and "outside schedule".
// - Recommend Option A initially to avoid hidden runtime surprises.
//
// B6) Tests for Phase B
// - Handler tests:
//   - create ruleset with valid schedule fields
//   - invalid cron -> 400
//   - invalid timezone -> 400
// - Storage tests:
//   - round-trip read/write schedule fields
// - Regression:
//   - rulesets without schedule still work
//
// B7) Docs updates
// - README API examples: include schedule fields in create/list examples.
// - Mention defaults and validation behavior.
//
// Exit criteria for B:
// - Migration applied cleanly.
// - Ruleset create/list APIs expose schedule fields.
// - Validation and tests are in place.
// - Existing clients remain compatible.
//
// ---------------------------------------------------------------------------
// PHASE C: Glue (wiring scheduling + webhooks + evaluate intent)
// ---------------------------------------------------------------------------
//
// Goal:
// Connect scheduling metadata + evaluation + notifications in a predictable way.
// This phase is orchestration and boundaries, not just schema.
//
// Why this phase:
// - Orbit should become "policy + timing + output" instead of only "policy now".
//
// C1) Define trigger modes (explicit contract)
// - Manual mode: existing POST /api/evaluate behavior.
// - Scheduled mode: system-triggered evaluations (future worker/job).
// - Ensure both paths share the same evaluation core logic.
//
// C2) Add evaluation context fields
// - Add optional metadata in evaluate request/internal context:
//   - trigger_source: "manual" | "schedule"
//   - triggered_at (server-generated)
// - Include this context in webhook payload for observability.
//
// C3) Webhook payload evolution
// - Current payload has ruleset_id, ok, evaluated_at, reason.
// - Extend carefully (backward compatible), e.g.:
//   - trigger_source
//   - schedule_snapshot (optional)
// - Keep old fields stable; only add fields.
//
// C4) Safety and idempotency decisions
// - Decide retry strategy for webhook failures (none, fixed retries, queue later).
// - If retries are added, define idempotency key shape early.
// - Keep initial version simple and logged.
//
// C5) Future scheduler seam (no heavy infra yet)
// - Introduce a small interface boundary for scheduled trigger producers.
//   Example: type TriggerSource interface { EnqueueRulesetEval(...) error }
// - Start with no-op/manual implementation to avoid lock-in.
//
// C6) Observability baseline
// - Standardize log fields for evaluate + webhook path:
//   ruleset_id, trigger_source, outcome, webhook_status.
// - Helps debugging once automatic triggers start running.
//
// C7) Tests for C
// - Evaluate flow includes trigger context consistently.
// - Webhook payload includes new optional fields.
// - Backward compatibility: old payload consumers are unaffected.
//
// Exit criteria for C:
// - Manual and scheduled paths have a shared core evaluation contract.
// - Webhook payload includes enough context to debug outcomes.
// - Design is ready for adding a real scheduler worker next.
