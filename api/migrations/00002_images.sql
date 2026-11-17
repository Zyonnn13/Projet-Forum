-- +goose Up
CREATE TABLE images (
    nom VARCHAR(100) PRIMARY KEY,
    content_type VARCHAR(100) NOT NULL,
    data BYTEA NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE images;