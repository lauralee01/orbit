CREATE TABLE rulesets (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE rules (
    id SERIAL PRIMARY KEY,
    ruleset_id INT NOT NULL REFERENCES rulesets(id) ON DELETE CASCADE,
    field TEXT NOT NULL,
    operator TEXT NOT NULL,
    value TEXT NOT NULL
);
