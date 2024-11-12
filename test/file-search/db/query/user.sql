-- name: CreateUser :one
INSERT INTO users (
  email,
  username,
  password,
  password_hash,
  phone,
  fullname,
  avatar,
  state,
  role,
  created_at,
  update_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;


-- name: GetUsersAsc :many
SELECT * FROM users
WHERE (username ILIKE '%' || $1 || '%' OR fullname ILIKE '%' || $1 || '%')
  AND state = $5
ORDER BY 
  CASE WHEN $4 = 'id' THEN id END ASC,
  CASE WHEN $4 = 'username' THEN username END ASC,
  CASE WHEN $4 = 'fullname' THEN fullname END ASC
LIMIT $2 OFFSET $3;

-- name: GetUsersDesc :many
SELECT * FROM users
WHERE (username ILIKE '%' || $1 || '%' OR fullname ILIKE '%' || $1 || '%')
  AND state = $5
ORDER BY
  CASE WHEN $4 = 'id' THEN id END DESC,
  CASE WHEN $4 = 'username' THEN username END DESC,
  CASE WHEN $4 = 'fullname' THEN fullname END DESC
LIMIT $2 OFFSET $3;


-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET
  fullname = COALESCE(sqlc.narg(fullname), fullname),
  email = COALESCE(sqlc.narg(email), email),
  phone = COALESCE(sqlc.narg(phone), phone),
  password = COALESCE(sqlc.narg(password), password),
  avatar = COALESCE(sqlc.narg(avatar), avatar),
  role = COALESCE(sqlc.narg(role), role),
  password_hash = COALESCE(sqlc.narg(password_hash), password_hash)
WHERE
  username = sqlc.arg(username)
RETURNING *;

-- name: DeleteUser :one
UPDATE users
SET state = $1
WHERE id = $2
RETURNING *;
