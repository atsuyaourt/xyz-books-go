package services

import (
	"database/sql"

	db "github.com/atsuyaourt/xyz-books/internal/db/sqlc"
	"github.com/atsuyaourt/xyz-books/internal/models"
	"github.com/atsuyaourt/xyz-books/internal/util"
	"golang.org/x/net/context"
)

func newAuthor(arg db.Author) models.Author {
	return models.Author{
		FirstName:  arg.FirstName,
		LastName:   arg.LastName,
		MiddleName: arg.LastName,
	}
}

type CreateAuthorReq struct {
	FirstName  string `json:"first_name" binding:"required,min=1"`
	LastName   string `json:"last_name" binding:"required,min=1"`
	MiddleName string `json:"middle_name" binding:"omitempty,min=1"`
} //@name CreateAuthorParams

func (s *DefaultService) CreateAuthor(ctx context.Context, req CreateAuthorReq) (*models.Author, error) {
	arg := db.CreateAuthorParams(req)

	author, err := s.store.CreateAuthor(ctx, arg)
	if err != nil {
		return nil, err
	}

	res := newAuthor(author)

	return &res, nil
}

func (s *DefaultService) GetAuthor(ctx context.Context, id int64) (*models.Author, error) {
	author, err := s.store.GetAuthor(ctx, id)
	if err != nil {
		return nil, err
	}

	res := newAuthor(author)

	return &res, nil
}

type ListAuthorsReq struct {
	Page    int32 `form:"page,default=1" binding:"omitempty,min=1"`     // page number
	PerPage int32 `form:"per_page,default=5" binding:"omitempty,min=1"` // limit
} //@name ListAuthorsParams

func (s *DefaultService) ListAuthors(ctx context.Context, req ListAuthorsReq) (*util.PaginatedList[models.Author], error) {
	offset := (req.Page - 1) * req.PerPage

	arg := db.ListAuthorsParams{
		Limit:  int64(req.PerPage),
		Offset: int64(offset),
	}
	authors, err := s.store.ListAuthors(ctx, arg)
	if err != nil {
		return nil, err
	}

	numAuthors := len(authors)
	items := make([]models.Author, numAuthors)
	for i, author := range authors {
		items[i] = newAuthor(author)
	}

	count, err := s.store.CountAuthors(ctx)
	if err != nil {
		return nil, err
	}

	res := util.NewPaginatedList(req.Page, req.PerPage, int32(count), items)

	return &res, nil
}

type UpdateAuthorReq struct {
	FirstName  string `json:"first_name" binding:"omitempty,min=1"`
	LastName   string `json:"last_name" binding:"omitempty,min=1"`
	MiddleName string `json:"middle_name" binding:"omitempty,min=1"`
} //@name UpdateAuthorParams

func (s *DefaultService) UpdateAuthor(ctx context.Context, oldID int64, req UpdateAuthorReq) (*models.Author, error) {
	arg := db.UpdateAuthorParams{
		AuthorID: oldID,
		FirstName: sql.NullString{
			String: req.FirstName,
			Valid:  len(req.FirstName) > 0,
		},
		LastName: sql.NullString{
			String: req.LastName,
			Valid:  len(req.LastName) > 0,
		},
		MiddleName: sql.NullString{
			String: req.MiddleName,
			Valid:  len(req.MiddleName) > 0,
		},
	}

	author, err := s.store.UpdateAuthor(ctx, arg)
	if err != nil {
		return nil, err
	}

	res := newAuthor(author)

	return &res, nil
}

func (s *DefaultService) DeleteAuthor(ctx context.Context, id int64) error {
	return s.store.DeleteAuthor(ctx, id)
}
