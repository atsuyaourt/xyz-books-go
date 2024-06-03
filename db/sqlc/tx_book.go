package db

import (
	"context"
	"errors"

	"github.com/atsuyaourt/xyz-books/internal/util"
)

type CreateBookTxParams struct {
	Book      CreateBookParams
	Authors   []util.Name
	Publisher string
}

func (store *SQLStore) CreateBookTx(ctx context.Context, arg CreateBookTxParams) (book Book, err error) {
	err = store.execTx(ctx, func(q *Queries) error {
		authors := make([]Author, len(arg.Authors))
		for i, authorInfo := range arg.Authors {
			authors[i], err = store.GetAuthorByName(ctx, GetAuthorByNameParams(authorInfo))
			if err != nil {
				if !errors.Is(err, ErrRecordNotFound) {
					return err
				}
				authors[i], err = store.CreateAuthor(ctx, CreateAuthorParams(authorInfo))
				if err != nil {
					return err
				}
			}
		}

		publisher, err := store.GetPublisherByName(ctx, arg.Publisher)
		if err != nil {
			if !errors.Is(err, ErrRecordNotFound) {
				return err
			}
			publisher, err = store.CreatePublisher(ctx, arg.Publisher)
			if err != nil {
				return err
			}
		}

		book, err = store.CreateBook(ctx, CreateBookParams{
			Title:           arg.Book.Title,
			Isbn13:          arg.Book.Isbn13,
			Isbn10:          arg.Book.Isbn10,
			Price:           arg.Book.Price,
			PublicationYear: arg.Book.PublicationYear,
			ImageUrl:        arg.Book.ImageUrl,
			Edition:         arg.Book.Edition,
			PublisherID:     publisher.PublisherID,
		})
		if err != nil {
			return err
		}

		for i := range authors {
			err = store.CreateAuthorBookRel(ctx, CreateAuthorBookRelParams{
				BookID:   book.BookID,
				AuthorID: authors[i].AuthorID,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return
}
