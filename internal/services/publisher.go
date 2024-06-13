package services

import (
	"database/sql"

	db "github.com/atsuyaourt/xyz-books/internal/db/sqlc"
	"github.com/atsuyaourt/xyz-books/internal/models"
	"github.com/atsuyaourt/xyz-books/internal/util"
	"golang.org/x/net/context"
)

func newPublisher(arg db.Publisher) models.Publisher {
	return models.Publisher{
		PublisherName: arg.PublisherName,
	}
}

type CreatePublisherReq struct {
	PublisherName string `json:"publisher_name" binding:"required,min=1"`
} //@name CreatePublisherParams

func (s *DefaultService) CreatePublisher(ctx context.Context, req CreatePublisherReq) (*models.Publisher, error) {
	publisher, err := s.store.CreatePublisher(ctx, req.PublisherName)
	if err != nil {
		return nil, err
	}

	res := newPublisher(publisher)

	return &res, nil
}

func (s *DefaultService) GetPublisher(ctx context.Context, id int64) (*models.Publisher, error) {
	publisher, err := s.store.GetPublisher(ctx, id)
	if err != nil {
		return nil, err
	}

	res := newPublisher(publisher)

	return &res, nil
}

type ListPublishersReq struct {
	Page    int32 `form:"page,default=1" binding:"omitempty,min=1"`            // page number
	PerPage int32 `form:"per_page,default=5" binding:"omitempty,min=1,max=30"` // limit
} //@name ListPublishersParams

func (s *DefaultService) ListPublishers(ctx context.Context, req ListPublishersReq) (*util.PaginatedList[models.Publisher], error) {
	offset := (req.Page - 1) * req.PerPage

	arg := db.ListPublishersParams{
		Limit:  int64(req.PerPage),
		Offset: int64(offset),
	}
	publishers, err := s.store.ListPublishers(ctx, arg)
	if err != nil {
		return nil, err
	}

	numPublishers := len(publishers)
	items := make([]models.Publisher, numPublishers)
	for i, publisher := range publishers {
		items[i] = newPublisher(publisher)
	}

	count, err := s.store.CountPublishers(ctx)
	if err != nil {
		return nil, err
	}

	res := util.NewPaginatedList(req.Page, req.PerPage, int32(count), items)

	return &res, nil
}

type UpdatePublisherReq struct {
	PublisherName string `json:"publisher_name" binding:"omitempty,min=1"`
} //@name UpdatePublisherParams

func (s *DefaultService) UpdatePublisher(ctx context.Context, oldID int64, req UpdatePublisherReq) (*models.Publisher, error) {
	arg := db.UpdatePublisherParams{
		PublisherID: oldID,
		PublisherName: sql.NullString{
			String: req.PublisherName,
			Valid:  len(req.PublisherName) > 0,
		},
	}

	publisher, err := s.store.UpdatePublisher(ctx, arg)
	if err != nil {
		return nil, err
	}

	res := newPublisher(publisher)

	return &res, nil
}

func (s *DefaultService) DeletePublisher(ctx context.Context, id int64) error {
	return s.store.DeletePublisher(ctx, id)
}
