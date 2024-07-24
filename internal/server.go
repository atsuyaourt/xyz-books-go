package internal

import (
	"log"
	"net/http"

	db "github.com/atsuyaourt/xyz-books/internal/db/sqlc"
	"github.com/atsuyaourt/xyz-books/internal/handlers"
	"github.com/atsuyaourt/xyz-books/internal/util"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	config  util.Config
	router  *gin.Engine
	store   db.Store
	handler handlers.Handler
}

// NewServer creates a new HTTP server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	server := &Server{
		config: config,
		store:  store,
	}
	gin.SetMode(config.GinMode)
	server.router = gin.Default()

	handler, err := handlers.NewDefaultHandler(store)
	if err != nil {
		return nil, err
	}
	server.handler = handler

	server.setupCORS()
	server.setupRouter()
	server.setupAPIRouter()

	return server, nil
}

func (s *Server) setupCORS() {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowMethods("OPTIONS")
	corsConfig.AddAllowHeaders("Authorization")
	s.router.Use(cors.New(corsConfig))
}

func (s *Server) setupRouter() {
	r := s.router

	r.Static("/assets", "internal/assets")
	r.GET("/", s.handler.Index)
	r.GET("/:isbn", s.handler.ShowBook)

	r.GET("/books", s.handler.ShowBooks)
	r.GET("/books/:isbn", s.handler.ShowBook)
}

func (s *Server) setupAPIRouter() {
	api := s.router.Group(s.config.APIBasePath)

	books := api.Group("/books")
	{
		books.GET("", s.handler.ListBooks)
		books.GET(":isbn", s.handler.GetBook)
		books.POST("", s.handler.CreateBook)
		books.PUT(":isbn", s.handler.UpdateBook)
		books.DELETE(":isbn", s.handler.DeleteBook)
	}

	authors := api.Group("/authors")
	{
		authors.GET("", s.handler.ListAuthors)
		authors.GET(":id", s.handler.GetAuthor)
		authors.POST("", s.handler.CreateAuthor)
		authors.PUT(":id", s.handler.UpdateAuthor)
		authors.DELETE(":id", s.handler.DeleteAuthor)
	}

	publishers := api.Group("/publishers")
	{
		publishers.GET("", s.handler.ListPublishers)
		publishers.GET(":id", s.handler.GetPublisher)
		publishers.POST("", s.handler.CreatePublisher)
		publishers.PUT(":id", s.handler.UpdatePublisher)
		publishers.DELETE(":id", s.handler.DeletePublisher)
	}

	api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (s *Server) Start(address string) *http.Server {
	srv := &http.Server{
		Addr:    address,
		Handler: s.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return srv
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
