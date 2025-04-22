package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/emrgen/authbase"
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
func NewPostService(cfg *authx.AuthbaseConfig, store store.UnstakStore, docClient docv1.DocumentServiceClient, authClient authbase.Client) *PostService {
	return &PostService{
		cfg:        cfg,
		store:      store,
		docClient:  docClient,
		authClient: authClient,
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

func (p *PostService) CreatePost(ctx context.Context, request *v1.CreatePostRequest) (*v1.CreatePostResponse, error) {
	var err error
	poolID, err := authx.GetAuthbasePoolID(ctx)
	if err != nil {
		return nil, err
	}
	accountID, err := authx.GetAuthbaseAccountID(ctx)
	if err != nil {
		return nil, err
	}

	docID := uuid.New().String()
	if request.GetPostId() != "" {
		docID = request.GetPostId()
	}

	req := &docv1.CreateDocumentRequest{
		DocumentId: &docID,
		ProjectId:  poolID.String(),
		Content:    request.GetContent(),
	}

	meta := map[string]string{
		"title":   request.GetTitle(),
		"authors": fmt.Sprintf("%d", accountID),
	}
	metaData, err := json.Marshal(meta)
	if err != nil {
		return nil, err
	}
	req.Meta = string(metaData)

	// create a document in the document service
	doc, err := p.docClient.CreateDocument(ctx, req)
	if err != nil {
		return nil, err
	}

	logrus.Infof("created document %s", doc.GetDocument().GetId())

	postID := uuid.New().String()
	if request.PostId != nil {
		postID = request.GetPostId()
	}

	post := &model.Post{
		ID:          postID,
		SpaceID:     request.GetSpaceId(),
		CreatedByID: accountID.String(),
		DocumentID:  doc.GetDocument().GetId(),
	}

	// get user default space-id if no space-id is provided
	err = p.store.Transaction(ctx, func(ctx context.Context, tx store.UnstakStore) error {
		err = tx.CreatePost(ctx, post)
		if err != nil {
			return err
		}

		return nil
	})

	// if the post creation fails, we need to erase the document created in the document service
	// TODO: use a transaction to rollback the document creation
	if err != nil {
		err2 := eraseDocument(ctx, p.docClient, doc.GetDocument().GetId())
		if err2 != nil {
			return nil, errors.Join(err, err2)
		}
		return nil, err
	}

	return &v1.CreatePostResponse{
		Post: &v1.Post{
			Id:        post.ID,
			CreatedAt: timestamppb.New(post.CreatedAt),
			UpdatedAt: timestamppb.New(post.UpdatedAt),
			Version:   doc.GetDocument().GetVersion(),
		},
	}, nil
}

func (p *PostService) GetPost(ctx context.Context, request *v1.GetPostRequest) (*v1.GetPostResponse, error) {
	post, err := p.store.GetPost(ctx, uuid.MustParse(request.GetId()))
	if err != nil {
		return nil, err
	}

	res, err := p.docClient.GetDocument(p.cfg.IntoContext(), &docv1.GetDocumentRequest{
		DocumentId: post.DocumentID,
	})
	if err != nil {
		return nil, err
	}

	poolID, err := authx.GetAuthbasePoolID(ctx)
	if err != nil {
		return nil, err
	}

	poolId := poolID.String()
	authorsRes, err := p.authClient.ListAccounts(p.cfg.IntoContext(), &authv1.ListAccountsRequest{
		PoolId:     &poolId,
		AccountIds: []string{post.CreatedByID},
	})
	if err != nil {
		return nil, err
	}

	var authors []*v1.Account
	for _, user := range authorsRes.GetAccounts() {
		authors = append(authors, &v1.Account{
			Id:    user.Id,
			Name:  user.Username,
			Email: user.Email,
		})
	}

	doc := res.GetDocument()
	if doc == nil {
		return nil, errors.New("document not found")
	}

	meta := map[string]string{}
	err = json.Unmarshal([]byte(doc.GetMeta()), &meta)
	if err != nil {
		return nil, err
	}
	title, ok := meta["title"]
	if !ok {
		title = ""
	}
	postProto := &v1.Post{
		Id:      post.ID,
		Content: doc.GetContent(),
		Tags:    make([]*v1.Tag, 0),
		Title:   title,
		//TODO: return main author
		//MainAuthor:
		Version: doc.GetVersion(),
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

	var posts []*model.Post
	if request.GetSpaceId() != "" {
		posts, err = p.store.ListPostBySpace(ctx, uuid.MustParse(request.GetSpaceId()), status)
	} else {
		posts, err = p.store.ListPostByOwnerID(ctx, userID, status)
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
	poolId := poolID.String()
	if res, err := p.authClient.ListAccounts(p.cfg.IntoContext(), &authv1.ListAccountsRequest{
		PoolId:     &poolId,
		AccountIds: userIDs.ToSlice(),
	}); err != nil {
		return nil, err
	} else {
		for _, user := range res.GetAccounts() {
			users[user.Id] = user
		}
	}

	var responsePosts []*v1.Post
	for _, post := range posts {
		doc, ok := documents[post.DocumentID]
		if !ok {
			continue
		}
		meta := map[string]string{}
		err = json.Unmarshal([]byte(doc.GetMeta()), &meta)
		if err != nil {
			return nil, err
		}
		title := ""
		if t, ok := meta["title"]; ok {
			title = t
		}

		postProto := &v1.Post{
			Id:     post.ID,
			Status: postStatusToProto(post.Status),
			Title:  title,
			//Summary:   doc.Summary,
			//Excerpt:   doc.Excerpt,
			//Thumbnail: doc.Thumbnail,
			Version:   doc.Version,
			CreatedAt: timestamppb.New(post.CreatedAt),
			UpdatedAt: timestamppb.New(post.UpdatedAt),
		}
		responsePosts = append(responsePosts, postProto)
		if user, ok := users[post.CreatedByID]; ok {
			postProto.MainAuthor = &v1.Account{
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
	meta := make(map[string]string)
	meta["title"] = request.GetTitle()
	meta["summary"] = request.GetSummary()
	meta["excerpt"] = request.GetExcerpt()
	meta["thumbnail"] = request.GetThumbnail()
	meta["authors"] = request.GetAuthors()

	// meta string
	marshal, err := json.Marshal(meta)
	if err != nil {
		return nil, err
	}
	metaStr := string(marshal)

	_, err = p.docClient.UpdateDocument(p.cfg.IntoContext(), &docv1.UpdateDocumentRequest{
		DocumentId: request.GetPostId(),
		Content:    request.Content,
		Meta:       &metaStr,
		Version:    request.GetVersion(),
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdatePostResponse{
		Post: &v1.Post{
			Id:      request.GetPostId(),
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

func eraseDocument(ctx context.Context, docClient docv1.DocumentServiceClient, id string) error {
	_, err := docClient.EraseDocument(ctx, &docv1.EraseDocumentRequest{
		Id: id,
	})
	return err
}
