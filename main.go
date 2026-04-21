package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/segmentio/kafka-go"
	"github.com/zimlewis/shortened/application"
	"github.com/zimlewis/shortened/global"
	"github.com/zimlewis/shortened/kafkastream"
)

func main() {
	eventChannel := make(chan []byte)
	defer close(eventChannel)

	config := &global.Config {
		BadgerOptions: badger.DefaultOptions("./tmp/badger/"),
		WriteMessageChannel: eventChannel,
		WriterConfig: kafka.WriterConfig{
			Brokers: []string{string("127.0.0.1:34439")},
			Topic:   "smth",
			Balancer: &kafka.Hash{},
			BatchSize:    100,               // how many messages to batch before flushing
			BatchTimeout: 10 * time.Millisecond, // flush at least every 10ms even if batch isn't full
			Async:        false,             // true = fire and forget, no error returned
		},
		ReaderConfig: kafka.ReaderConfig{
			Brokers:   []string{string("127.0.0.1:34439")},
			Topic:     "smth",
			GroupID:   "group-0",
			MinBytes:  1,
			MaxBytes:  1e6,
			MaxWait:   time.Millisecond,
		},
	}

	app := application.New(eventChannel, config)

	// Writer goroutin
	go func() {

		publisher := kafkastream.NewPublisher(
			context.Background(),
			config,
		)

		err := publisher.Start()
		if err != nil {
			fmt.Printf("Cannot start kafka connection: %s", err.Error())
			panic(1)
		}
	}()

	// Reader goroutin
	go func () {
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
