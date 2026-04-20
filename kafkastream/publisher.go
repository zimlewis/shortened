package kafkastream

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)


type Publisher struct {
	ctx    context.Context
	config kafka.WriterConfig
}

func NewPublisher(
	ctx    context.Context,
	config kafka.WriterConfig,
) *Publisher {
	return &Publisher{
		ctx:       ctx,
		config: config,
	}
}

func (self *Publisher) Start(c chan []byte) error {
	w := kafka.NewWriter(self.config)

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
