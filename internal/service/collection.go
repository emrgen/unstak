package service

import (
	"context"
	authx "github.com/emrgen/authbase/x"
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/emrgen/unpost/internal/model"
	"github.com/emrgen/unpost/internal/store"
	"github.com/google/uuid"
)

// NewCollectionService creates a new collection service
func NewCollectionService(store store.UnPostStore) *CollectionService {
	return &CollectionService{
		store: store,
	}
}

var _ v1.CollectionServiceServer = new(CollectionService)

// CollectionService is the service that provides collection operations
type CollectionService struct {
	store store.UnPostStore
	v1.UnimplementedCollectionServiceServer
}

// CreateCollection creates a collection
func (c *CollectionService) CreateCollection(ctx context.Context, request *v1.CreateCollectionRequest) (*v1.CreateCollectionResponse, error) {
	collection := &model.Collection{
		ID:   uuid.New().String(),
		Name: request.GetName(),
	}

	err := c.store.CreateCollection(ctx, collection)
	if err != nil {
		return nil, err
	}

	return &v1.CreateCollectionResponse{
		Collection: &v1.Collection{
			Id:   collection.ID,
			Name: request.GetName(),
		},
	}, nil
}

// GetCollection gets a collection
func (c *CollectionService) GetCollection(ctx context.Context, request *v1.GetCollectionRequest) (*v1.GetCollectionResponse, error) {
	collection, err := c.store.GetCollection(ctx, uuid.MustParse(request.GetId()))
	if err != nil {
		return nil, err
	}

	return &v1.GetCollectionResponse{
		Collection: &v1.Collection{
			Id:   collection.ID,
			Name: collection.Name,
		},
	}, nil

}

// ListCollection lists collections
func (c *CollectionService) ListCollection(ctx context.Context, request *v1.ListCollectionRequest) (*v1.ListCollectionResponse, error) {
	userID, err := authx.GetAuthbaseUserID(ctx)
	if err != nil {
		return nil, err
	}

	collections, err := c.store.ListCollectionsByOwnerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var responseCollections []*v1.Collection
	for _, collection := range collections {
		responseCollections = append(responseCollections, &v1.Collection{
			Id:   collection.ID,
			Name: collection.Name,
		})
	}

	return &v1.ListCollectionResponse{
		Collections: responseCollections,
	}, nil
}

// UpdateCollection updates a collection
func (c *CollectionService) UpdateCollection(ctx context.Context, request *v1.UpdateCollectionRequest) (*v1.UpdateCollectionResponse, error) {
	collection := &model.Collection{
		ID:   uuid.MustParse(request.GetId()).String(),
		Name: request.GetName(),
	}

	err := c.store.UpdateCollection(ctx, collection)
	if err != nil {
		return nil, err
	}

	return &v1.UpdateCollectionResponse{
		Collection: &v1.Collection{
			Id: collection.ID,
		},
	}, nil
}

// DeleteCollection deletes a collection
func (c *CollectionService) DeleteCollection(ctx context.Context, request *v1.DeleteCollectionRequest) (*v1.DeleteCollectionResponse, error) {
	err := c.store.DeleteCollection(ctx, uuid.MustParse(request.GetId()))
	if err != nil {
		return nil, err
	}

	return &v1.DeleteCollectionResponse{}, nil
}

// AddCollectionTag adds a tag to a collection
func (c *CollectionService) AddCollectionTag(ctx context.Context, request *v1.AddCollectionTagRequest) (*v1.AddCollectionTagResponse, error) {
	collectionID := uuid.MustParse(request.GetCollectionId())
	tagID := uuid.MustParse(request.GetTagId())
	err := c.store.Transaction(ctx, func(ctx context.Context, tx store.UnPostStore) error {
		_, err := tx.GetCollection(ctx, collectionID)
		if err != nil {
			return err
		}

		_, err = tx.GetTag(ctx, tagID)
		if err != nil {
			return err
		}

		collectionTag := &model.CollectionTag{
			CollectionID: collectionID.String(),
			TagID:        tagID.String(),
		}
		err = tx.AddCollectionTag(ctx, collectionTag)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.AddCollectionTagResponse{}, nil
}

// RemoveCollectionTag removes a tag from a collection
func (c *CollectionService) RemoveCollectionTag(ctx context.Context, request *v1.RemoveCollectionTagRequest) (*v1.RemoveCollectionTagResponse, error) {
	collectionID := uuid.MustParse(request.GetCollectionId())
	tagID := uuid.MustParse(request.GetTagId())
	err := c.store.Transaction(ctx, func(ctx context.Context, tx store.UnPostStore) error {
		_, err := tx.GetCollection(ctx, collectionID)
		if err != nil {
			return err
		}

		_, err = tx.GetTag(ctx, tagID)
		if err != nil {
			return err
		}

		collectionTag := &model.CollectionTag{
			CollectionID: collectionID.String(),
			TagID:        tagID.String(),
		}
		err = tx.RemoveCollectionTag(ctx, collectionTag)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.RemoveCollectionTagResponse{}, nil
}
