# Orbit

**Orbit** is a rules engine—built with **Go**.

The goal is to accept structured input (“facts”), apply configurable rules, and return clear outcomes. Scope will stay small at first and expand as the foundations feel solid.

## Why this project

- Learn Go from scratch alongside real service patterns (HTTP, JSON, persistence, tests).
- Explore how rule evaluation, storage, and APIs fit together without frontend complexity dominating the story.

## Status

Early setup. No API or runtime guarantees yet.

## Requirements

- [Go](https://go.dev/dl/) (install the latest stable release for your OS).

## Getting started

From the repo root:

```bash
go run ./cmd/orbit
```

Then in another terminal:

```bash
curl -s http://localhost:8080/health
```

Optional: set a port with `PORT=3000 go run ./cmd/orbit`.

Build a binary:

```bash
go build -o orbit ./cmd/orbit
./orbit
```

## License

To be decided.
