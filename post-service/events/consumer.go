package events

import (
	"context"
	"fmt"

	conf "github.com/Hatsker01/Docker_implemintation/post-service/config"
	"github.com/Hatsker01/Docker_implemintation/post-service/events/handler"
	"github.com/Hatsker01/Docker_implemintation/post-service/pkg/logger"
	"github.com/Hatsker01/Docker_implemintation/post-service/storage"
	"github.com/jmoiron/sqlx"
	kafka "github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	kafkaConsumer *kafka.Reader
	eventHandler  *handler.EventHandler
	log           logger.Logger
}
type KafkaConsumera interface {
	Consume(ctx context.Context, topic string)
}

// func NewKafkaConsumer(ctx context.Context, conf config.Config, log logger.Logger, topic string)  {
// 	connString := fmt.Sprintf("%s:%d", conf.KafkaHost, conf.KafkaPort)

// 	r := kafka.NewReader(kafka.ReaderConfig{
// 		Brokers: []string{connString},
// 		Topic:   topic,
// 	})
// 	for {
// 		// the `ReadMessage` method blocks until we receive the next event
// 		msg, err := r.ReadMessage(ctx)
// 		if err != nil {
// 			panic("could not read message " + err.Error())
// 		}
// 		// after receiving the message, log its value
// 		fmt.Println("received: ", string(msg.Value))
// 	}

// }

// func (p *KafkaConsumer) Stop() error {
// 	err := p.kafkaReader.Close()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
func NewKafkaConsumer(db *sqlx.DB, conf *conf.Config, log logger.Logger, topic string) *KafkaConsumer {
	connString := fmt.Sprintf("%s:%d", conf.KafkaHost, conf.KafkaPort)
	return &KafkaConsumer{
		kafkaConsumer: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{connString},
			Topic:          topic,
			MinBytes:       10e3,
			MaxBytes:       10e6,
			Partition:      0,
			CommitInterval: 0,
		}),
		eventHandler: handler.NewEventHandlerFunc(*conf, storage.NewStoragePg(db), log),
		log:          log,
	}
}

func (k *KafkaConsumer) Start() {
	fmt.Println(">>> Kafka consumer started. ")
	for {
		m, err := k.kafkaConsumer.ReadMessage(context.Background())
		if err != nil {
			k.log.Error("failed on consuming a message:", logger.Error(err))
			break
		}

		err = k.eventHandler.Handler(m.Value)
		if err != nil {
			k.log.Error("failed to handle consumed topic:",
				logger.String("on topic", m.Topic), logger.Error(err))
		} else {
			fmt.Println()
			k.log.Info("Successfuly consumed message",
				logger.String("on topic", m.Topic),
				logger.String("message", "success"))
			fmt.Println()
		}

	}
	err := k.kafkaConsumer.Close()
	if err != nil {
		k.log.Error("Error in closing kafka", logger.Error(err))
	}
}

// func (p *KafkaConsumer) Consume(ctx context.Context, topic string) {
// 	// initialize a new reader with the brokers and topic
// 	// the groupID identifies the consumer and prevents
// 	// it from receiving duplicate messages
// 	connString := fmt.Sprintf("%s:%d", conf.Load().KafkaHost, conf.Load().KafkaPort)
// 	r := kafka.NewReader(kafka.ReaderConfig{

// 		Brokers: []string{connString},
// 		Topic:   topic,
// 	})
// 	for {
// 		// the `ReadMessage` method blocks until we receive the next event
// 		msg, err := r.ReadMessage(ctx)
// 		if err != nil {
// 			panic("could not read message " + err.Error())
// 		}
// 		// after receiving the message, log its value
// 		fmt.Println("received: ", string(msg.Value))
// 	}
// }
