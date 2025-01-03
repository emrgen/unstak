package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/emrgen/unpost/internal/model"
	redis "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var _ DocumentQueue = (*Redis)(nil)

type Redis struct {
	client *redis.Client
}

// PublishChange implements Document.
func (r *Redis) PublishChange(ctx context.Context, change *model.Document) error {
	doc, err := json.Marshal(change)
	if err != nil {
		return err
	}

	_, err = r.client.TxPipelined(ctx, func(tx redis.Pipeliner) error {
		if err := tx.RPush(ctx, DocumentUpdateCacheQueue, doc).Err(); err != nil {
			return err
		}

		if err := tx.RPush(ctx, DocumentUpdateDatabaseQueue, doc).Err(); err != nil {
			return err
		}

		return nil
	})

	return err
}

// SubscribeUpdateCacheQueue implements DocumentQueue.
func (r *Redis) SubscribeUpdateCacheQueue(ctx context.Context) (<-chan *model.Document, error) {
	return r.subscribeToQueue(ctx, DocumentUpdateCacheQueue)
}

// SubscribeUpdateDatabaseQueue implements DocumentQueue.
func (r *Redis) SubscribeUpdateDatabaseQueue(ctx context.Context) (<-chan *model.Document, error) {
	return r.subscribeToQueue(ctx, DocumentUpdateDatabaseQueue)
}

func (r *Redis) subscribeToQueue(ctx context.Context, queue string) (<-chan *model.Document, error) {
	updateChan := make(chan *model.Document)

	go func() {
		for {
			updates := r.client.LRange(ctx, queue, 0, 1)
			if updates.Err() != nil {
				logrus.Errorf("going to sleep for 1sec, failed to get updates: %v", updates.Err())
				time.Sleep(1 * time.Microsecond)
				continue
			}

			for _, update := range updates.Val() {
				var doc model.Document
				if err := json.Unmarshal([]byte(update), &doc); err != nil {
					logrus.Errorf("failed to unmarshal update: %v", err)
				}

				updateChan <- &doc
			}

			time.Sleep(time.Millisecond)
		}
	}()

	return updateChan, nil
}
