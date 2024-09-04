-- name: GetUsers :many
SELECT * FROM Users;

-- name: GetUser :one
SELECT * FROM Users WHERE username = ?;

-- name: CreateUser :one
INSERT INTO Users (
  username,
  email,
  password,
  salt,
  isAdmin,
  active
) VALUES (
  ?,
  ?,
  ?,
  ?,
  ?,
  ?
)
RETURNING *;

-- name: DeleteUser :exec
DELETE From Users WHERE id = ?;

-- name: GetUserById :one
SELECT * FROM Users WHERE id = ?;

-- name: ChangePassword :one
UPDATE Users SET password = ?, salt = ? WHERE id = ? RETURNING *;

-- name: MarkActive :exec
Update Users SET active = true WHERE id = ?;
