package queue

import (
	"context"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/emrgen/unpost/internal/model"
)

type Kafka struct {
	producer *kafka.Producer
	consumer *kafka.Consumer
	stop     chan struct{}
}

// PublishChange implements Document.
func (k *Kafka) PublishChange(ctx context.Context, change *model.Document) error {
	var err error
	changeJSON, err := json.Marshal(change)
	if err != nil {
		return err
	}

	err = k.producer.InitTransactions(ctx)
	if err != nil {
		return err
	}
	tx := k.producer.BeginTransaction()
	if tx == nil {
		return err
	}
	err = k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &DocumentUpdateCacheQueue, Partition: kafka.PartitionAny},
		Value:          changeJSON,
		Key:            []byte(change.ID),
	}, nil)
	if err != nil {
		err := k.producer.AbortTransaction(ctx)
		if err != nil {
			return err
		}
		return err
	}

	err = k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &DocumentUpdateDatabaseQueue, Partition: kafka.PartitionAny},
		Value:          changeJSON,
		Key:            []byte(change.ID),
	}, nil)
	if err != nil {
		err := k.producer.AbortTransaction(ctx)
		if err != nil {
			return err
		}
		return err
	}

	err = k.producer.CommitTransaction(ctx)
	if err != nil {
		return err
	}

	return nil
}

// SubscribeUpdateCacheQueue implements DocumentQueue.
func (k *Kafka) SubscribeUpdateCacheQueue(ctx context.Context) (<-chan *model.Document, error) {
	return k.subscribe(ctx, DocumentUpdateCacheQueue)
}

// SubscribeUpdateDatabaseQueue implements DocumentQueue.
func (k *Kafka) SubscribeUpdateDatabaseQueue(ctx context.Context) (<-chan *model.Document, error) {
	return k.subscribe(ctx, DocumentUpdateDatabaseQueue)
}

func (k *Kafka) subscribe(_ context.Context, topic string) (<-chan *model.Document, error) {
	err := k.consumer.Subscribe(topic, nil)
	if err != nil {
		return nil, err
	}

	update := make(chan *model.Document)

	go func() {
		for {
			select {
			case <-k.stop:
				close(update)
				return
			default:
				ev := k.consumer.Poll(100)
				if ev == nil {
					continue
				}

				switch e := ev.(type) {
				case *kafka.Message:
					var doc model.Document
					if err := json.Unmarshal(e.Value, &doc); err != nil {
						continue
					}

					update <- &doc
				}
			}

		}
	}()

	return update, nil
}
