// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"
)

type Querier interface {
	CountAuthors(ctx context.Context) (int64, error)
	CountBooks(ctx context.Context, arg CountBooksParams) (int64, error)
	CountPublishers(ctx context.Context) (int64, error)
	CreateAuthor(ctx context.Context, arg CreateAuthorParams) (Author, error)
	CreateAuthorBookRel(ctx context.Context, arg CreateAuthorBookRelParams) error
	CreateBook(ctx context.Context, arg CreateBookParams) (Book, error)
	CreatePublisher(ctx context.Context, publisherName string) (Publisher, error)
	DeleteAuthor(ctx context.Context, authorID int64) error
	DeleteBookByISBN(ctx context.Context, arg DeleteBookByISBNParams) error
	DeletePublisher(ctx context.Context, publisherID int64) error
	GetAuthor(ctx context.Context, authorID int64) (Author, error)
	GetAuthorByName(ctx context.Context, arg GetAuthorByNameParams) (Author, error)
	GetBookByISBN(ctx context.Context, arg GetBookByISBNParams) (GetBookByISBNRow, error)
	GetPublisher(ctx context.Context, publisherID int64) (Publisher, error)
	GetPublisherByName(ctx context.Context, publisherName string) (Publisher, error)
	ListAuthors(ctx context.Context, arg ListAuthorsParams) ([]Author, error)
	ListAuthorsWithBookID(ctx context.Context, bookID int64) ([]ListAuthorsWithBookIDRow, error)
	ListBooks(ctx context.Context, arg ListBooksParams) ([]ListBooksRow, error)
	ListPublishers(ctx context.Context, arg ListPublishersParams) ([]Publisher, error)
	UpdateAuthor(ctx context.Context, arg UpdateAuthorParams) (Author, error)
	UpdateBookByISBN(ctx context.Context, arg UpdateBookByISBNParams) (Book, error)
	UpdatePublisher(ctx context.Context, arg UpdatePublisherParams) (Publisher, error)
}

var _ Querier = (*Queries)(nil)
