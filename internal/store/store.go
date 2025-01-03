package store

import (
	"context"
	"github.com/google/uuid"

	"github.com/emrgen/unpost/internal/model"
)

type TinyPostStore interface {
	OutletStore
	PostStore
	OutletMemberStore
	CollectionStore
	TagStore
	Transaction(ctx context.Context, f func(ctx context.Context, store TinyPostStore) error) error
	Migrate() error
}

type OutletStore interface {
	// CreateOutlet creates a new space.
	CreateOutlet(ctx context.Context, space *model.Subscription) error
	// GetOutlet retrieves a space by ID.
	GetOutlet(ctx context.Context, spaceID uuid.UUID) (*model.Subscription, error)
	// ListOutlets retrieves a list of spaces by user ID.
	ListOutlets(ctx context.Context, userID uuid.UUID) ([]*model.Subscription, error)
	// UpdateOutlet updates a space.
	UpdateOutlet(ctx context.Context, space *model.Subscription) error
	// DeleteOutlet deletes a space by ID.
	DeleteOutlet(ctx context.Context, spaceID uuid.UUID) error
	// GetDefaultOutlet retrieves the default space of a user.
	GetDefaultOutlet(ctx context.Context, userID uuid.UUID) (*model.Subscription, error)
}

type OutletMemberStore interface {
	// AddOutletMember creates a new member.
	AddOutletMember(ctx context.Context, member *model.OutletMember) error
	// GetOutletMember retrieves a member by ID.
	GetOutletMember(ctx context.Context, spaceID, userID uuid.UUID) (*model.OutletMember, error)
	// ListOutletMembers retrieves a list of members by space ID.
	ListOutletMembers(ctx context.Context, spaceID uuid.UUID) ([]*uuid.UUID, error)
	// UpdateOutletMember updates a member.
	UpdateOutletMember(ctx context.Context, member *model.OutletMember) error
	// RemoveOutletMember deletes a member by ID.
	RemoveOutletMember(ctx context.Context, spaceID, userID uuid.UUID) error
}

type PostStore interface {
	// CreatePost creates a new post.
	CreatePost(ctx context.Context, doc *model.Post) error
	// GetPost retrieves a post by ID.
	GetPost(ctx context.Context, id uuid.UUID) (*model.Post, error)
	// ListPostByOwnerID retrieves a list of tinyposts by owner ID.
	ListPostByOwnerID(ctx context.Context, userID uuid.UUID, status *model.PostStatus) ([]*model.Post, error)
	// ListPostByUserID retrieves a list of tinyposts by user ID.
	ListPostByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Post, error)
	// ListPostsByOutletID retrieves a list of tinyposts by space ID.
	ListPostsByOutletID(ctx context.Context, spaceID uuid.UUID) ([]*model.Post, error)
	// UpdatePost updates a post.
	UpdatePost(ctx context.Context, doc *model.Post) error
	// DeletePost deletes a post by ID.
	DeletePost(ctx context.Context, id uuid.UUID) error
	// UpdatePostReaction updates a post reaction.
	UpdatePostReaction(ctx context.Context, userID, postID uuid.UUID, reaction *model.Reaction) error
	// UpdatePostTags updates the tags of a post.
	UpdatePostTags(ctx context.Context, postID uuid.UUID, tags []*model.Tag) error
}

type CollectionStore interface {
	// CreateCollection creates a new collection.
	CreateCollection(ctx context.Context, collection *model.Collection) error
	// GetCollection retrieves a collection by ID.
	GetCollection(ctx context.Context, id uuid.UUID) (*model.Collection, error)
	// ListCollectionsByUserID retrieves a list of collections by user ID.
	ListCollectionsByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Collection, error)
	// ListCollectionsByOwnerID retrieves a list of collections by owner ID.
	ListCollectionsByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*model.Collection, error)
	// UpdateCollection updates a collection.
	UpdateCollection(ctx context.Context, collection *model.Collection) error
	// DeleteCollection deletes a collection by ID.
	DeleteCollection(ctx context.Context, id uuid.UUID) error
	// AddCollectionTag	adds a tag to a collection.
	AddCollectionTag(ctx context.Context, tag *model.CollectionTag) error
	// RemoveCollectionTag removes a tag from a collection.
	RemoveCollectionTag(ctx context.Context, tag *model.CollectionTag) error
}

type TagStore interface {
	// CreateTag creates a new tag.
	CreateTag(ctx context.Context, tag *model.Tag) error
	// GetTag retrieves a tag by ID.
	GetTag(ctx context.Context, id uuid.UUID) (*model.Tag, error)
	// ListTags retrieves a list of tags by space ID.
	ListTags(ctx context.Context, pageNumber, pageSize uint64) ([]*model.Tag, error)
	// UpdateTag updates a tag.
	UpdateTag(ctx context.Context, tag *model.Tag) error
	// DeleteTag deletes a tag by ID.
	DeleteTag(ctx context.Context, id uuid.UUID) error
}
