package nats

import (
	"context"
	"encoding/json"
	"github.com/Verce11o/Hezzl-Go/internal/models"
	"github.com/Verce11o/Hezzl-Go/internal/product"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"sync"
)

const (
	workers = 10
	bufMax  = 10
)

type Consumer struct {
	stream jetstream.Stream
	log    *zap.SugaredLogger
	cons   jetstream.Consumer
	buf    []models.Product
	click  product.ClickHouseRepository
}

func NewConsumer(stream jetstream.Stream, log *zap.SugaredLogger, click product.ClickHouseRepository) *Consumer {
	return &Consumer{stream: stream, log: log, buf: make([]models.Product, 0, bufMax), click: click}
}

func (c *Consumer) Consume(ctx context.Context) error {
	cons, err := c.stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:     "ClickHouse",
		Description: "Consumer to add data into ClickHouse",
	})

	if err != nil {
		c.log.Errorf("cannot create consumer: %v", err)
		return err
	}

	c.cons = cons

	wg := &sync.WaitGroup{}
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go c.worker(ctx, wg, cons)
	}

	wg.Wait()

	return nil
}

func (c *Consumer) worker(ctx context.Context, wg *sync.WaitGroup, cons jetstream.Consumer) {
	defer wg.Done()

	_, err := cons.Consume(c.processMessage(ctx))

	if err != nil {
		c.log.Errorf("cannot process message: %v", err)
		return
	}

}

func (c *Consumer) processMessage(ctx context.Context) func(msg jetstream.Msg) {
	return func(msg jetstream.Msg) {
		defer msg.Ack()

		var prod models.Product
		err := json.Unmarshal(msg.Data(), &prod)
		if err != nil {
			c.log.Errorf("cannot unmarshal request: %v", prod)
			return
		}

		if len(c.buf) == cap(c.buf) {
			if err := c.click.UploadEvent(ctx, c.buf); err != nil {
				c.log.Errorf("cannot upload event: %v", err)
				return
			}
			c.buf = c.buf[:0]
			return
		}

		c.buf = append(c.buf, prod)

	}
}
