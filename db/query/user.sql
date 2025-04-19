-- db/query/user.sql

-- name: CreateUser :one
INSERT INTO users (
  username,
  email,
  hashed_password,
  full_name -- Add full_name
) VALUES (
  $1, $2, $3, $4 -- Add $4 for full_name
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- We can add ListUsers, UpdateUser, DeleteUser etc. later as needed.