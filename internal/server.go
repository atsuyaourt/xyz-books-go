package internal

import (
	"fmt"
	"log"
	"net/http"

	"github.com/atsuyaourt/xyz-books/internal/api"
	db "github.com/atsuyaourt/xyz-books/internal/db/sqlc"
	"github.com/atsuyaourt/xyz-books/internal/util"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	gin.SetMode(config.GinMode)
	server.router = gin.Default()

	server.setupCORS()
	server.setupRouter()

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

	api.NewServer(s.config, s.store, r)

	if s.config.GinMode != gin.TestMode {
		r.LoadHTMLFiles(fmt.Sprintf("%s/index.html", s.config.WebDistPath))
		r.Static("/assets", fmt.Sprintf("%s/assets", s.config.WebDistPath))
		r.GET("/", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "index.html", gin.H{})
		})
	}
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
