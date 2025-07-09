package service

import (
	"context"
	"github.com/emrgen/authbase"
	authx "github.com/emrgen/authbase/x"
	docv1 "github.com/emrgen/document/apis/v1"
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/emrgen/unpost/internal/model"
	"github.com/emrgen/unpost/internal/store"
	"github.com/emrgen/unpost/internal/x"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NewPostService creates a new post service
func NewPostService(cfg *authx.AuthbaseConfig, store store.UnstakStore) *PostService {
	return &PostService{
		cfg:   cfg,
		store: store,
	}
}

var _ v1.PostServiceServer = new(PostService)

// PostService is the service that provides post operations
type PostService struct {
	cfg        *authx.AuthbaseConfig
	store      store.UnstakStore
	docClient  docv1.DocumentServiceClient
	authClient authbase.Client
	v1.UnimplementedPostServiceServer
}

func (p *PostService) CreatePost(ctx context.Context, req *v1.CreatePostRequest) (*v1.CreatePostResponse, error) {
	var postID uuid.UUID
	var err error
	if req.PostId != nil {
		postID, err = uuid.Parse(req.GetPostId())
		if err != nil {
			return nil, err
		}
	} else {
		postID = uuid.New()
	}

	post := &model.Post{
		ID:      postID.String(),
		Title:   req.GetTitle(),
		Summary: req.GetSummary(),
		Content: req.GetContent(),
		Slug:    req.GetSlug() + "-" + x.RandomString(10),
		Status:  model.PostStatusDraft,
		Tags:    nil,
		Version: 0,
	}

	err = p.store.CreatePost(ctx, post)
	if err != nil {
		return nil, err
	}

	return &v1.CreatePostResponse{
		Post: &v1.Post{
			Id:        postID.String(),
			Title:     post.Title,
			CreatedAt: timestamppb.New(post.CreatedAt),
			UpdatedAt: timestamppb.New(post.UpdatedAt),
			Version:   post.Version,
		},
	}, nil
}

func (p *PostService) GetPost(ctx context.Context, request *v1.GetPostRequest) (*v1.GetPostResponse, error) {
	post, err := p.store.GetPost(ctx, uuid.MustParse(request.GetId()))
	if err != nil {
		return nil, err
	}

	postProto := &v1.Post{
		Id:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		Tags:    make([]*v1.Tag, 0),
		Version: 1,
		Status:  postStatusToProto(post.Status),
	}

	for _, tag := range post.Tags {
		postProto.Tags = append(postProto.Tags, &v1.Tag{
			Id:   tag.ID,
			Name: tag.Name,
		})
	}

	return &v1.GetPostResponse{
		Post: postProto,
	}, nil
}

// ListPost retrieves a list of posts within a space
func (p *PostService) ListPost(ctx context.Context, request *v1.ListPostRequest) (*v1.ListPostResponse, error) {
	posts, err := p.store.ListPosts(ctx, &store.PostFiler{})
	if err != nil {
		return nil, err
	}

	postProtos := make([]*v1.Post, 0)
	for _, post := range posts {
		postProto := &v1.Post{
			Id:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			CreatedAt: timestamppb.New(post.CreatedAt),
			UpdatedAt: timestamppb.New(post.UpdatedAt),
		}

		postProtos = append(postProtos, postProto)
	}

	return &v1.ListPostResponse{
		Posts: postProtos,
	}, nil
}

func (p *PostService) UpdatePost(ctx context.Context, req *v1.UpdatePostRequest) (*v1.UpdatePostResponse, error) {
	postID, err := uuid.Parse(req.GetPostId())
	if err != nil {
		return nil, err
	}
	err = p.store.Transaction(ctx, func(ctx context.Context, store store.UnstakStore) error {
		post, err := store.GetPost(ctx, postID)
		if err != nil {
			return err
		}

		if req.Title != nil {
			post.Title = req.GetTitle()
		}

		if req.Content != nil {
			post.Content = req.GetContent()
		}

		return store.UpdatePost(ctx, post)
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdatePostResponse{
		Post: &v1.Post{
			Id: req.GetPostId(),
		},
	}, nil
}

func (p *PostService) DeletePost(ctx context.Context, request *v1.DeletePostRequest) (*v1.DeletePostResponse, error) {
	err := p.store.DeletePost(ctx, uuid.MustParse(request.GetId()))
	if err != nil {
		return nil, err
	}

	return &v1.DeletePostResponse{}, nil
}

func (p *PostService) AddPostTag(ctx context.Context, request *v1.AddPostTagRequest) (*v1.AddPostTagResponse, error) {
	postID := uuid.MustParse(request.GetPostId())
	tagID := uuid.MustParse(request.GetTagId())

	err := p.store.Transaction(ctx, func(ctx context.Context, tx store.UnstakStore) error {
		post, err := tx.GetPost(ctx, postID)
		if err != nil {
			return err
		}

		tag, err := tx.GetTag(ctx, tagID)
		if err != nil {
			return err
		}

		post.Tags = append(post.Tags, tag)

		err = tx.UpdatePost(ctx, post)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.AddPostTagResponse{
		Post: &v1.Post{
			Id: request.GetPostId(),
		},
	}, nil
}

func (p *PostService) RemovePostTag(ctx context.Context, request *v1.RemovePostTagRequest) (*v1.RemovePostTagResponse, error) {
	tags := make([]*v1.Tag, 0)
	err := p.store.Transaction(ctx, func(ctx context.Context, tx store.UnstakStore) error {
		postID := uuid.MustParse(request.GetPostId())
		tagID := uuid.MustParse(request.GetTagId())
		post, err := tx.GetPost(ctx, postID)
		if err != nil {
			return err
		}

		tag, err := tx.GetTag(ctx, tagID)
		if err != nil {
			return err
		}

		for i, t := range post.Tags {
			if t.ID == tag.ID {
				post.Tags = append(post.Tags[:i], post.Tags[i+1:]...)
				logrus.Info("tag removed")
				break
			}
		}

		err = tx.UpdatePostTags(ctx, postID, post.Tags)
		if err != nil {
			return err
		}

		for _, t := range post.Tags {
			tags = append(tags, &v1.Tag{
				Id:   t.ID,
				Name: t.Name,
			})
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.RemovePostTagResponse{
		Post: &v1.Post{
			Id:   request.GetPostId(),
			Tags: tags,
		},
	}, nil
}

// UpdatePostReaction updates the reaction of a post
// A cron job will be responsible for updating the post's aggregated reaction counts
func (p *PostService) UpdatePostReaction(ctx context.Context, request *v1.UpdatePostReactionRequest) (*v1.UpdatePostReactionResponse, error) {
	userID, err := authx.GetAuthbaseAccountID(ctx)
	if err != nil {
		return nil, err
	}
	postID, err := uuid.Parse(request.GetPostId())
	if err != nil {
		return nil, err
	}

	err = p.store.Transaction(ctx, func(ctx context.Context, tx store.UnstakStore) error {
		_, err := tx.GetPost(ctx, uuid.MustParse(request.GetPostId()))
		if err != nil {
			return err
		}

		// Check if the user has already reacted to the post
		reaction := &model.Reaction{
			PostID: postID.String(),
			UserID: userID.String(),
			Name:   request.GetReactionName(),
			State:  request.GetCount(),
		}

		err = p.store.UpdatePostReaction(ctx, userID, postID, reaction)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdatePostReactionResponse{
		Post: &v1.Post{
			Id: request.GetPostId(),
		},
	}, nil
}

func (p *PostService) UpdatePostStatus(ctx context.Context, request *v1.UpdatePostStatusRequest) (*v1.UpdatePostStatusResponse, error) {
	postID := uuid.MustParse(request.GetPostId())
	err := p.store.Transaction(ctx, func(ctx context.Context, tx store.UnstakStore) error {
		post, err := tx.GetPost(ctx, postID)
		if err != nil {
			return err
		}

		post.Status = postStatusFromProto(request.GetStatus())
		err = tx.UpdatePost(ctx, post)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdatePostStatusResponse{
		Post: &v1.Post{
			Id: request.GetPostId(),
		},
	}, nil
}

func postStatusFromProto(status v1.PostStatus) model.PostStatus {
	switch status {
	case v1.PostStatus_PUBLISHED:
		return model.PostStatusPublished
	case v1.PostStatus_ARCHIVED:
		return model.PostStatusArchived
	default:
		return model.PostStatusDraft
	}
}

func postStatusToProto(status model.PostStatus) v1.PostStatus {
	switch status {
	case model.PostStatusPublished:
		return v1.PostStatus_PUBLISHED
	case model.PostStatusArchived:
		return v1.PostStatus_ARCHIVED
	default:
		return v1.PostStatus_DRAFT
	}
}
