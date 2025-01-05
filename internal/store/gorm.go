package store

import (
	"context"
	"github.com/emrgen/unpost/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// NewGormStore creates a new GormStore.
func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{
		db: db,
	}
}

var (
	_ UnPostStore = (*GormStore)(nil)
)

type GormStore struct {
	db *gorm.DB
}

func (g *GormStore) CreateUser(ctx context.Context, user *model.User) error {
	return g.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(user).Error
}

func (g *GormStore) GetUser(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	if err := g.db.Where("id = ?", id.String()).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (g *GormStore) CreatePlatformTag(ctx context.Context, tag *model.PlatformTag) error {
	//TODO implement me
	panic("implement me")
}

func (g *GormStore) GetPlatformTag(ctx context.Context, id uuid.UUID) (*model.PlatformTag, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GormStore) ListPlatformTags(ctx context.Context, pageNumber, pageSize uint64) ([]*model.PlatformTag, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GormStore) UpdatePlatformTag(ctx context.Context, tag *model.PlatformTag) error {
	//TODO implement me
	panic("implement me")
}

func (g *GormStore) DeletePlatformTag(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (g *GormStore) CreateSpace(ctx context.Context, space *model.Space) error {
	return g.db.Create(space).Error
}

func (g *GormStore) GetSpace(ctx context.Context, spaceID uuid.UUID) (*model.Space, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GormStore) ListSpaces(ctx context.Context, userID uuid.UUID) ([]*model.Space, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GormStore) UpdateSpace(ctx context.Context, space *model.Space) error {
	//TODO implement me
	panic("implement me")
}

func (g *GormStore) DeleteSpace(ctx context.Context, spaceID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (g *GormStore) AddSpaceMember(ctx context.Context, member *model.SpaceMember) error {
	return g.db.Create(member).Error
}

func (g *GormStore) GetSpaceMember(ctx context.Context, subMemberID uuid.UUID) (*model.SpaceMember, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GormStore) ListSpaceMembers(ctx context.Context, subID uuid.UUID) ([]*model.SpaceMember, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GormStore) UpdateSpaceMember(ctx context.Context, member *model.SpaceMember) error {
	//TODO implement me
	panic("implement me")
}

func (g *GormStore) RemoveSpaceMember(ctx context.Context, subMemberID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (g *GormStore) CreateCourse(ctx context.Context, course *model.Course) error {
	return g.db.Create(course).Error
}

func (g *GormStore) GetCourse(ctx context.Context, id uuid.UUID) (*model.Course, error) {
	var course model.Course
	if err := g.db.Where("id = ?", id.String()).First(&course).Error; err != nil {
		return nil, err
	}

	return &course, nil
}

func (g *GormStore) ListCourses(ctx context.Context, spaceID uuid.UUID) ([]*model.Course, error) {
	var courses []*model.Course
	if err := g.db.Where("space_id = ?", spaceID.String()).Find(&courses).Error; err != nil {
		return nil, err
	}

	return courses, nil
}

func (g *GormStore) UpdateCourse(ctx context.Context, course *model.Course) error {
	return g.db.Save(course).Error
}

func (g *GormStore) DeleteCourse(ctx context.Context, id uuid.UUID) error {
	course := &model.Course{
		ID: id.String(),
	}
	return g.db.Delete(course).Error
}

func (g *GormStore) UpdateCourseTags(ctx context.Context, courseID uuid.UUID, tags []*model.Tag) error {
	return g.db.Model(&model.Course{ID: courseID.String()}).Association("Tags").Replace(tags)
}

func (g *GormStore) CreatePage(ctx context.Context, page *model.Page) error {
	return g.db.Create(page).Error
}

func (g *GormStore) GetPage(ctx context.Context, id uuid.UUID) (*model.Page, error) {
	var page model.Page
	if err := g.db.Where("id = ?", id.String()).First(&page).Error; err != nil {
		return nil, err
	}

	return &page, nil
}

func (g *GormStore) UpdatePage(ctx context.Context, page *model.Page) error {
	return g.db.Save(page).Error
}

func (g *GormStore) DeletePage(ctx context.Context, id uuid.UUID) error {
	page := &model.Page{
		ID: id.String(),
	}
	return g.db.Delete(page).Error
}

func (g *GormStore) UpdatePageTags(ctx context.Context, pageID uuid.UUID, tags []*model.Tag) error {
	return g.db.Model(&model.Page{ID: pageID.String()}).Association("Tags").Replace(tags)
}

// -----------------------
// TierStore
// -----------------------

func (g *GormStore) CreatePost(ctx context.Context, post *model.Post) error {
	return g.db.Create(post).Error
}

func (g *GormStore) GetPost(ctx context.Context, id uuid.UUID) (*model.Post, error) {
	var post model.Post
	if err := g.db.Where("id = ?", id.String()).Preload("Tags").First(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

func (g *GormStore) ListPostByOwnerID(ctx context.Context, userID uuid.UUID, status *model.PostStatus) ([]*model.Post, error) {
	var posts []*model.Post
	if status != nil {
		if err := g.db.Where("created_by_id = ? AND status = ?", userID.String(), status).Find(&posts).Error; err != nil {
			return nil, err
		}

		return posts, nil
	}

	if err := g.db.Where("created_by_id = ?", userID.String()).Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func (g *GormStore) UpdatePostTags(ctx context.Context, postID uuid.UUID, tags []*model.Tag) error {
	return g.db.Model(&model.Post{ID: postID.String()}).Association("Tags").Replace(tags)
}

// ListPostByUserID retrieves a list of tinyposts by user ID.
// returns a list of tinyposts the user has access to.
func (g *GormStore) ListPostByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Post, error) {
	var posts []*model.Post
	if err := g.db.Where("created_by_id = ?", userID.String()).Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

// ListPostsByTierID retrieves a list of tinyposts by space ID.
func (g *GormStore) ListPostsByTierID(ctx context.Context, spaceID uuid.UUID) ([]*model.Post, error) {
	var posts []*model.Post
	if err := g.db.Where("space_id = ?", spaceID.String()).Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func (g *GormStore) UpdatePost(ctx context.Context, post *model.Post) error {
	return g.db.Save(post).Error
}

func (g *GormStore) DeletePost(ctx context.Context, id uuid.UUID) error {
	post := &model.Post{
		ID: id.String(),
	}
	return g.db.Delete(post).Error
}

func (g *GormStore) UpdatePostReaction(ctx context.Context, userID, postID uuid.UUID, reaction *model.Reaction) error {
	return g.db.Create(reaction).Error
}

func (g *GormStore) AddMember(ctx context.Context, spaceID, userID uuid.UUID, permission uint64) error {
	member := &model.TierMember{
		TierID: spaceID.String(),
		UserID: userID.String(),
	}

	return g.db.Create(member).Error
}

func (g *GormStore) GetMember(ctx context.Context, spaceID, userID uuid.UUID) (*model.TierMember, error) {
	var member model.TierMember
	if err := g.db.Where("space_id = ? AND user_id = ?", spaceID.String(), userID.String()).First(&member).Error; err != nil {
		return nil, err
	}

	return &member, nil

}

func (g *GormStore) ListMembers(ctx context.Context, spaceID uuid.UUID) ([]*uuid.UUID, error) {
	var members []*model.TierMember
	if err := g.db.Where("space_id = ?", spaceID.String()).Find(&members).Error; err != nil {
		return nil, err
	}

	var memberIDs []*uuid.UUID
	for _, member := range members {
		id, err := uuid.Parse(member.UserID)
		if err != nil {
			return nil, err
		}
		memberIDs = append(memberIDs, &id)
	}

	return memberIDs, nil

}

func (g *GormStore) UpdateMember(ctx context.Context, member *model.TierMember) error {
	return g.db.Save(member).Error
}

func (g *GormStore) RemoveMember(ctx context.Context, spaceID, userID uuid.UUID) error {
	member := &model.TierMember{
		TierID: spaceID.String(),
		UserID: userID.String(),
	}
	return g.db.Delete(member).Error
}

// -----------------------
// CollectionStore
// -----------------------

func (g *GormStore) CreateCollection(ctx context.Context, collection *model.Collection) error {
	return g.db.Create(collection).Error
}

func (g *GormStore) GetCollection(ctx context.Context, id uuid.UUID) (*model.Collection, error) {
	var collection model.Collection
	if err := g.db.Where("id = ?", id.String()).First(&collection).Error; err != nil {
		return nil, err
	}

	return &collection, nil
}

func (g *GormStore) ListCollectionsByOwnerID(ctx context.Context, userID uuid.UUID) ([]*model.Collection, error) {
	var collections []*model.Collection
	if err := g.db.Where("created_by_id = ?", userID.String()).Find(&collections).Error; err != nil {
		return nil, err
	}

	return collections, nil
}

func (g *GormStore) ListCollectionsByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Collection, error) {
	var collections []*model.Collection
	if err := g.db.Where("created_by_id = ?", userID.String()).Find(&collections).Error; err != nil {
		return nil, err
	}

	return collections, nil
}

func (g *GormStore) UpdateCollection(ctx context.Context, collection *model.Collection) error {
	return g.db.Save(collection).Error
}

func (g *GormStore) DeleteCollection(ctx context.Context, id uuid.UUID) error {
	collection := &model.Collection{
		ID: id.String(),
	}
	return g.db.Delete(collection).Error
}

func (g *GormStore) AddCollectionTag(ctx context.Context, tag *model.CollectionTag) error {
	return g.db.Create(tag).Error
}

func (g *GormStore) RemoveCollectionTag(ctx context.Context, tag *model.CollectionTag) error {
	return g.db.Delete(tag).Error
}

// -----------------------
// TagStore
// -----------------------

func (g *GormStore) CreateTag(ctx context.Context, tag *model.Tag) error {
	return g.db.Create(tag).Error
}

func (g *GormStore) GetTag(ctx context.Context, id uuid.UUID) (*model.Tag, error) {
	var tag model.Tag
	if err := g.db.Where("id = ?", id.String()).First(&tag).Error; err != nil {
		return nil, err
	}

	return &tag, nil
}

func (g *GormStore) ListTags(ctx context.Context, spaceID uuid.UUID, pageNumber, pageSize uint64) ([]*model.Tag, error) {
	var tags []*model.Tag
	if err := g.db.Where("space_id = ?", spaceID.String()).Limit(int(pageSize)).Offset(int(pageNumber * pageSize)).Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

func (g *GormStore) UpdateTag(ctx context.Context, tag *model.Tag) error {
	return g.db.Save(tag).Error
}

func (g *GormStore) DeleteTag(ctx context.Context, id uuid.UUID) error {
	return g.db.Delete(&model.Tag{ID: id.String()}).Error
}

func (g *GormStore) UpdateTierMember(ctx context.Context, member *model.TierMember) error {
	//TODO implement me
	panic("implement me")
}

func (g *GormStore) CreateTier(ctx context.Context, space *model.Tier) error {
	return g.db.Create(space).Error
}

func (g *GormStore) GetTier(ctx context.Context, id uuid.UUID) (*model.Tier, error) {
	var space model.Tier
	if err := g.db.Where("id = ?", id.String()).First(&space).Error; err != nil {
		return nil, err
	}

	return &space, nil
}

func (g *GormStore) ListTiers(ctx context.Context, userID uuid.UUID) ([]*model.Tier, error) {
	var spaces []*model.Tier
	if err := g.db.Where("created_by_id = ?", userID.String()).Find(&spaces).Error; err != nil {
		return nil, err
	}

	return spaces, nil
}

func (g *GormStore) UpdateTier(ctx context.Context, space *model.Tier) error {
	return g.db.Save(space).Error
}

func (g *GormStore) DeleteTier(ctx context.Context, id uuid.UUID) error {
	post := &model.Tier{
		ID: id.String(),
	}
	return g.db.Delete(post).Error
}

func (g *GormStore) GetDefaultTier(ctx context.Context, userID uuid.UUID) (*model.Tier, error) {
	var space model.Tier
	if err := g.db.Where("created_by_id = ? AND user_default = true", userID.String()).First(&space).Error; err != nil {
		return nil, err
	}

	return &space, nil
}

func (g *GormStore) AddTierMember(ctx context.Context, member *model.TierMember) error {
	return g.db.Create(member).Error
}

func (g *GormStore) GetTierMember(ctx context.Context, subMemberID uuid.UUID) (*model.TierMember, error) {
	var member model.TierMember
	if err := g.db.Where("id = ?", subMemberID.String()).Preload("Tier").First(&member).Error; err != nil {
		return nil, err
	}

	return &member, nil

}

func (g *GormStore) ListTierMembers(ctx context.Context, subID uuid.UUID) ([]*model.TierMember, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GormStore) RemoveTierMember(ctx context.Context, subMemberID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (g *GormStore) Transaction(ctx context.Context, f func(ctx context.Context, store UnPostStore) error) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		return f(ctx, NewGormStore(tx))
	})
}

func (g *GormStore) Migrate() error {
	return model.Migrate(g.db)
}
