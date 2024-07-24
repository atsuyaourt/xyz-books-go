package handlers

import (
	"net/http"

	"github.com/atsuyaourt/xyz-books/internal/services"
	"github.com/atsuyaourt/xyz-books/internal/views"
	"github.com/gin-gonic/gin"
)

func (h *DefaultHandler) Index(ctx *gin.Context) {
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

	render(ctx, http.StatusOK, views.Books(*res))
}
