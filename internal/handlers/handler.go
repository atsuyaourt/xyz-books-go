package handlers

import (
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

	Index(ctx *gin.Context)

	ShowBooks(ctx *gin.Context)
	ShowBook(ctx *gin.Context)
}
