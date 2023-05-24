CREATE TYPE roles AS ENUM ('hr', 'admin');

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(64),
    password_salt VARCHAR(32),
    email VARCHAR(320) UNIQUE,
    is_verified BOOLEAN,
    roles roles NOT NULL DEFAULT 'hr',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);