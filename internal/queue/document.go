package queue

import (
	"context"

	"github.com/emrgen/unpost/internal/model"
)

var DocumentUpdateCacheQueue = "unpost:update:queue"
var DocumentUpdateDatabaseQueue = "unpost:sync:queue"

type DocumentQueue interface {
	// PublishChange appends a unpost change to the queue.
	PublishChange(ctx context.Context, change *model.Document) error
	SubscribeUpdateCacheQueue(ctx context.Context) (<-chan *model.Document, error)
	SubscribeUpdateDatabaseQueue(ctx context.Context) (<-chan *model.Document, error)
}
