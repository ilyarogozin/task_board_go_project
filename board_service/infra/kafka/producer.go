package kafka

import (
	"time"

	k "github.com/segmentio/kafka-go"
)

func NewKafkaProducer(brokers []string, topic string) *k.Writer {
	return &k.Writer{
		Addr:         k.TCP(brokers...),
		Topic:        topic,
		Balancer:     &k.Hash{},
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: k.RequireAll,
	}
}