-- +goose Up
CREATE TABLE utilisateurs (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    banned BOOLEAN NOT NULL DEFAULT FALSE,
    avatar VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    nom VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE fils_de_discussion (
    id SERIAL PRIMARY KEY,
    titre VARCHAR(255) NOT NULL,
    etat VARCHAR(50) NOT NULL DEFAULT 'ouvert',
    auteur_id INT NOT NULL REFERENCES utilisateurs(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE fils_categories (
    fil_id INT NOT NULL REFERENCES fils_de_discussion(id) ON DELETE CASCADE,
    categorie_id INT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (fil_id, categorie_id)
);

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    contenu TEXT NOT NULL,
    fil_id INT NOT NULL REFERENCES fils_de_discussion(id) ON DELETE CASCADE,
    auteur_id INT NOT NULL REFERENCES utilisateurs(id) ON DELETE CASCADE,
    score_popularite INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE reactions (
    utilisateur_id INT NOT NULL REFERENCES utilisateurs(id) ON DELETE CASCADE,
    message_id INT NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    PRIMARY KEY (utilisateur_id, message_id)
);

-- +goose StatementBegin
CREATE FUNCTION set_updated_at() RETURNS trigger AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER fils_updated_at BEFORE UPDATE ON fils_de_discussion
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER messages_updated_at BEFORE UPDATE ON messages
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS reactions, messages, fils_categories, fils_de_discussion, categories, utilisateurs;
DROP FUNCTION IF EXISTS set_updated_at;