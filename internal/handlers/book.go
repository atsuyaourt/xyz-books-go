package handlers

import (
	"errors"
	"net/http"

	db "github.com/atsuyaourt/xyz-books/internal/db/sqlc"
	"github.com/atsuyaourt/xyz-books/internal/services"
	"github.com/gin-gonic/gin"
)

// CreateBook
//
//	@Summary	Create book
//	@Tags		books
//	@Accept		json
//	@Produce	json
//	@Param		req	body		services.CreateBookReq	true	"Create book parameters"
//	@Success	201	{object}	models.Book
//	@Router		/books [post]
func (h *DefaultHandler) CreateBook(ctx *gin.Context) {
	var req services.CreateBookReq
	var err error
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res, err := h.service.CreateBook(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

type getBookReq struct {
	ISBN13 string `uri:"isbn" binding:"required,isbn13"`
}

// GetBook
//
//	@Summary	Get book
//	@Tags		books
//	@Accept		json
//	@Produce	json
//	@Param		isbn	path		string	true	"ISBN-13"
//	@Success	200		{object}	models.Book
//	@Router		/books/{isbn} [get]
func (h *DefaultHandler) GetBook(ctx *gin.Context) {
	var req getBookReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res, err := h.service.GetBook(ctx, req.ISBN13)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("book not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// ListBooks
//
//	@Summary	List books
//	@Tags		books
//	@Accept		json
//	@Produce	json
//	@Param		req	query		services.ListBooksReq	false	"List books parameters"
//	@Success	200	{object}	models.PaginatedBooks
//	@Router		/books [get]
func (h *DefaultHandler) ListBooks(ctx *gin.Context) {
	var req services.ListBooksReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res, err := h.service.ListBooks(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}

type updateBookUri struct {
	ISBN13 string `uri:"isbn" binding:"required,isbn13"`
}

// UpdateBook
//
//	@Summary	Update book
//	@Tags		books
//	@Accept		json
//	@Produce	json
//	@Param		isbn	path		string			true	"ISBN-13"
//	@Param		req		body		services.UpdateBookReq	true	"Update book parameters"
//	@Success	200		{object}	models.Book
//	@Router		/books/{isbn} [put]
func (h *DefaultHandler) UpdateBook(ctx *gin.Context) {
	var uri updateBookUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req services.UpdateBookReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res, err := h.service.UpdateBook(ctx, uri.ISBN13, req)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("book not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}

type deleteBookUri struct {
	ISBN13 string `uri:"isbn" binding:"required,isbn13"`
}

// DeleteBook
//
//	@Summary	Delete book
//	@Tags		books
//	@Accept		json
//	@Produce	json
//	@Param		isbn	path	string	true	"ISBN-13"
//	@Success	204
//	@Router		/books/{isbn} [delete]
func (h *DefaultHandler) DeleteBook(ctx *gin.Context) {
	var req deleteBookUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := h.service.DeleteBook(ctx, req.ISBN13)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
