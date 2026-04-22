package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/segmentio/kafka-go"
	"github.com/zimlewis/shortened/application"
	"github.com/zimlewis/shortened/global"
	"github.com/zimlewis/shortened/kafkastream"
	"github.com/zimlewis/shortened/repository"
)

func main() {
	badgerLocation := os.Getenv("BADGER_DIR")
	if badgerLocation == "" {
		fmt.Println("Cannot get badger directory in environment variables")
		panic(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("Cannot get port in environment variables")
		panic(1)
	}

	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		fmt.Println("Cannot get broker in environment variables")
		panic(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	eventChannel := make(chan []byte)
	defer close(eventChannel)

	badgerOptions := badger.DefaultOptions(badgerLocation)
	badgerOptions.Logger = nil
	repo, err := repository.NewBadger(&badgerOptions)
	if err != nil {
		fmt.Printf("Cannot start badger db: %s", err.Error())
		os.Exit(1)
	}
	defer repo.Close()

	config := &global.Config {
		Repository: repo,
		WriteMessageChannel: eventChannel,
		WriterConfig: kafka.WriterConfig{
			Brokers: []string{broker},
			Topic:   "url-click",
			Balancer: &kafka.Hash{},
			BatchSize:    1,               // how many messages to batch before flushing
			BatchTimeout: 500 * time.Millisecond, // flush at least every 10ms even if batch isn't full
			Async:        false,             // true = fire and forget, no error returned
		},
		ReaderConfig: kafka.ReaderConfig{
			Brokers:   []string{broker},
			Topic:     "url-click",
			GroupID:   "group-1",
			MinBytes:  1,
			MaxBytes:  1e6,
			MaxWait:   500 * time.Millisecond,
			StartOffset: kafka.FirstOffset,
		},
	}

	app := application.New(eventChannel, config)

	// Writer goroutin
	go func() {

		publisher := kafkastream.NewPublisher(
			config,
		)

		err := publisher.Start(ctx)
		if err != nil {
			fmt.Printf("Cannot start kafka connection: %s", err.Error())
			panic(1)
		}
	}()

	// Reader goroutin
	go func () {
		consumer := kafkastream.NewConsumer(
			config,
		)

		err := consumer.Start(ctx)
		if err != nil {
			fmt.Printf("Cannot start kafka connection: %s", err.Error())
			panic(1)
		}
	}()

	err = app.Start(ctx, port)
	if err != nil {
		fmt.Printf("Cannot start application: %s", err.Error())
	}
}
