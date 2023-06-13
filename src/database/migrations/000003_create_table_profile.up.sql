CREATE TABLE profiles (
    userId VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255),
    gender VARCHAR(255),
    dateOfBirth DATE,
    height INT,
    weight INT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);