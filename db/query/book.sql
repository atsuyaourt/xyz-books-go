-- name: CreateBook :one
INSERT INTO books (
  title,
  isbn13,
  isbn10,
  price,
  publication_year,
  image_url,
  edition,
  publisher_id
) VALUES (
  ?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8
) RETURNING *;

-- name: GetBookByISBN :one
SELECT
	sqlc.embed(b),
	GROUP_CONCAT(a.first_name || CASE WHEN a.middle_name IS NOT NULL THEN
			' ' || a.middle_name || ' '
		ELSE
			' '
		END || a.last_name) AS authors,
	p.publisher_name AS publisher_name
FROM
	books AS b
	JOIN author_book AS ab ON b.book_id = ab.book_id
	JOIN authors AS a ON ab.author_id = a.author_id
	JOIN publishers AS p ON b.publisher_id = p.publisher_id
WHERE
	b.isbn13 = ?1 OR b.isbn10 = ?2
GROUP BY
	b.title,
	p.publisher_name;

-- name: ListBooks :many
SELECT
  sqlc.embed(b),
  GROUP_CONCAT(a.first_name || CASE WHEN a.middle_name IS NOT NULL THEN
			' ' || a.middle_name || ' '
		ELSE
			' '
		END || a.last_name) AS authors,
  p.publisher_name AS publisher_name
FROM
  books b
JOIN author_book ab ON b.book_id = ab.book_id
JOIN authors a ON ab.author_id = a.author_id
JOIN publishers p ON b.publisher_id = p.publisher_id
WHERE
  (b.title LIKE '%' || sqlc.narg(title) || '%' OR sqlc.narg(title) IS NULL)
  AND (b.price >= sqlc.narg(min_price) OR sqlc.narg(min_price) IS NULL)
  AND (b.price <= sqlc.narg(max_price) OR sqlc.narg(max_price) IS NULL)
  AND (b.publication_year >= sqlc.narg(min_publication_year) OR sqlc.narg(min_publication_year) IS NULL)
  AND (b.publication_year <= sqlc.narg(max_publication_year) OR sqlc.narg(max_publication_year) IS NULL)
  AND (a.first_name || ' ' || a.middle_name || ' ' || a.last_name LIKE '%' || sqlc.narg(author) || '%' OR sqlc.narg(author) IS NULL)
  AND (p.publisher_name LIKE '%' || sqlc.narg(publisher) || '%' OR sqlc.narg(publisher) IS NULL)
GROUP BY
	b.title,
	p.publisher_name
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: UpdateBookByISBN :one
UPDATE books
SET
  title = COALESCE(sqlc.narg(title), title),
  isbn13 = COALESCE(sqlc.narg(new_isbn13), isbn13),
  isbn10 = COALESCE(sqlc.narg(new_isbn10), isbn10),
  price = COALESCE(sqlc.narg(price), price),
  publication_year = COALESCE(sqlc.narg(publication_year), publication_year),
  image_url = COALESCE(sqlc.narg(image_url), image_url)
WHERE
  isbn13 = @isbn13 OR isbn10 = @isbn10
RETURNING *;

-- name: DeleteBookByISBN :exec
DELETE FROM books 
WHERE 
  isbn13 = ?1
  OR isbn10 = ?2;

-- name: CountBooks :one
SELECT
  COUNT(DISTINCT b.book_id)
FROM
  books b
JOIN author_book ab ON b.book_id = ab.book_id
JOIN authors a ON ab.author_id = a.author_id
JOIN publishers p ON b.publisher_id = p.publisher_id
WHERE
  (b.title LIKE '%' || sqlc.narg(title) || '%' OR sqlc.narg(title) IS NULL)
  AND (b.price >= sqlc.narg(min_price) OR sqlc.narg(min_price) IS NULL)
  AND (b.price <= sqlc.narg(max_price) OR sqlc.narg(max_price) IS NULL)
  AND (b.publication_year >= sqlc.narg(min_publication_year) OR sqlc.narg(min_publication_year) IS NULL)
  AND (b.publication_year <= sqlc.narg(max_publication_year) OR sqlc.narg(max_publication_year) IS NULL)
  AND (a.first_name || ' ' || a.middle_name || ' ' || a.last_name LIKE '%' || sqlc.narg(author) || '%' OR sqlc.narg(author) IS NULL)
  AND (p.publisher_name LIKE '%' || sqlc.narg(publisher) || '%' OR sqlc.narg(publisher) IS NULL);
