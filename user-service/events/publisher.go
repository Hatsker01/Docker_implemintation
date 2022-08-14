package events

import (
	"context"
	"fmt"
	"time"

	"github.com/Hatsker01/Docker_implemintation/user-service/config"
	"github.com/Hatsker01/Docker_implemintation/user-service/pkg/logger"
	broker "github.com/Hatsker01/Docker_implemintation/user-service/pkg/messagebroker"
	"github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
	kafkaWriter *kafka.Writer
	log         logger.Logger
}

func NewKafkaProducerBroker(conf config.Config, log logger.Logger, topic string) broker.Publisher {
	connString := fmt.Sprintf("%s:%d", conf.KafkaHost, conf.KafkaPort)

	return &KafkaPublisher{
		kafkaWriter: &kafka.Writer{
			Addr:         kafka.TCP(connString),
			Topic:        topic,
			BatchTimeout: 10 * time.Millisecond,
			AllowAutoTopicCreation: true,
		},
		log: log,
	}
}

func (p *KafkaPublisher) Start() error {
	return nil
}
func (p *KafkaPublisher) Stop() error {
	err := p.kafkaWriter.Close()
	if err != nil {
		return err
	}
	return nil
}
func (p *KafkaPublisher) Publish(key, body []byte, logBody string) error {
	message := kafka.Message{
		Key:   key,
		Value: body,
	}

	if err := p.kafkaWriter.WriteMessages(context.Background(), message); err != nil {
		return err
	}
	p.log.Info("Message published(key/body): " + string(key) + "/" + logBody)
	return nil
}
