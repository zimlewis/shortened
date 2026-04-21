package global

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/segmentio/kafka-go"
)

type Config struct {
	BadgerOptions badger.Options
	WriteMessageChannel chan []byte
	WriterConfig kafka.WriterConfig
	ReaderConfig kafka.ReaderConfig
}
