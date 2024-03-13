package server

import (
	"context"
	"fmt"
	"github.com/Verce11o/Hezzl-Go/api"
	"github.com/Verce11o/Hezzl-Go/internal/config"
	productHandler "github.com/Verce11o/Hezzl-Go/internal/product/handler"
	productRepository "github.com/Verce11o/Hezzl-Go/internal/product/repository"
	productService "github.com/Verce11o/Hezzl-Go/internal/product/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	middleware "github.com/oapi-codegen/gin-middleware"

	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

type Server struct {
	log        *zap.SugaredLogger
	cfg        *config.Config
	db         *pgxpool.Pool
	httpServer *http.Server
}

func NewServer(log *zap.SugaredLogger, cfg *config.Config, db *pgxpool.Pool) *Server {
	return &Server{log: log, cfg: cfg, db: db}
}

func (s *Server) Run(handler http.Handler) error {
	addr := fmt.Sprintf("%v:%v", s.cfg.Server.Host, s.cfg.Server.Port)
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	s.log.Infof("Server running on: %v", addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) InitRoutes() *gin.Engine {

	swagger, err := api.GetSwagger()
	if err != nil {
		log.Fatalf("error loading swagger spec\n:%v", err)
	}

	router := gin.Default()

	router.Use(middleware.OapiRequestValidator(swagger))

	productRepo := productRepository.NewProductRepository(s.db)
	productServices := productService.NewService(s.log, productRepo)

	//router.GET("/swagger", swagger)

	apiGroup := router.Group("/api/v1")

	api.RegisterHandlers(apiGroup, productHandler.NewHandler(s.log, productServices))

	return router
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
