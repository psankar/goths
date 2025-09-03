BEGIN;

CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE posts (
    post_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(user_id),
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    published_at TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO users (email, password) VALUES 
    ('user1@example.com', 'password'),
    ('user2@example.com', 'password');

COMMIT;
