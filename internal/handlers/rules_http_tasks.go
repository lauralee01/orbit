package handlers

// Phase 4 (finish) — HTTP CRUD for **rules** (child of a ruleset).
// Implement handlers in a new file (e.g. rules.go). Reuse writeJSON, errorResponse, maxBodyBytes.
//
// # Design choice: where does ruleset_id live?
//
// Pick ONE style (both are valid):
//
//   **A) REST-style path** — e.g. POST /api/rulesets/{id}/rules and GET /api/rulesets/{id}/rules
//       You must parse the ruleset id from r.URL.Path or use Go 1.22+ ServeMux path patterns.
//       Read: net/http ServeMux patterns and PathValue (Go 1.22+) if available.
//
//   **B) Simpler body/query** — e.g. POST /api/rules with JSON
//       {"ruleset_id":1,"field":"age","operator":"equals","value":"30"}
//       and GET /api/rules?ruleset_id=1
//       Easiest to implement first; refactor to nested paths later if you want.
//
// # POST — create a rule
//
// - Decode JSON into a struct (field, operator, value, plus ruleset_id if using style B).
// - Validate: non-empty field/operator/value; ruleset_id > 0.
// - Call storage.InsertRule(r.Context(), db, …).
// - Return 201 + JSON with the new rule id (and echo fields you like).
// - 404 if ruleset does not exist (if you add an explicit check).
// - MaxBytesReader on the body (same as CreateRuleset).
//
// # GET — list rules for a ruleset
//
// - Read ruleset id from path or query (see design choice above).
// - Call storage.ListRulesByRulesetID(r.Context(), db, rulesetID).
// - Return JSON array (empty slice is fine — not an error).
//
// # Wire in cmd/orbit/main.go
//
// Register your routes with mux.HandleFunc, passing db via closures like ListRulesets(db).
// Use method-specific patterns if you use stdlib mux: "POST /api/rules", etc.
