package server

import (
	"context"
	"github.com/Verce11o/Hezzl-Go/api"
	"github.com/Verce11o/Hezzl-Go/internal/config"
	productHandler "github.com/Verce11o/Hezzl-Go/internal/product/handler/http/v1"
	"github.com/Verce11o/Hezzl-Go/internal/product/handler/nats"
	productRepository "github.com/Verce11o/Hezzl-Go/internal/product/repository"
	productService "github.com/Verce11o/Hezzl-Go/internal/product/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	middleware "github.com/oapi-codegen/gin-middleware"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

type Server struct {
	log        *zap.SugaredLogger
	cfg        *config.Config
	db         *pgxpool.Pool
	redis      *redis.Client
	httpServer *http.Server
	publisher  *nats.Publisher
}

func NewServer(log *zap.SugaredLogger, cfg *config.Config, db *pgxpool.Pool, redis *redis.Client, publisher *nats.Publisher) *Server {
	return &Server{log: log, cfg: cfg, db: db, redis: redis, publisher: publisher}
}

func (s *Server) Run(handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:         s.cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	s.log.Infof("Server running on: %v", s.cfg.Server.Port)
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
	productCache := productRepository.NewProductsRedis(s.redis)

	productServices := productService.NewService(s.log, productRepo, productCache, s.publisher)

	apiGroup := router.Group("/api/v1")

	api.RegisterHandlers(apiGroup, productHandler.NewHandler(s.log, productServices))

	return router
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
