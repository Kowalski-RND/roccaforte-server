CREATE TABLE users (
    id UUID PRIMARY KEY NOT NULL,
    username TEXT NOT NULL,
    fullname TEXT NOT NULL,
    password TEXT NOT NULL,
    public_key TEXT NOT NULL
);
CREATE UNIQUE INDEX user_username_key ON users (username);

CREATE TABLE secrets (
    id UUID PRIMARY KEY NOT NULL,
    author UUID NOT NULL,
    cipher_text TEXT NOT NULL,
    iv TEXT NOT NULL,
    CONSTRAINT secret_author_fkey FOREIGN KEY (author) REFERENCES users (id)
);

CREATE TABLE keys (
    id UUID PRIMARY KEY NOT NULL,
    owner UUID NOT NULL,
    secret UUID NOT NULL,
    key TEXT NOT NULL,
    CONSTRAINT key_owner_fkey FOREIGN KEY (owner) REFERENCES users (id),
    CONSTRAINT keys_secret_fkey FOREIGN KEY (secret) REFERENCES secrets (id) ON DELETE CASCADE
);