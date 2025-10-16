-- Redirect rules table: exact-path matches
CREATE TABLE redirect_rules_full_match (
    id SERIAL PRIMARY KEY,
    from_raw TEXT NOT NULL,
    from_canonical TEXT NOT NULL,
    target TEXT NOT NULL
);

-- Unique index to ensure single rule per path
CREATE UNIQUE INDEX full_match_unique_from_canonical_idx ON redirect_rules_full_match(from_canonical);

-- Welcome page URL table (single row expected)
CREATE TABLE welcome_page_url (
    id SERIAL PRIMARY KEY,
    target TEXT NOT NULL
);