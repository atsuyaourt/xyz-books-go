package handlers

import (
	"errors"
	"net/http"

	db "github.com/atsuyaourt/xyz-books/internal/db/sqlc"
	"github.com/atsuyaourt/xyz-books/internal/models"
	"github.com/atsuyaourt/xyz-books/internal/services"
	"github.com/gin-gonic/gin"
)

func newPublisher(arg db.Publisher) models.Publisher {
	return models.Publisher{
		PublisherName: arg.PublisherName,
	}
}

// CreatePublisher
//
//	@Summary	Create publisher
//	@Tags		publishers
//	@Accept		json
//	@Produce	json
//	@Param		req	body		services.CreateAuthorReq	true	"Create publisher parameters"
//	@Success	201	{object}	models.Publisher
//	@Router		/publishers [post]
func (h *DefaultHandler) CreatePublisher(ctx *gin.Context) {
	var req services.CreatePublisherReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res, err := h.service.CreatePublisher(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, res)
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

	res, err := h.service.GetPublisher(ctx, req.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("publisher not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// ListPublishers
//
//	@Summary	List publishers
//	@Tags		publishers
//	@Accept		json
//	@Produce	json
//	@Param		req	query		services.ListPublishersReq	false	"List publishers parameters"
//	@Success	200	{object}	models.PaginatedPublishers
//	@Router		/publishers [get]
func (h *DefaultHandler) ListPublishers(ctx *gin.Context) {
	var req services.ListPublishersReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res, err := h.service.ListPublishers(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}

type updatePublisherUri struct {
	ID int64 `uri:"id" binding:"required,numeric"`
}

// UpdatePublisher
//
//	@Summary	Update publisher
//	@Tags		publishers
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"publisher ID"
//	@Param		req		body		services.UpdatePublisherReq	true	"Update publisher parameters"
//	@Success	200		{object}	models.Publisher
//	@Router		/publishers/{id} [put]
func (h *DefaultHandler) UpdatePublisher(ctx *gin.Context) {
	var uri updatePublisherUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req services.UpdatePublisherReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res, err := h.service.UpdatePublisher(ctx, uri.ID, req)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("publisher not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
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

	err := h.service.DeletePublisher(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
