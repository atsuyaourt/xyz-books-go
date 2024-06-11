package api

import (
	db "github.com/atsuyaourt/xyz-books/internal/db/sqlc"
	"github.com/atsuyaourt/xyz-books/internal/util"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	config util.Config
	router *gin.Engine
	store  db.Store
}

// NewServer creates a new HTTP server and setup routing
func NewServer(config util.Config, store db.Store, router *gin.Engine) (*Server, error) {
	server := &Server{
		config: config,
		store:  store,
	}

	if router == nil {
		gin.SetMode(config.GinMode)
		server.router = gin.Default()
	} else {
		server.router = router
	}

	server.setupRouter()

	return server, nil
}

func (s *Server) setupRouter() {
	api := s.router.Group(s.config.APIBasePath)

	books := api.Group("/books")
	{
		books.GET("", s.ListBooks)
		books.GET(":isbn", s.GetBook)
		books.POST("", s.CreateBook)
		books.PUT(":isbn", s.UpdateBook)
		books.DELETE(":isbn", s.DeleteBook)
	}

	authors := api.Group("/authors")
	{
		authors.GET("", s.ListAuthors)
		authors.GET(":id", s.GetAuthor)
		authors.POST("", s.CreateAuthor)
		authors.PUT(":id", s.UpdateAuthor)
		authors.DELETE(":id", s.DeleteAuthor)
	}

	publishers := api.Group("/publishers")
	{
		publishers.GET("", s.ListPublishers)
		publishers.GET(":id", s.GetPublisher)
		publishers.POST("", s.CreatePublisher)
		publishers.PUT(":id", s.UpdatePublisher)
		publishers.DELETE(":id", s.DeletePublisher)
	}

	api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
