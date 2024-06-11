-- name: CreateAuthor :one
INSERT INTO authors (
  first_name,
  last_name,
  middle_name
) VALUES (
  ?1, ?2, ?3
) RETURNING *;

-- name: GetAuthor :one
SELECT * FROM authors
WHERE author_id = ?1 LIMIT 1;

-- name: GetAuthorByName :one
SELECT * FROM authors
WHERE
  first_name = @first_name AND
  last_name = @last_name AND
  middle_name = COALESCE(@middle_name, middle_name)
LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY author_id
LIMIT ?1
OFFSET ?2;

-- name: UpdateAuthor :one
UPDATE authors
SET
  first_name = COALESCE(sqlc.narg(first_name), first_name),
  last_name = COALESCE(sqlc.narg(last_name), last_name),
  middle_name = COALESCE(sqlc.narg(middle_name), middle_name)
WHERE
  author_id = sqlc.arg(author_id)
RETURNING *;

-- name: DeleteAuthor :exec
DELETE FROM authors WHERE author_id = ?1;

-- name: CountAuthors :one
SELECT count(*) FROM authors;