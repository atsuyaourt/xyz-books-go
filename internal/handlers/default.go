package handlers

import (
	"github.com/a-h/templ"
	db "github.com/atsuyaourt/xyz-books/internal/db/sqlc"
	"github.com/atsuyaourt/xyz-books/internal/services"

	"github.com/gin-gonic/gin"
)

type DefaultHandler struct {
	service services.Service
}

func NewDefaultHandler(store db.Store) (*DefaultHandler, error) {
	s, err := services.NewDefaultService(store)
	if err != nil {
		return nil, err
	}
	h := &DefaultHandler{
		service: s,
	}

	return h, nil
}

func render(ctx *gin.Context, status int, template templ.Component) error {
	ctx.Status(status)
	return template.Render(ctx.Request.Context(), ctx.Writer)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
