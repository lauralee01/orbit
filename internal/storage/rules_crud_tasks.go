package storage

// Phase 4 (finish) — persistence for the `rules` table (migrations/001_init.sql).
// Implement the functions below in a NEW file such as rules.go (or extend this package).
// Keep SQL parameterized ($1, $2, …); check rows.Err() after every rows.Next loop.
//
// # Schema reminder
//
//   rules ( id, ruleset_id FK, field, operator, value )
//
// # A) Define a storage struct for a loaded row
//
// Example name: StoredRule — fields: ID (int64), RulesetID (int64), Field, Operator, Value (strings).
// Add `json:"..."` tags if you return this type from HTTP directly, or map to a handler DTO later.
//
// # B) InsertRule
//
// Signature idea:
//   InsertRule(ctx, db, rulesetID int64, field, operator, value string) (int64, error)
//
// - INSERT INTO rules (ruleset_id, field, operator, value) VALUES ($1,$2,$3,$4) RETURNING id
// - Return the new id (Scan into int64).
// - Optional: verify the ruleset exists first (SELECT from rulesets WHERE id = $1) and return a
//   clear error (or rely on FK violation from Postgres — less friendly error messages).
//
// # C) ListRulesByRulesetID
//
// Signature idea:
//   ListRulesByRulesetID(ctx, db, rulesetID int64) ([]StoredRule, error)
//
// - SELECT id, ruleset_id, field, operator, value FROM rules WHERE ruleset_id = $1
//   (ORDER BY id ASC if you want stable evaluation order).
// - Loop rows.Next, Scan, append; then rows.Err().
//
// # D) Optional later
//
// DeleteRule, UpdateRule — skip until you need them.
//
// # Tests (optional for this step)
//
// You can add storage tests with sqlmock or a real test DB; not required to move on to HTTP.
