-- name: GetUsers :many
SELECT * FROM Users;

-- name: GetUser :one
SELECT * FROM Users WHERE username = ? AND password = ?;

-- name: CreateUser :one
INSERT INTO Users (
  username,
  password
) VALUES (
  ?,
  ?
)
RETURNING *;

-- name: DeleteUser :exec
DELETE From Users WHERE id = ?;

-- name: GetUserById :one
SELECT * FROM Users WHERE id = ?;

-- name: ChangePassword :one
UPDATE Users SET password = ? WHERE id = ? RETURNING *;
