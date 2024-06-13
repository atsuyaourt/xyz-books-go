package handlers

import (
	"errors"
	"net/http"

	db "github.com/atsuyaourt/xyz-books/internal/db/sqlc"
	"github.com/atsuyaourt/xyz-books/internal/services"
	"github.com/gin-gonic/gin"
)

// CreateAuthor
//
//	@Summary	Create author
//	@Tags		authors
//	@Accept		json
//	@Produce	json
//	@Param		req	body		services.CreateAuthorReq	true	"Create author parameters"
//	@Success	201	{object}	models.Author
//	@Router		/authors [post]
func (h *DefaultHandler) CreateAuthor(ctx *gin.Context) {
	var req services.CreateAuthorReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res, err := h.service.CreateAuthor(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

type getAuthorReq struct {
	ID int64 `uri:"id" binding:"required,numeric"`
}

// GetAuthor
//
//	@Summary	Get author
//	@Tags		authors
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"author ID"
//	@Success	200	{object}	models.Author
//	@Router		/authors/{id} [get]
func (h *DefaultHandler) GetAuthor(ctx *gin.Context) {
	var req getAuthorReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res, err := h.service.GetAuthor(ctx, req.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("author not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// ListAuthors
//
//	@Summary	List authors
//	@Tags		authors
//	@Accept		json
//	@Produce	json
//	@Param		req	query		services.ListAuthorsReq	false	"List authors parameters"
//	@Success	200	{object}	models.PaginatedAuthors
//	@Router		/authors [get]
func (h *DefaultHandler) ListAuthors(ctx *gin.Context) {
	var req services.ListAuthorsReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res, err := h.service.ListAuthors(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}

type updateAuthorUri struct {
	ID int64 `uri:"id" binding:"required,numeric"`
}

// UpdateAuthor
//
//	@Summary	Update author
//	@Tags		authors
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int				true	"author ID"
//	@Param		req	body		services.UpdateAuthorReq	true	"Update author parameters"
//	@Success	200	{object}	models.Author
//	@Router		/authors/{id} [put]
func (h *DefaultHandler) UpdateAuthor(ctx *gin.Context) {
	var uri updateAuthorUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req services.UpdateAuthorReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res, err := h.service.UpdateAuthor(ctx, uri.ID, req)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("author not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}

type deleteAuthorUri struct {
	ID int64 `uri:"id" binding:"required,numeric"`
}

// DeleteAuthor
//
//	@Summary	Delete author
//	@Tags		authors
//	@Accept		json
//	@Produce	json
//	@Param		id	path	int	true	"author ID"
//	@Success	204
//	@Router		/authors/{id} [delete]
func (h *DefaultHandler) DeleteAuthor(ctx *gin.Context) {
	var req deleteAuthorUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := h.service.DeleteAuthor(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
