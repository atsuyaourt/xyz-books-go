// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: publisher.sql

package db

import (
	"context"
	"database/sql"
)

const countPublishers = `-- name: CountPublishers :one
SELECT count(*) FROM publishers
`

func (q *Queries) CountPublishers(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, countPublishers)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createPublisher = `-- name: CreatePublisher :one
INSERT INTO publishers (
  publisher_name
) VALUES (
  ?1
) RETURNING publisher_id, publisher_name
`

func (q *Queries) CreatePublisher(ctx context.Context, publisherName string) (Publisher, error) {
	row := q.db.QueryRowContext(ctx, createPublisher, publisherName)
	var i Publisher
	err := row.Scan(&i.PublisherID, &i.PublisherName)
	return i, err
}

const deletePublisher = `-- name: DeletePublisher :exec
DELETE FROM publishers WHERE publisher_id = ?1
`

func (q *Queries) DeletePublisher(ctx context.Context, publisherID int64) error {
	_, err := q.db.ExecContext(ctx, deletePublisher, publisherID)
	return err
}

const getPublisher = `-- name: GetPublisher :one
SELECT publisher_id, publisher_name FROM publishers
WHERE publisher_id = ?1 LIMIT 1
`

func (q *Queries) GetPublisher(ctx context.Context, publisherID int64) (Publisher, error) {
	row := q.db.QueryRowContext(ctx, getPublisher, publisherID)
	var i Publisher
	err := row.Scan(&i.PublisherID, &i.PublisherName)
	return i, err
}

const getPublisherByName = `-- name: GetPublisherByName :one
SELECT publisher_id, publisher_name FROM publishers
WHERE publisher_name = ?1 LIMIT 1
`

func (q *Queries) GetPublisherByName(ctx context.Context, publisherName string) (Publisher, error) {
	row := q.db.QueryRowContext(ctx, getPublisherByName, publisherName)
	var i Publisher
	err := row.Scan(&i.PublisherID, &i.PublisherName)
	return i, err
}

const listPublishers = `-- name: ListPublishers :many
SELECT publisher_id, publisher_name FROM publishers
ORDER BY publisher_id
LIMIT ?1
OFFSET ?2
`

type ListPublishersParams struct {
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

func (q *Queries) ListPublishers(ctx context.Context, arg ListPublishersParams) ([]Publisher, error) {
	rows, err := q.db.QueryContext(ctx, listPublishers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Publisher{}
	for rows.Next() {
		var i Publisher
		if err := rows.Scan(&i.PublisherID, &i.PublisherName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePublisher = `-- name: UpdatePublisher :one
UPDATE publishers
SET
  publisher_name = COALESCE(?1, publisher_name)
WHERE
  publisher_id = ?2
RETURNING publisher_id, publisher_name
`

type UpdatePublisherParams struct {
	PublisherName sql.NullString `json:"publisher_name"`
	PublisherID   int64          `json:"publisher_id"`
}

func (q *Queries) UpdatePublisher(ctx context.Context, arg UpdatePublisherParams) (Publisher, error) {
	row := q.db.QueryRowContext(ctx, updatePublisher, arg.PublisherName, arg.PublisherID)
	var i Publisher
	err := row.Scan(&i.PublisherID, &i.PublisherName)
	return i, err
}
