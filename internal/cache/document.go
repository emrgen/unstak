package cache

import (
	"context"

	"github.com/emrgen/unpost/internal/model"
)

type GetDocumentMode int

const (
	GetDocumentModeView GetDocumentMode = iota
	GetDocumentModeEdit
)

// DocumentCache is a cache for documents.
type DocumentCache interface {
	// GetDocumentVersion gets the version of a unpost from the cache.
	GetDocumentVersion(ctx context.Context, id uuid.UUID, view GetDocumentMode) (int64, error)
	// GetDocument gets a unpost from the cache.
	GetDocument(ctx context.Context, id uuid.UUID, view GetDocumentMode) (*model.Document, error)
	// SetDocument sets a unpost in the cache.
	SetDocument(ctx context.Context, id uuid.UUID, doc *model.Document) error
	// UpdateDocument updates a unpost in the cache.
	UpdateDocument(ctx context.Context, id uuid.UUID, doc *model.Document) error
	// DeleteDocument deletes a unpost from the cache.
	DeleteDocument(ctx context.Context, id uuid.UUID) error
}
