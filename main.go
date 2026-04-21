package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/segmentio/kafka-go"
	"github.com/zimlewis/shortened/application"
	"github.com/zimlewis/shortened/kafkastream"
)

func main() {
	eventChannel := make(chan []byte)
	defer close(eventChannel)
	app := application.New(eventChannel, badger.DefaultOptions("./tmp/badger"))

	go func() {
		config := kafka.WriterConfig{
			Brokers: []string{string("127.0.0.1:34439")},
			Topic:   "smth",
			Balancer: &kafka.Hash{},
			BatchSize:    100,               // how many messages to batch before flushing
			BatchTimeout: 10 * time.Millisecond, // flush at least every 10ms even if batch isn't full
			Async:        false,             // true = fire and forget, no error returned
		}

		publisher := kafkastream.NewPublisher(
			context.Background(),
			config,
		)

		err := publisher.Start(eventChannel)
		if err != nil {
			fmt.Printf("Cannot start kafka connection: %s", err.Error())
			panic(1)
		}
	}()

	go func () {
		config := kafka.ReaderConfig{
			Brokers:   []string{string("127.0.0.1:34439")},
			Topic:     "smth",
			GroupID:   "group-0",
			MinBytes:  1,
			MaxBytes:  1e6,
			MaxWait:   time.Millisecond,
		}
		consumer := kafkastream.NewConsumer(	
			context.Background(),
			config,
		)
		
		err := consumer.Start()
		if err != nil {
			fmt.Printf("Cannot start kafka connection: %s", err.Error())
			panic(1)
		}
	}()
	err := app.Start()
	if err != nil {
		fmt.Printf("Cannot start application: %s", err.Error())
	}
}
