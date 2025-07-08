package service

import (
	"context"
	authx "github.com/emrgen/authbase/x"
	docv1 "github.com/emrgen/document/apis/v1"
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/emrgen/unpost/internal/model"
	"github.com/emrgen/unpost/internal/store"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NewPageService creates a new book service
func NewPageService(cfg *authx.AuthbaseConfig, store store.UnstakStore, docClient docv1.DocumentServiceClient) *PageService {
	return &PageService{
		cfg:       cfg,
		docClient: docClient,
		store:     store,
	}
}

var _ v1.PageServiceServer = new(PageService)

type PageService struct {
	cfg       *authx.AuthbaseConfig
	store     store.UnstakStore
	docClient docv1.DocumentServiceClient
	v1.UnimplementedPageServiceServer
}

func (p *PageService) CreatePage(ctx context.Context, request *v1.CreatePageRequest) (*v1.CreatePageResponse, error) {
	userID, err := authx.GetAuthbaseAccountID(ctx)
	if err != nil {
		return nil, err
	}

	page := &model.Page{
		ID:          uuid.New().String(),
		CreatedByID: userID.String(),
	}

	if err := p.store.CreatePage(ctx, page); err != nil {
		return nil, err
	}

	return &v1.CreatePageResponse{
		Page: &v1.Page{
			Id: page.ID,
		},
	}, nil
}

func (p *PageService) GetPage(ctx context.Context, request *v1.GetPageRequest) (*v1.GetPageResponse, error) {
	pageID := uuid.MustParse(request.Id)
	page, err := p.store.GetPage(ctx, pageID)
	if err != nil {
		return nil, err
	}

	// TODO: parse the doc.Meta field to get the title, summary, excerpt, and thumbnail
	pageProto := &v1.Page{
		Id:          page.ID,
		CreatedById: page.CreatedByID,
		Status:      postStatusToProto(page.Status),
		//Title:       doc.GetTitle(),
		//Summary:     doc.Summary,
		//Excerpt:     doc.Excerpt,
		//Thumbnail:   doc.Thumbnail,
		CreatedAt: timestamppb.New(page.CreatedAt),
		UpdatedAt: timestamppb.New(page.UpdatedAt),
	}

	return &v1.GetPageResponse{
		Page: pageProto,
	}, nil
}

func (p *PageService) UpdatePage(ctx context.Context, request *v1.UpdatePageRequest) (*v1.UpdatePageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PageService) DeletePage(ctx context.Context, request *v1.DeletePageRequest) (*v1.DeletePageResponse, error) {
	pageID := uuid.MustParse(request.Id)
	if err := p.store.DeletePage(ctx, pageID); err != nil {
		return nil, err
	}

	return &v1.DeletePageResponse{}, nil
}

func (p *PageService) AddPageTag(ctx context.Context, request *v1.AddPageTagRequest) (*v1.AddPageTagResponse, error) {
	pageID := uuid.MustParse(request.GetPageId())
	tagID := uuid.MustParse(request.GetTagId())

	err := p.store.Transaction(ctx, func(ctx context.Context, tx store.UnstakStore) error {
		post, err := tx.GetPost(ctx, pageID)
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

	return &v1.AddPageTagResponse{}, nil
}

func (p *PageService) RemovePageTag(ctx context.Context, request *v1.RemovePageTagRequest) (*v1.RemovePageTagResponse, error) {
	pageID := uuid.MustParse(request.GetPageId())
	tagID := uuid.MustParse(request.GetTagId())

	err := p.store.Transaction(ctx, func(ctx context.Context, tx store.UnstakStore) error {
		post, err := tx.GetPost(ctx, pageID)
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
				break
			}
		}

		err = tx.UpdatePost(ctx, post)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.RemovePageTagResponse{}, nil
}
