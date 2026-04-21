package kafkastream

import (
	"context"
	"fmt"
	// "time"

	"github.com/segmentio/kafka-go"
	"github.com/zimlewis/shortened/global"
)

type Consumer struct {
	ctx    context.Context
	config *global.Config
}

func NewConsumer(
	ctx context.Context,
	config *global.Config,
) *Consumer {
	return &Consumer{
		ctx:    ctx,
		config: config,
	}
}

func (self *Consumer) Start() error {
	r := kafka.NewReader(self.config.ReaderConfig)
	defer r.Close()

	repo := self.config.Repository

	for {
		msg, err := r.FetchMessage(self.ctx) // blocks until message OR ctx cancelled

		if err != nil {
			
			fmt.Printf("Cannot read message: %s\n", err.Error())
			continue
		}

		count, err := repo.IncreaseLinkClick(self.ctx, string(msg.Value))
		if err != nil {
			fmt.Printf("Cannot increase the count: %s\n", err.Error())
			continue
		}

		fmt.Printf("%s has clicked %d times\n", string(msg.Value), count)
		err = r.CommitMessages(self.ctx, msg)
		if err != nil {
			fmt.Printf("Cannot commit message: %s\n", err.Error())
			continue
		}

	}
}
