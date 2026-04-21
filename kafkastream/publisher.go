package kafkastream

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/zimlewis/shortened/global"
)

type Publisher struct {
	ctx    context.Context
	config *global.Config
}

func NewPublisher(
	ctx context.Context,
	config *global.Config,
) *Publisher {
	return &Publisher{
		ctx:    ctx,
		config: config,
	}
}

func (self *Publisher) Start() error {
	w := kafka.NewWriter(self.config.WriterConfig)
	c := self.config.WriteMessageChannel

	for e := range c {
		go func() {
			err := w.WriteMessages(self.ctx, kafka.Message{
				Value: e,
			})

			if err != nil {
				fmt.Printf("Cannot write message %s\n", err.Error())
			}
		}()
	}

	return nil
}
