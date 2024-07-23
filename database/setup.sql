CREATE DATABASE ecdsadb;

\c ecdsadb;

CREATE USER ecdsa_user WITH PASSWORD 'password';

GRANT ALL PRIVILEGES ON DATABASE ecdsadb TO ecdsa_user;

\c ecdsadb ecdsa_user;

CREATE TABLE statuses (
    id SERIAL PRIMARY KEY,
    status_id VARCHAR(255) NOT NULL,
    status_list BYTEA NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_status_id ON statuses (status_id);

INSERT INTO statuses (status_id, status_list, created_at, updated_at)
VALUES ('testStatusId', '\x01', NOW(), NOW());

SELECT * FROM statuses;
