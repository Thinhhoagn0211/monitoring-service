-- name: InsertFile :one
INSERT INTO file (
  name,
  extension,
  size,
  path,
  attributes,
  content,
  created_at,
  modified_at,
  accessed_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;