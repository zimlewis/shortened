package global

import (
	"github.com/segmentio/kafka-go"
	"github.com/zimlewis/shortened/repository"
)

type Config struct {
	WriteMessageChannel chan []byte
	WriterConfig kafka.WriterConfig
	ReaderConfig kafka.ReaderConfig
	Repository   repository.Repository
}
