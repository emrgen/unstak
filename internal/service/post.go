package service

import (
	"context"
	mapset "github.com/deckarep/golang-set/v2"
	authv1 "github.com/emrgen/authbase/apis/v1"
	authx "github.com/emrgen/authbase/x"
	docv1 "github.com/emrgen/document/apis/v1"
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/emrgen/unpost/internal/model"
	"github.com/emrgen/unpost/internal/store"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NewPostService creates a new post service
func NewPostService(cfg *authx.AuthbaseConfig, store store.UnPostStore, docClient docv1.DocumentServiceClient) *PostService {
	return &PostService{
		cfg:       cfg,
		store:     store,
		docClient: docClient,
	}
}

var _ v1.PostServiceServer = new(PostService)

// PostService is the service that provides post operations
type PostService struct {
	cfg       *authx.AuthbaseConfig
	store     store.UnPostStore
	docClient docv1.DocumentServiceClient
	v1.UnimplementedPostServiceServer
}

func (p *PostService) CreatePost(ctx context.Context, request *v1.CreatePostRequest) (*v1.CreatePostResponse, error) {
	var err error

	poolID, err := authx.GetAuthbasePoolID(ctx)
	if err != nil {
		return nil, err
	}
	userID, err := authx.GetAuthbaseAccountID(ctx)
	if err != nil {
		return nil, err
	}

	logrus.Infof("creating post %s", request.GetTitle())

	doc, err := p.docClient.CreateDocument(p.cfg.IntoContext(), &docv1.CreateDocumentRequest{
		ProjectId: poolID.String(),
		Title:     request.GetTitle(),
		Content:   request.GetContent(),
	})
	if err != nil {
		return nil, err
	}

	user, err := p.store.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	post := &model.Post{
		ID:          uuid.New().String(),
		CreatedByID: userID.String(),
		DocumentID:  doc.GetDocument().GetId(),
		Authors:     []*model.User{user},
	}

	// get user default space-id if no space-id is provided
	err = p.store.Transaction(ctx, func(ctx context.Context, tx store.UnPostStore) error {
		err = tx.CreatePost(ctx, post)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreatePostResponse{
		Post: &v1.Post{
			Id:        post.ID,
			Title:     doc.GetDocument().GetTitle(),
			CreatedAt: timestamppb.New(post.CreatedAt),
			UpdatedAt: timestamppb.New(post.UpdatedAt),
		},
	}, nil
}

func (p *PostService) GetPost(ctx context.Context, request *v1.GetPostRequest) (*v1.GetPostResponse, error) {
	post, err := p.store.GetPost(ctx, uuid.MustParse(request.GetId()))
	if err != nil {
		return nil, err
	}

	res, err := p.docClient.GetDocument(p.cfg.IntoContext(), &docv1.GetDocumentRequest{
		Id: post.DocumentID,
	})
	if err != nil {
		return nil, err
	}

	doc := res.GetDocument()
	postProto := &v1.Post{
		Id:        post.ID,
		Title:     doc.GetTitle(),
		Content:   doc.GetContent(),
		Summary:   doc.GetSummary(),
		Excerpt:   doc.GetExcerpt(),
		Thumbnail: doc.GetThumbnail(),
		Tags:      make([]*v1.Tag, 0),
		Version:   doc.GetVersion(),
		Status:    postStatusToProto(post.Status),
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
	poolID, err := authx.GetAuthbasePoolID(ctx)
	if err != nil {
		return nil, err
	}

	userID, err := authx.GetAuthbaseAccountID(ctx)
	if err != nil {
		return nil, err
	}

	var status *model.PostStatus = nil
	if request.Status != nil {
		statusModel := postStatusFromProto(request.GetStatus())
		status = &statusModel
	}

	posts, err := p.store.ListPostByOwnerID(ctx, userID, status)
	if err != nil {
		return nil, err
	}

	docIDs := make([]string, 0)
	for _, post := range posts {
		docIDs = append(docIDs, post.DocumentID)
	}

	// get the documents from the document service
	res, err := p.docClient.ListDocuments(p.cfg.IntoContext(), &docv1.ListDocumentsRequest{
		ProjectId:   poolID.String(),
		DocumentIds: docIDs,
	})
	if err != nil {
		return nil, err
	}

	// document map
	documents := make(map[string]*docv1.Document)
	for _, doc := range res.GetDocuments() {
		documents[doc.GetId()] = doc
	}

	userIDs := mapset.NewSet[string]()
	for _, post := range posts {
		userIDs.Add(post.CreatedByID)
	}

	users := make(map[string]*authv1.Account)
	//if res, err := p.authClient.ListUsers(ctx, &authv1.ListUsersRequest{}); err != nil {
	//	return nil, err
	//} else {
	//	for _, user := range res.GetUsers() {
	//		users[user.Id] = user
	//	}
	//}

	var responsePosts []*v1.Post
	for _, post := range posts {
		doc, ok := documents[post.DocumentID]
		if !ok {
			continue
		}
		postProto := &v1.Post{
			Id:        post.ID,
			Status:    postStatusToProto(post.Status),
			Title:     doc.GetTitle(),
			Summary:   doc.Summary,
			Excerpt:   doc.Excerpt,
			Thumbnail: doc.Thumbnail,
			Version:   doc.Version,
			CreatedAt: timestamppb.New(post.CreatedAt),
			UpdatedAt: timestamppb.New(post.UpdatedAt),
		}
		responsePosts = append(responsePosts, postProto)
		if user, ok := users[post.CreatedByID]; ok {
			postProto.OriginalAuthor = &v1.User{
				Id:    user.Id,
				Name:  user.Username,
				Email: user.Email,
			}
		}
	}

	return &v1.ListPostResponse{
		Posts: responsePosts,
	}, nil
}

func (p *PostService) UpdatePost(ctx context.Context, request *v1.UpdatePostRequest) (*v1.UpdatePostResponse, error) {
	err := p.store.Transaction(ctx, func(ctx context.Context, tx store.UnPostStore) error {
		post, err := tx.GetPost(ctx, uuid.MustParse(request.GetId()))
		if err != nil {
			return err
		}

		logrus.Infof("updating post %d", request.GetVersion())

		_, err = p.docClient.UpdateDocument(p.cfg.IntoContext(), &docv1.UpdateDocumentRequest{
			Id:        post.DocumentID,
			Title:     request.Title,
			Content:   request.Content,
			Summary:   request.Summary,
			Excerpt:   request.Excerpt,
			Thumbnail: request.Thumbnail,
			Version:   request.GetVersion(),
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdatePostResponse{
		Post: &v1.Post{
			Id:      request.GetId(),
			Version: request.GetVersion(),
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

	err := p.store.Transaction(ctx, func(ctx context.Context, tx store.UnPostStore) error {
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
	err := p.store.Transaction(ctx, func(ctx context.Context, tx store.UnPostStore) error {
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

	err = p.store.Transaction(ctx, func(ctx context.Context, tx store.UnPostStore) error {
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
	err := p.store.Transaction(ctx, func(ctx context.Context, tx store.UnPostStore) error {
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
