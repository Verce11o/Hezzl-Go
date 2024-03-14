package nats

import (
	"context"
	"errors"
	"github.com/Verce11o/Hezzl-Go/internal/config"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	streamName        = "GOODS"
	streamDescription = "Stream for storing goods"
	subject           = "GOODS.*"
	msgAmount         = 50
)

type Nats struct {
	Conn   *nats.Conn
	Js     jetstream.JetStream
	Stream jetstream.Stream
}

func NewNats(ctx context.Context, cfg *config.Config) *Nats {

	nc, err := nats.Connect(cfg.Nats.Url)
	if err != nil {
		panic(err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		panic(err)
	}

	stream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:        streamName,
		Description: streamDescription,
		Subjects:    []string{subject},
		MaxMsgs:     msgAmount,
		Discard:     jetstream.DiscardNew,
		Storage:     jetstream.FileStorage,
		Retention:   jetstream.WorkQueuePolicy,
	})

	if !errors.Is(err, jetstream.ErrStreamNameAlreadyInUse) && err != nil {
		panic(err)
	}

	return &Nats{
		Conn:   nc,
		Js:     js,
		Stream: stream,
	}

}
