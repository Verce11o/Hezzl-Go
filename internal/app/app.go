package app

import (
	"context"
	"github.com/Verce11o/Hezzl-Go/internal/config"
	productNats "github.com/Verce11o/Hezzl-Go/internal/product/handler/nats"
	"github.com/Verce11o/Hezzl-Go/internal/product/repository"
	"github.com/Verce11o/Hezzl-Go/internal/server"
	"github.com/Verce11o/Hezzl-Go/lib/db/clickhouse"
	"github.com/Verce11o/Hezzl-Go/lib/db/postgres"
	"github.com/Verce11o/Hezzl-Go/lib/db/redis"
	"github.com/Verce11o/Hezzl-Go/lib/nats"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func Run(log *zap.SugaredLogger, cfg *config.Config) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := postgres.NewPostgresPool(ctx, cfg)
	redisClient := redis.NewRedisClient(ctx, cfg)
	natsClient := nats.NewNats(ctx, cfg)
	clickHouseClient := clickhouse.NewClickHouseClient(ctx, cfg)
	clickHouseDB := repository.NewProductClickHouse(clickHouseClient)

	publisher := productNats.NewPublisher(natsClient.Stream, natsClient.Js, log)
	srv := server.NewServer(log, cfg, db, redisClient, publisher)

	consumer := productNats.NewConsumer(natsClient.Stream, log, clickHouseDB)

	go func() {
		if err := consumer.Consume(ctx); err != nil {
			log.Fatalf("Error while listening to nats: %v", err)
		}
	}()

	go func() {
		if err := srv.Run(srv.InitRoutes()); err != nil {
			log.Fatalf("Error while start server: %v", err)
		}
	}()
	//
	//for i := 0; i < 100; i++ {
	//	publisher.Publish(ctx, models.Product{
	//		ID:          rand.IntN(9999),
	//		ProjectID:   1,
	//		Name:        fmt.Sprintf("Test: %d", rand.IntN(9999)),
	//		Description: "Wow Description",
	//		Priority:    0,
	//		Removed:     false,
	//		CreatedAt:   time.Time{},
	//	})
	//}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Infof("Error occured on server shutting down: %v", err)
	}

}
