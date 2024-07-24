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
