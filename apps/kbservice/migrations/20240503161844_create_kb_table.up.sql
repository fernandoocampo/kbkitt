-- Add up migration script here
CREATE TABLE IF NOT EXISTS kbs (
    KB_ID VARCHAR(36) PRIMARY KEY,
    KB_KEY VARCHAR(64) NOT NULL UNIQUE,
    KB_VALUE TEXT NOT NULL,
    NOTES TEXT NOT NULL,
    KIND VARCHAR(64) NOT NULL,
    TAGS TSVECTOR NOT NULL,
    CREATED_ON TIMESTAMP NOT NULL DEFAULT NOW()
);