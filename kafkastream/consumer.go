package kafkastream

import (
	"context"
	"fmt"
	// "time"

	"github.com/segmentio/kafka-go"
	"github.com/zimlewis/shortened/global"
)

type Consumer struct {
	config *global.Config
}

func NewConsumer(
	config *global.Config,
) *Consumer {
	return &Consumer{
		config: config,
	}
}

func (self *Consumer) Start(ctx context.Context) error {
	reader := kafka.NewReader(self.config.ReaderConfig)
	defer func () {
		fmt.Println("Closing kafka reader...")
		err := reader.Close()
		if err != nil {
			fmt.Printf("Cannot close reader: %s\n", err.Error())
		}
	}()

	repo := self.config.Repository

	for {
		msg, err := reader.FetchMessage(ctx) // blocks until message OR ctx cancelled
		if err != nil {
			// ctx was cancelled or reader closed
			if ctx.Err() != nil {
				break
			}
			fmt.Printf("Cannot read message: %s\n", err.Error())
			continue
		}

		_, err = repo.IncreaseLinkClick(ctx, string(msg.Value))
		if err != nil {
			fmt.Printf("Cannot increase the count: %s\n", err.Error())
			continue
		}

		// fmt.Printf("%s has clicked %d times\n", string(msg.Value), count)
		err = reader.CommitMessages(ctx, msg)
		if err != nil {
			fmt.Printf("Cannot commit message: %s\n", err.Error())
			continue
		}

	}

	return nil
}
