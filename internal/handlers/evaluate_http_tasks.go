package handlers

// Phase 4 (finish) — HTTP **evaluation**: load persisted rules for a ruleset, run internal/rules.Evaluate.
// Implement in a new file (e.g. evaluate.go). This closes the “facts in → result out” loop from docs/PLAN.md.
//
// # Request shape (you define the exact JSON)
//
// Client sends **facts** (same idea as rules.Facts: map[string]any or a JSON object) plus a way
// to pick which ruleset to use, e.g.:
//   { "ruleset_id": 1, "facts": { "age": 30, "name": "Ada" } }
//
// # Handler steps
//
// 1. POST only; MaxBytesReader + json.Decoder; validate ruleset_id and facts present.
// 2. storage.ListRulesByRulesetID(ctx, db, rulesetID) → []storage.StoredRule (you defined this).
// 3. Convert []StoredRule to rules.Rules:
//    loop and append rules.Rule{ Field: s.Field, Operator: s.Operator, Value: s.Value }.
//    (Order: use same order as SQL ORDER BY if you added it.)
// 4. facts map: if your JSON uses "facts": { ... }, range and build rules.Facts (map[string]any).
// 5. Call rules.Evaluate(facts, ruleSlice) — handle (bool, error); map errors to 400/500 as appropriate.
// 6. Respond JSON, e.g. { "ok": true } or { "ok": false, "error": "..." } — your choice; document it.
//
// # Edge cases
//
// - Ruleset exists but has zero rules: Evaluate may return (true, nil) with your current AND semantics
//   over an empty slice — confirm that matches what you want, or define behavior explicitly.
// - Unknown operator in Evaluate vs DB: same issue as before; document supported operators.
//
// # Wire in main
//
// e.g. mux.HandleFunc("POST /api/evaluate", handlers.Evaluate(db)) — path name is yours.
//
// # Optional polish
//
// - GET ruleset by id before evaluate to ensure it exists (404).
// - context.WithTimeout for DB calls if you want a hard ceiling (ties into Phase 5).
