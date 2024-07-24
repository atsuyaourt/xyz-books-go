package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	db "github.com/atsuyaourt/xyz-books/internal/db/sqlc"
	"github.com/atsuyaourt/xyz-books/internal/services"
	"github.com/atsuyaourt/xyz-books/internal/views"
	"github.com/atsuyaourt/xyz-books/internal/views/components"
	"github.com/gin-gonic/gin"
)

func (h *DefaultHandler) ShowBooks(ctx *gin.Context) {
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

	var qSlice []string
	for k, v := range ctx.Request.URL.Query() {
		if len(v) == 1 && len(v[0]) > 0 {
			qSlice = append(qSlice, fmt.Sprintf("%v=%v", k, v[0]))
		}
	}
	qStr := "?" + strings.Join(qSlice, "&")

	ctx.Writer.Header().Set("HX-Replace-Url", qStr)

	render(ctx, http.StatusOK, components.Books(*res))
}

func (h *DefaultHandler) ShowBook(ctx *gin.Context) {
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

	render(ctx, http.StatusOK, views.Book(res))
}
