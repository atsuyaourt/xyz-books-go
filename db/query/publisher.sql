-- name: CreatePublisher :one
INSERT INTO publishers (
  publisher_name
) VALUES (
  ?1
) RETURNING *;

-- name: GetPublisher :one
SELECT * FROM publishers
WHERE publisher_id = ?1 LIMIT 1;

-- name: GetPublisherByName :one
SELECT * FROM publishers
WHERE publisher_name = ?1 LIMIT 1;

-- name: ListPublishers :many
SELECT * FROM publishers
ORDER BY publisher_id
LIMIT ?1
OFFSET ?2;

-- name: UpdatePublisher :one
UPDATE publishers
SET
  publisher_name = COALESCE(sqlc.narg(publisher_name), publisher_name)
WHERE
  publisher_id = sqlc.arg(publisher_id)
RETURNING *;

-- name: DeletePublisher :exec
DELETE FROM publishers WHERE publisher_id = ?1;

-- name: CountPublishers :one
SELECT count(*) FROM publishers;