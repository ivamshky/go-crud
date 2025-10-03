-- name: GetUserById :one
SELECT * FROM users
WHERE id = ?;

-- name: ListUsers :many
SELECT * FROM users
WHERE
    (sqlc.narg(id) IS NULL OR id = sqlc.narg(id))
AND (sqlc.narg(name) IS NULL OR name = sqlc.narg(name))
AND (sqlc.narg(email) IS NULL OR email = sqlc.narg(email))
LIMIT ?
OFFSET ?;

-- name: CreateUser :execresult
INSERT INTO users (name,email,age)
VALUES (?, ?, ?);

-- name: DeleteUser :exec
DELETE FROM users
where id = ?;