CREATE TABLE ruleset_facts (
    ruleset_id BIGINT PRIMARY KEY REFERENCES rulesets(id) ON DELETE CASCADE,
    facts JSONB NOT NULL
);
