-- name: CreateAuthorBookRel :exec
INSERT INTO author_book (
  author_id,
  book_id
) VALUES (
  ?1, ?2
) RETURNING *;

-- name: ListAuthorsWithBookID :many
SELECT sqlc.embed(a)
FROM
  author_book ab
  JOIN authors a ON ab.author_id = a.author_id
WHERE ab.book_id = ?1;
