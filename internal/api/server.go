package api

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/emiliogozo/xyz-books/db/sqlc"
	"github.com/emiliogozo/xyz-books/internal/util"

	"github.com/gin-contrib/cors"
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
func NewServer(config util.Config, store db.Store) (*Server, error) {
	server := &Server{
		config: config,
		store:  store,
	}

	server.setupRouter()

	return server, nil
}

func (s *Server) setupRouter() {
	gin.SetMode(s.config.GinMode)
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowMethods("OPTIONS")
	corsConfig.AddAllowHeaders("Authorization")
	r.Use(cors.New(corsConfig))

	api := r.Group(s.config.APIBasePath)

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

	if s.config.GinMode != gin.TestMode {
		r.LoadHTMLFiles(fmt.Sprintf("%s/index.html", s.config.WebDistPath))
		r.Static("/assets", fmt.Sprintf("%s/assets", s.config.WebDistPath))
		r.GET("/", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "index.html", gin.H{})
		})
	}

	s.router = r
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
