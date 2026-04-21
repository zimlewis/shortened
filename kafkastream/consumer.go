package kafkastream

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/zimlewis/shortened/global"
)

type Consumer struct {
	ctx context.Context
	config *global.Config
}

func NewConsumer(
	ctx    context.Context,
	config *global.Config,
) *Consumer {
	return &Consumer {
		ctx: ctx,
		config: config,
	}
}

func (self *Consumer) Start() error {
	r := kafka.NewReader(self.config.ReaderConfig)
	defer r.Close()

	for {
		msg, err := r.FetchMessage(self.ctx) // blocks until message OR ctx cancelled
		if err != nil {
			return fmt.Errorf("Cannot read message: %w", err) 
		}

		fmt.Printf("Get message: %s\n", string(msg.Value))
		err = r.CommitMessages(self.ctx, msg)
		if err != nil {
			return fmt.Errorf("Cannot commit message: %w", err)
		}

	}
}
