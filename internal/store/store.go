package store

import (
	"context"
	"github.com/google/uuid"

	"github.com/emrgen/unpost/internal/model"
)

type UnstakStore interface {
	PostStore
	TierStore
	TierMemberStore
	CourseStore
	PageStore
	TagStore
	PlatformTagStore
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

type PostFiler struct {
	TierID  *uuid.UUID
	OwnerID *uuid.UUID
	UserID  *uuid.UUID
	Status  *model.PostStatus
}

type PostStore interface {
	// CreatePost creates a new post.
	CreatePost(ctx context.Context, doc *model.Post) error
	// GetPost retrieves a post by ID.
	GetPost(ctx context.Context, id uuid.UUID) (*model.Post, error)
	// GetPostBySlugID retries the post by slug id.
	GetPostBySlugID(ctx context.Context, id string) (*model.Post, error)
	// ListPosts retrieves a list of tinyposts by space ID.
	ListPosts(ctx context.Context, filer *PostFiler) ([]*model.Post, error)
	// UpdatePost updates a post.
	UpdatePost(ctx context.Context, doc *model.Post) error
	// DeletePost deletes a post by ID.
	DeletePost(ctx context.Context, id uuid.UUID) error
	// UpdatePostReaction updates a post reaction.
	UpdatePostReaction(ctx context.Context, userID, postID uuid.UUID, reaction *model.Reaction) error
	// UpdatePostTags updates the tags of a post.
	UpdatePostTags(ctx context.Context, postID uuid.UUID, tags []*model.Tag) error
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
