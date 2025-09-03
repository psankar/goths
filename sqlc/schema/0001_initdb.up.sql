BEGIN;

CREATE TABLE authors (
    author_id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE
);

CREATE TABLE posts (
    post_id SERIAL PRIMARY KEY,
    author_id INT NOT NULL REFERENCES authors(author_id),
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    published_at TIMESTAMPTZ DEFAULT NOW()
);

COMMIT;
