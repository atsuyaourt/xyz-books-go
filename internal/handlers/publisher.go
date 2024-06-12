package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/atsuyaourt/xyz-books/internal/db/sqlc"
	"github.com/atsuyaourt/xyz-books/internal/models"
	"github.com/atsuyaourt/xyz-books/internal/util"
	"github.com/gin-gonic/gin"
)

func newPublisher(arg db.Publisher) models.Publisher {
	return models.Publisher{
		PublisherName: arg.PublisherName,
	}
}

type createPublisherReq struct {
	PublisherName string `json:"publisher_name" binding:"required,min=1"`
} //@name CreatePublisherParams

// CreatePublisher
//
//	@Summary	Create publisher
//	@Tags		publishers
//	@Accept		json
//	@Produce	json
//	@Param		req	body		createAuthorReq	true	"Create publisher parameters"
//	@Success	201	{object}	models.Publisher
//	@Router		/publishers [post]
func (h *DefaultHandler) CreatePublisher(ctx *gin.Context) {
	var req createPublisherReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	publisher, err := h.store.CreatePublisher(ctx, req.PublisherName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, newPublisher(publisher))
}

type getPublisherReq struct {
	ID int64 `uri:"id" binding:"required,numeric"`
}

// GetPublisher
//
//	@Summary	Get publisher
//	@Tags		publishers
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"publisher ID"
//	@Success	200		{object}	models.Publisher
//	@Router		/publishers/{id} [get]
func (h *DefaultHandler) GetPublisher(ctx *gin.Context) {
	var req getPublisherReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	publisher, err := h.store.GetPublisher(ctx, req.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("publisher not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newPublisher(publisher))
}

type listPublishersReq struct {
	Page    int32 `form:"page,default=1" binding:"omitempty,min=1"`            // page number
	PerPage int32 `form:"per_page,default=5" binding:"omitempty,min=1,max=30"` // limit
} //@name ListPublishersParams

// ListPublishers
//
//	@Summary	List publishers
//	@Tags		publishers
//	@Accept		json
//	@Produce	json
//	@Param		req	query		listPublishersReq	false	"List publishers parameters"
//	@Success	200	{object}	models.PaginatedPublishers
//	@Router		/publishers [get]
func (h *DefaultHandler) ListPublishers(ctx *gin.Context) {
	var req listPublishersReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	offset := (req.Page - 1) * req.PerPage

	arg := db.ListPublishersParams{
		Limit:  int64(req.PerPage),
		Offset: int64(offset),
	}
	publishers, err := h.store.ListPublishers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	numPublishers := len(publishers)
	items := make([]models.Publisher, numPublishers)
	for i, publisher := range publishers {
		items[i] = newPublisher(publisher)
	}

	count, err := h.store.CountPublishers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := util.NewPaginatedList(req.Page, req.PerPage, int32(count), items)

	ctx.JSON(http.StatusOK, res)
}

type updatePublisherUri struct {
	ID int64 `uri:"id" binding:"required,numeric"`
}

type updatePublisherReq struct {
	PublisherName string `json:"publisher_name" binding:"omitempty,min=1"`
} //@name UpdatePublisherParams

// UpdatePublisher
//
//	@Summary	Update publisher
//	@Tags		publishers
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"publisher ID"
//	@Param		req		body		updatePublisherReq	true	"Update publisher parameters"
//	@Success	200		{object}	models.Publisher
//	@Router		/publishers/{id} [put]
func (h *DefaultHandler) UpdatePublisher(ctx *gin.Context) {
	var uri updatePublisherUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updatePublisherReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdatePublisherParams{
		PublisherID: uri.ID,
		PublisherName: sql.NullString{
			String: req.PublisherName,
			Valid:  len(req.PublisherName) > 0,
		},
	}

	publisher, err := h.store.UpdatePublisher(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("publisher not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newPublisher(publisher))
}

type deletePublisherUri struct {
	ID int64 `uri:"id" binding:"required,numeric"`
}

// DeletePublisher
//
//	@Summary	Delete publisher
//	@Tags		publishers
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"publisher ID"
//	@Success	204
//	@Router		/publishers/{id} [delete]
func (h *DefaultHandler) DeletePublisher(ctx *gin.Context) {
	var req deletePublisherUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := h.store.DeletePublisher(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
