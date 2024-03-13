package app

import (
	"context"
	"github.com/Verce11o/Hezzl-Go/internal/config"
	"github.com/Verce11o/Hezzl-Go/internal/server"
	"github.com/Verce11o/Hezzl-Go/lib/db/postgres"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func Run(log *zap.SugaredLogger, cfg *config.Config) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := postgres.NewPostgresPool(ctx, cfg)
	srv := server.NewServer(log, cfg, db)

	go func() {
		if err := srv.Run(srv.InitRoutes()); err != nil {
			log.Fatalf("Error while start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Infof("Error occured on server shutting down: %v", err)
	}

}
