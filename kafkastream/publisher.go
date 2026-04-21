package kafkastream

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/zimlewis/shortened/global"
)

type Publisher struct {
	config *global.Config
}

func NewPublisher(
	config *global.Config,
) *Publisher {
	return &Publisher{
		config: config,
	}
}

func (self *Publisher) Start(ctx context.Context) error {
	writer := kafka.NewWriter(self.config.WriterConfig)
	defer func () {
		fmt.Println("Closing kafka publisher...")
		err := writer.Close()
		if err != nil {
			fmt.Printf("Cannot close publisher: %s\n", err.Error())
		}
	}()
	messageChan := self.config.WriteMessageChannel

	loop:
	for {
		var messageToWrite []byte

		select {
		case <- ctx.Done(): break loop 
		case messageToWrite = <- messageChan:
		}

		err := writer.WriteMessages(ctx, kafka.Message{
			Value: messageToWrite,
		})

		if err != nil {
			fmt.Printf("Cannot write message %s\n", err.Error())
		}
	}

	return nil
}
