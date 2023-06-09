CREATE INDEX sessions_token_hash_idx ON sessions(token_hash);

CREATE TABLE dogs(
    id SERIAL PRIMARY KEY,
    name TEXT,
    user_id REFERENCES users(id) ON DELETE CASCADE,
    -- Does NOT work with Postgres
    INDEX dogs_user_id_idx (user_id)
);