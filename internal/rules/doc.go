// Package rules is where you implement **phase 3** of the roadmap in docs/PLAN.md:
// structs for facts and rules, evaluation in memory (no HTTP, no DB), and tests.
//
// Do not import database or net/http here—keep this package “pure” so tests stay fast
// and logic is easy to reason about. You can wire HTTP to this package in a later phase.
//
// # Suggested order of work
//
// 1) **Facts** — Decide what a “fact” is for Orbit v1. Examples: a map of string keys to
// values (string/number/bool), or a small struct with named fields (user_id, country, …).
// Put the type(s) in a file such as facts.go. Keep v1 embarrassingly small.
//
// 2) **Rules** — Decide what a “rule” is before persistence: e.g. a field name, an
// operator (equals, greater than), and a value to compare against. Or a simple slice of
// conditions. Put types in rules.go (or the same file if you prefer at first).
//
// 3) **Evaluate** — Write a function with a clear signature, something like:
//   func Evaluate(facts Fact, rules []Rule) Result
// where Result describes what matched (e.g. matched rule IDs, or a boolean + message).
// Implement only in-memory logic: loops, comparisons, no I/O.
//
// 4) **Tests** — Add evaluate_test.go in this package. Use **table-driven tests**:
// a slice of struct { name string; facts …; rules …; want … } and range over it,
// calling t.Run(name, func(t *testing.T) { … }). Run: go test ./internal/rules -v
//
// 5) **Optional (still phase 3)** — Only after tests pass, you *may* add a tiny HTTP
// handler that accepts JSON facts + rules in the request body and returns the evaluation
// result as JSON—or wait until you are comfortable and add that in a follow-up commit.
//
// # Files you will likely create
//
//   internal/rules/facts.go      — fact type(s)
//   internal/rules/rules.go      — rule type(s) (name can vary)
//   internal/rules/evaluate.go    — Evaluate and helpers
//   internal/rules/evaluate_test.go — table-driven tests
//
// # How you know you are done with phase 3
//
// - go test ./internal/rules passes.
// - You can explain your Fact and Rule structs and one example evaluation path in plain English.
package rules

