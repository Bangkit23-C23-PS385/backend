CREATE TABLE IF EXISTS verifications (
    id SERIAL PRIMARY KEY,
    email VARCHAR(320) UNIQUE,
    token VARCHAR(64),
    attempt_left INTEGER DEFAULT 3,
    CONSTRAINT fk_email_users_verifications FOREIGN KEY(email) REFERENCES users(email),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);