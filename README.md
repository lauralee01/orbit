# Orbit

**Orbit** is a small **rules engine** service in **Go**: you store **rulesets** and **rules** in **PostgreSQL**, send **facts** as JSON, and **evaluate** whether those facts satisfy the persisted rules.

## Why this project

- Learn Go alongside real backend patterns: HTTP + JSON, `database/sql`, and tests.
- Keep scope small: no custom rule language in v1—just field / operator / value conditions.

## Status

Core flows are implemented: persistence, REST-style JSON APIs for rulesets and rules, and **POST /api/evaluate** to run evaluation against stored rules.

**Rulesets** support optional **`webhook_url`** and **schedule** metadata (`schedule_cron`, `schedule_tz`, `schedule_enabled`). Create accepts these fields in JSON. If **`schedule_enabled`** is true, **`schedule_cron`** (standard 5-field cron) and **`schedule_tz`** (IANA zone, e.g. `America/New_York`) are required and validated; **`POST /api/evaluate`** does not gate on the schedule yet (schedule is metadata for future automation).

## Requirements

- [Go](https://go.dev/dl/) **1.26+** (see `go.mod`).
- **PostgreSQL** reachable via a connection URI.
- **Environment:** the process reads **`DATABASE_URL`** (required) and **`PORT`** (optional; defaults to `8080` when unset).

## Clone and run locally

```bash
git clone https://github.com/lauralee01/orbit.git
cd orbit
```

Copy **`.env.example`** to **`.env`**, set **`DATABASE_URL`** to your local Postgres URL, then:

```bash
go run ./cmd/orbit
```

`godotenv` loads **`.env`** when the file exists; if it does not (e.g. in Docker), only OS environment variables are used.

Optional port override:

```bash
PORT=3000 go run ./cmd/orbit
```

Build a binary:

```bash
go build -o orbit ./cmd/orbit
./orbit
```

## Database migrations

Create an empty database, then apply SQL in **order** (Neon SQL Editor, `psql`, or any Postgres client):

1. `migrations/001_init.sql` — tables `rulesets` and `rules`
2. `migrations/002_webhook.sql` — `webhook_url` on `rulesets`
3. `migrations/003_schedule.sql` — `schedule_cron`, `schedule_tz`, `schedule_enabled` on `rulesets`

Example with `psql`:

```bash
psql "$DATABASE_URL" -f migrations/001_init.sql
psql "$DATABASE_URL" -f migrations/002_webhook.sql
psql "$DATABASE_URL" -f migrations/003_schedule.sql
```

## Environment variables

| Variable | Required | Notes |
|----------|----------|--------|
| `DATABASE_URL` | Yes | Postgres URI, e.g. `postgres://user:pass@host:5432/dbname?sslmode=disable` (local) or your host’s URI with `sslmode=require` when TLS is required. |
| `PORT` | No | Listen port; Render and similar platforms set this automatically. |

See **`.env.example`** for a template. Extra keys there are for documentation only until wired in code. Do **not** commit secrets; **`.env`** is gitignored.

## Docker

From the repository root:

```bash
docker build -t orbit .
docker run --rm -p 8080:8080 -e DATABASE_URL="postgres://..." orbit
```

Then open `http://localhost:8080/health`.

## Deployment (example: Render + Neon)

This repo includes a **`Dockerfile`** for container builds. A typical setup:

1. **Neon:** create a project, run the migrations above on that database, copy the connection string (often includes `sslmode=require`).
2. **Render:** new **Web Service**, connect the Git repo, runtime **Docker**, root directory empty if the `Dockerfile` is at the repo root.
3. **Environment** on Render: set **`DATABASE_URL`** to the Neon URI. Do not set **`PORT`** unless you have a reason—Render injects it.
4. **Health check:** HTTP GET path **`/health`**.
5. Deploy; smoke-test **`GET https://<your-service>.onrender.com/health`** and the API examples below.

On first boot without a `.env` file you may see a `godotenv: open .env: no such file` log line; that is expected and harmless in production.

## API (overview)

| Method | Path | Purpose |
|--------|------|--------|
| GET | `/health` | Liveness check |
| GET / POST | `/api/rulesets` | List / create rulesets (create accepts optional `webhook_url`, schedule fields; see below) |
| GET / POST | `/api/rules` | List rules (`GET ?ruleset_id=`) / create rule (JSON body includes `ruleset_id`) |
| POST | `/api/evaluate` | Body: `ruleset_id` + `facts` → `{ "ok": true \| false, ... }` (errors may include `detail`) |

### Ruleset schedule fields (create body)

Optional JSON keys:

| Field | Type | Notes |
|-------|------|--------|
| `webhook_url` | string | Omitted or empty if unused. |
| `schedule_cron` | string | Standard 5-field cron (validated when schedule is enabled). |
| `schedule_tz` | string | IANA timezone (e.g. `UTC`, `Europe/London`); required when `schedule_enabled` is true. |
| `schedule_enabled` | bool | Default false. If true, `schedule_cron` and `schedule_tz` must be non-empty and valid. |

If `schedule_enabled` is false, empty `schedule_cron` is stored as disabled; empty `schedule_tz` defaults to `UTC` in storage.

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

Optional: create a ruleset with a schedule (validated cron + timezone):

```bash
curl -s -X POST http://localhost:8080/api/rulesets -H 'Content-Type: application/json' \
  -d '{"name":"timed","schedule_cron":"0 9 * * *","schedule_tz":"America/New_York","schedule_enabled":true}'
```

## Tests

```bash
go test ./...
```

## License

To be decided.
