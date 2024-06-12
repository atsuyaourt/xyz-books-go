package handlers

import (
	db "github.com/atsuyaourt/xyz-books/internal/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	CreateBook(ctx *gin.Context)
	ListBooks(ctx *gin.Context)
	GetBook(ctx *gin.Context)
	UpdateBook(ctx *gin.Context)
	DeleteBook(ctx *gin.Context)
	CreateAuthor(ctx *gin.Context)
	ListAuthors(ctx *gin.Context)
	GetAuthor(ctx *gin.Context)
	UpdateAuthor(ctx *gin.Context)
	DeleteAuthor(ctx *gin.Context)
	CreatePublisher(ctx *gin.Context)
	ListPublishers(ctx *gin.Context)
	GetPublisher(ctx *gin.Context)
	UpdatePublisher(ctx *gin.Context)
	DeletePublisher(ctx *gin.Context)
}

type DefaultHandler struct {
	store db.Store
}

func NewDefaultHandler(store db.Store) (Handler, error) {
	h := &DefaultHandler{
		store: store,
	}

	return h, nil
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
