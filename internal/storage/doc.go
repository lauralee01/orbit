// Package storage is where you implement **phase 4** of docs/PLAN.md: persist rules (and
// related metadata) in a relational database, using Go’s database/sql API and SQL.
//
// internal/rules stays pure (no imports from here into rules). This package only knows
// how to load and save rows; evaluation still happens in rules.Evaluate.
//
// # 0) Prereqs (your machine)
//
// - Install PostgreSQL locally, or run it with Docker (official image is fine).
// - Create a database for Orbit (e.g. orbit_dev) and a user/password you can connect with.
// - Decide a connection string; common pattern is env var DATABASE_URL, e.g.:
//     postgres://USER:PASSWORD@localhost:5432/orbit_dev?sslmode=disable
//
// # 1) Add a Postgres driver (dependency)
//
// - Pick one driver that works with database/sql:
//   - jackc/pgx (stdlib adapter) is widely used today, or
//   - lib/pq (older, still common).
// - Run: go get <module> and add a blank import in your OpenDB code path, e.g.:
//     import _ "github.com/jackc/pgx/v5/stdlib"
//   The underscore import registers the driver with database/sql by side effect.
//
// # 2) Open the pool and ping once
//
// - In a new file here (e.g. db.go), write a function Open(ctx, databaseURL) (*sql.DB, error)
//   that calls sql.Open("pgx", url) or sql.Open("postgres", url) matching your driver.
// - Immediately sql.DB.PingContext(ctx) so misconfiguration fails fast at startup.
// - sql.DB is a **connection pool**, not a single connection—reuse one *sql.DB for the app.
//
// # 3) Design a minimal schema (SQL)
//
// - Plan for at least: storing many rules, and optionally grouping them (a “ruleset” or
//   “policy” with a name or ID). Example shapes (you choose the names):
//   - rulesets: id, name, created_at
//   - rules: id, ruleset_id (FK), field, operator, value, maybe sort order
// - Write plain .sql files or a short README of CREATE TABLE statements you applied.
// - Apply them manually with psql first, or add a tiny “migrate” step (optional tools:
//   golang-migrate, goose—only if you want; raw SQL in main once is OK for learning).
//
// # 4) Implement CRUD (or start with Create + List)
//
// - New file(s) in this package, e.g. rulesets.go / rules.go, with functions that take
//   (ctx context.Context, db *sql.DB, ...) and run parameterized queries.
// - Always use **placeholders** ($1, $2, …) in your SQL—never string-concatenate user input
//   into queries (SQL injection).
// - Return data as your own small structs (e.g. StoredRule with ID, Field, Operator, Value)
//   or reuse/adapt types from internal/rules if it stays clean (often a separate storage DTO
//   is clearer than coupling DB tags to rule engine types).
//
// # 5) Learn context in the DB API
//
// - Prefer QueryContext, ExecContext, PingContext over the non-Context variants.
// - Handlers will pass r.Context() (or a derived timeout context) into storage calls so
//   slow DB work can be cancelled when the client disconnects.
//
// # 6) Wire it from cmd/orbit/main.go (no HTTP required at first)
//
// - Read DATABASE_URL from os.Getenv in main.
// - Open *sql.DB, defer db.Close().
// - Call a storage function to insert a test ruleset/rule, then list them back, log.Println—
//   proves the stack works before you add REST.
//
// # 7) Optional same phase: HTTP for rules
//
// - In internal/handlers, add POST/GET routes that decode JSON, call storage, return JSON
//   (same pattern as Echo). Register routes in main.
//
// # 8) Done when
//
// - You can restart the app and still see rules in Postgres (data survives restart).
// - go test ./... still passes (add storage tests later with a test DB or sqlite—optional;
//   integration tests often use docker-compose or env-gated tests).
//
// Work in the order above; commit after each step that compiles and runs.
package storage
