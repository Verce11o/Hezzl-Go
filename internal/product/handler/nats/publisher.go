package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Verce11o/Hezzl-Go/internal/models"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

type Publisher struct {
	stream jetstream.Stream
	js     jetstream.JetStream
	log    *zap.SugaredLogger
}

func NewPublisher(stream jetstream.Stream, js jetstream.JetStream, log *zap.SugaredLogger) *Publisher {
	return &Publisher{stream: stream, js: js, log: log}
}

func (p *Publisher) Publish(ctx context.Context, product models.Product) error {

	productBytes, err := json.Marshal(product)

	if err != nil {
		p.log.Errorf("error marshaling product: %v", err)
		return err
	}

	fmt.Println(p.js)
	_, err = p.js.PublishAsync("GOODS.NEW", productBytes)

	if err != nil {
		p.log.Errorf("cannot publish new message to stream: %v", err)
		return err
	}

	return nil
}
