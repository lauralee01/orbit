package actions

// Phase C.1 — Facts provider implementation tasks (comment-only).
//
// Current state:
// - Scheduler requires a FactsProvider function.
// - cmd/orbit/main.go currently uses a hardcoded placeholder fact map.
//
// Goal:
// - Replace placeholder facts with a real, explicit source.
//
// Suggested implementation path:
// 1) Define one concrete source for v1:
//    - Option A: static facts from DB table (recommended first)
//    - Option B: facts passed from external API call
//    - Option C: facts built from another service
//
// 2) Create a dedicated package for provider logic, e.g.:
//    internal/facts/provider.go
//    Keep scheduler unaware of how facts are produced.
//
// 3) Design provider contract details:
//    - Input: context + ruleset
//    - Output: rules.Facts + error
//    - "no data available" behavior:
//      - return nil, nil  => scheduler logs and skips run
//      - return error     => scheduler logs provider error
//
// 4) Add minimal validation in provider:
//    - Ensure returned map keys match expected rule fields where possible.
//    - Normalize numeric/boolean shapes if your source is loosely typed.
//
// 5) Wire provider into main:
//    - Replace inline factsProvider in cmd/orbit/main.go with package function.
//    - Keep function small and easy to mock in tests.
//
// Exit criteria:
// - Scheduled evaluations use real data, not hardcoded placeholders.
// - Behavior for missing data is documented and deterministic.
