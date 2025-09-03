-- name: Login :one
SELECT * FROM authors WHERE email = $1;
