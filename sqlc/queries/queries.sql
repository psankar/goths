-- name: Login :one
-- password should be hashed and stored and compared in real applications
SELECT * FROM users WHERE email = $1 AND password = $2;
