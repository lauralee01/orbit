# Orbit

**Orbit** is a small **rules engine** service in **Go**: you store **rulesets** and **rules** in **PostgreSQL**, send **facts** as JSON, and **evaluate** whether those facts satisfy the persisted rules.

## Why this project

- Learn Go alongside real backend patterns: HTTP + JSON, `database/sql`, and tests.
- Keep scope small: no custom rule language in v1—just field / operator / value conditions.

## Status

Core flows are implemented: persistence, REST-style JSON APIs for rulesets and rules, and **POST /api/evaluate** to run **Evaluate** against stored rules.


## Requirements

- [Go](https://go.dev/dl/) (1.22+ recommended for method-specific `ServeMux` routes).
- **PostgreSQL** and a database created for local dev.
- **`DATABASE_URL`** — e.g. `postgres://USER:PASSWORD@localhost:5432/DBNAME?sslmode=disable`

Apply the schema in `migrations/001_init.sql` once (e.g. via `psql` or your GUI) before running the app.

Optional: a `.env` file with `DATABASE_URL=` (loaded via `godotenv` in `main`); otherwise export the variable in your shell.

## Run

```bash
go run ./cmd/orbit
```

Optional port: `PORT=3000 go run ./cmd/orbit`

Build a binary:

```bash
go build -o orbit ./cmd/orbit
./orbit
```

## API (overview)

| Method | Path | Purpose |
|--------|------|--------|
| GET | `/health` | Liveness check |
| GET / POST | `/api/rulesets` | List / create rulesets |
| GET / POST | `/api/rules` | List rules (`GET ?ruleset_id=`) / create rule (JSON body includes `ruleset_id`) |
| POST | `/api/evaluate` | Body: `ruleset_id` + `facts` object → `{ "ok": true \| false }` (errors may include `detail`) |

## Quick manual check

```bash
curl -s http://localhost:8080/health
```

Create a ruleset, add a rule, then evaluate (adjust IDs to match your DB):

```bash
curl -s -X POST http://localhost:8080/api/rulesets -H 'Content-Type: application/json' -d '{"name":"policy"}'
curl -s -X POST http://localhost:8080/api/rules -H 'Content-Type: application/json' \
  -d '{"ruleset_id":1,"field":"age","operator":"equals","value":"30"}'
curl -s -X POST http://localhost:8080/api/evaluate -H 'Content-Type: application/json' \
  -d '{"ruleset_id":1,"facts":{"age":30}}'
```

## License

To be decided.
