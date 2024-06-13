package services

import (
	"github.com/atsuyaourt/xyz-books/internal/models"
	"github.com/atsuyaourt/xyz-books/internal/util"
	"golang.org/x/net/context"
)

type Service interface {
	CreateBook(ctx context.Context, req CreateBookReq) (*models.Book, error)
	GetBook(ctx context.Context, isbn13 string) (*models.Book, error)
	ListBooks(ctx context.Context, req ListBooksReq) (*util.PaginatedList[models.Book], error)
	UpdateBook(ctx context.Context, oldISBN13 string, req UpdateBookReq) (*models.Book, error)
	DeleteBook(ctx context.Context, isbn13 string) error

	CreateAuthor(ctx context.Context, req CreateAuthorReq) (*models.Author, error)
	GetAuthor(ctx context.Context, id int64) (*models.Author, error)
	ListAuthors(ctx context.Context, req ListAuthorsReq) (*util.PaginatedList[models.Author], error)
	UpdateAuthor(ctx context.Context, oldID int64, req UpdateAuthorReq) (*models.Author, error)
	DeleteAuthor(ctx context.Context, id int64) error

	CreatePublisher(ctx context.Context, req CreatePublisherReq) (*models.Publisher, error)
	GetPublisher(ctx context.Context, id int64) (*models.Publisher, error)
	ListPublishers(ctx context.Context, req ListPublishersReq) (*util.PaginatedList[models.Publisher], error)
	UpdatePublisher(ctx context.Context, oldID int64, req UpdatePublisherReq) (*models.Publisher, error)
	DeletePublisher(ctx context.Context, id int64) error
}
