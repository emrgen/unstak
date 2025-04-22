package store

import (
	"context"
	"github.com/google/uuid"

	"github.com/emrgen/unpost/internal/model"
)

type UnstakStore interface {
	TierStore
	PostStore
	TierMemberStore
	CollectionStore
	CourseStore
	PageStore
	TagStore
	PlatformTagStore
	SpaceStore
	SpaceMemberStore
	Transaction(ctx context.Context, f func(ctx context.Context, store UnstakStore) error) error
	Migrate() error
}

type TierStore interface {
	// CreateTier creates a new space.
	CreateTier(ctx context.Context, space *model.Tier) error
	// GetTier retrieves a space by ID.
	GetTier(ctx context.Context, spaceID uuid.UUID) (*model.Tier, error)
	// ListTiers retrieves a list of spaces by user ID.
	ListTiers(ctx context.Context, userID uuid.UUID) ([]*model.Tier, error)
	// UpdateTier updates a space.
	UpdateTier(ctx context.Context, space *model.Tier) error
	// DeleteTier deletes a space by ID.
	DeleteTier(ctx context.Context, spaceID uuid.UUID) error
	// GetDefaultTier retrieves the default space of a user.
	GetDefaultTier(ctx context.Context, userID uuid.UUID) (*model.Tier, error)
}

type TierMemberStore interface {
	// AddTierMember creates a new member.
	AddTierMember(ctx context.Context, member *model.TierMember) error
	// GetTierMember retrieves a member by ID.
	GetTierMember(ctx context.Context, subMemberID uuid.UUID) (*model.TierMember, error)
	// ListTierMembers retrieves a list of members by space ID.
	ListTierMembers(ctx context.Context, subID uuid.UUID) ([]*model.TierMember, error)
	// UpdateTierMember updates a member.
	UpdateTierMember(ctx context.Context, member *model.TierMember) error
	// RemoveTierMember deletes a member by ID.
	RemoveTierMember(ctx context.Context, subMemberID uuid.UUID) error
}

type PostStore interface {
	// CreatePost creates a new post.
	CreatePost(ctx context.Context, doc *model.Post) error
	// GetPost retrieves a post by ID.
	GetPost(ctx context.Context, id uuid.UUID) (*model.Post, error)
	// ListPostBySpace retrieves a list of tinyposts by space ID.
	ListPostBySpace(ctx context.Context, spaceID uuid.UUID, status *model.PostStatus) ([]*model.Post, error)
	// ListPostByOwnerID retrieves a list of tinyposts by owner ID.
	ListPostByOwnerID(ctx context.Context, userID uuid.UUID, status *model.PostStatus) ([]*model.Post, error)
	// ListPostByUserID retrieves a list of tinyposts by user ID.
	ListPostByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Post, error)
	// ListPostsByTierID retrieves a list of tinyposts by space ID.
	ListPostsByTierID(ctx context.Context, spaceID uuid.UUID) ([]*model.Post, error)
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

type CourseStore interface {
	// CreateCourse creates a new course.
	CreateCourse(ctx context.Context, course *model.Course) error
	// GetCourse retrieves a course by ID.
	GetCourse(ctx context.Context, id uuid.UUID) (*model.Course, error)
	// ListCourses retrieves a list of courses by space ID.
	ListCourses(ctx context.Context, spaceID uuid.UUID) ([]*model.Course, error)
	// UpdateCourse updates a course.
	UpdateCourse(ctx context.Context, course *model.Course) error
	// DeleteCourse deletes a course by ID.
	DeleteCourse(ctx context.Context, id uuid.UUID) error
	// UpdateCourseTags updates the tags of a course.
	UpdateCourseTags(ctx context.Context, courseID uuid.UUID, tags []*model.Tag) error
}

type PageStore interface {
	// CreatePage creates a new page.
	CreatePage(ctx context.Context, page *model.Page) error
	// GetPage retrieves a page by ID.
	GetPage(ctx context.Context, id uuid.UUID) (*model.Page, error)
	// UpdatePage updates a page.
	UpdatePage(ctx context.Context, page *model.Page) error
	// DeletePage deletes a page by ID.
	DeletePage(ctx context.Context, id uuid.UUID) error
	// UpdatePageTags updates the tags of a page.
	UpdatePageTags(ctx context.Context, pageID uuid.UUID, tags []*model.Tag) error
}

type TagStore interface {
	// CreateTag creates a new tag.
	CreateTag(ctx context.Context, tag *model.Tag) error
	// GetTag retrieves a tag by ID.
	GetTag(ctx context.Context, id uuid.UUID) (*model.Tag, error)
	// ListTags retrieves a list of tags by space ID.
	ListTags(ctx context.Context, spaceID uuid.UUID, pageNumber, pageSize uint64) ([]*model.Tag, error)
	// UpdateTag updates a tag.
	UpdateTag(ctx context.Context, tag *model.Tag) error
	// DeleteTag deletes a tag by ID.
	DeleteTag(ctx context.Context, id uuid.UUID) error
}

type PlatformTagStore interface {
	// CreatePlatformTag creates a new platform tag.
	CreatePlatformTag(ctx context.Context, tag *model.PlatformTag) error
	// GetPlatformTag retrieves a platform tag by ID.
	GetPlatformTag(ctx context.Context, id uuid.UUID) (*model.PlatformTag, error)
	// ListPlatformTags retrieves a list of platform tags by space ID.
	ListPlatformTags(ctx context.Context, pageNumber, pageSize uint64) ([]*model.PlatformTag, error)
	// UpdatePlatformTag updates a platform tag.
	UpdatePlatformTag(ctx context.Context, tag *model.PlatformTag) error
	// DeletePlatformTag deletes a platform tag by ID.
	DeletePlatformTag(ctx context.Context, id uuid.UUID) error
}

type SpaceStore interface {
	// CreateSpace creates a new space.
	CreateSpace(ctx context.Context, space *model.Space) error
	// GetMasterSpace retrieves the master space.
	GetMasterSpace(ctx context.Context) (*model.Space, error)
	// GetSpaceByName retrieves a space by name.
	GetSpaceByName(ctx context.Context, spaceName string) (*model.Space, error)
	// GetSpace retrieves a space by ID.
	GetSpace(ctx context.Context, spaceID uuid.UUID) (*model.Space, error)
	// ListSpaces retrieves a list of spaces by user ID.
	ListSpaces(ctx context.Context, userID uuid.UUID) ([]*model.Space, error)
	// UpdateSpace updates a space.
	UpdateSpace(ctx context.Context, space *model.Space) error
	// DeleteSpace deletes a space by ID.
	DeleteSpace(ctx context.Context, spaceID uuid.UUID) error
}

type SpaceMemberStore interface {
	// AddSpaceMember creates a new member.
	AddSpaceMember(ctx context.Context, member *model.SpaceMember) error
	// GetSpaceMember retrieves a member by ID.
	GetSpaceMember(ctx context.Context, spaceID, memberID uuid.UUID) (*model.SpaceMember, error)
	// ListSpaceMembers retrieves a list of members by space ID.
	ListSpaceMembers(ctx context.Context, subID uuid.UUID) ([]*model.SpaceMember, error)
	// UpdateSpaceMember updates a member, if exists it updates the member else creates a new one.
	UpdateSpaceMember(ctx context.Context, member *model.SpaceMember) error
	// RemoveSpaceMember deletes a member by ID.
	RemoveSpaceMember(ctx context.Context, spaceID, memberID uuid.UUID) error
}
