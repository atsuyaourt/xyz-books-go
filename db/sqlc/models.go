// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package db

import (
	"database/sql"
)

type Author struct {
	AuthorID   int64  `json:"author_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
}

type AuthorBook struct {
	AuthorID int64 `json:"author_id"`
	BookID   int64 `json:"book_id"`
}

type Book struct {
	BookID          int64          `json:"book_id"`
	Title           string         `json:"title"`
	Isbn13          sql.NullString `json:"isbn13"`
	Isbn10          sql.NullString `json:"isbn10"`
	Price           float64        `json:"price"`
	PublicationYear int64          `json:"publication_year"`
	ImageUrl        sql.NullString `json:"image_url"`
	Edition         sql.NullString `json:"edition"`
	PublisherID     int64          `json:"publisher_id"`
}

type Publisher struct {
	PublisherID   int64  `json:"publisher_id"`
	PublisherName string `json:"publisher_name"`
}
