package service

import (
	"context"
	authx "github.com/emrgen/authbase/x"
	docv1 "github.com/emrgen/document/apis/v1"
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/emrgen/unpost/internal/model"
	"github.com/emrgen/unpost/internal/store"
	"github.com/google/uuid"
)

func NewCourseService(cfg *authx.AuthbaseConfig, store store.UnstakStore, docClient docv1.DocumentServiceClient) *CourseService {
	return &CourseService{
		cfg:       cfg,
		store:     store,
		docClient: docClient,
	}
}

var _ v1.CourseServiceServer = new(CourseService)

type CourseService struct {
	cfg       *authx.AuthbaseConfig
	store     store.UnstakStore
	docClient docv1.DocumentServiceClient
	v1.UnimplementedCourseServiceServer
}

func (c *CourseService) CreateCourse(ctx context.Context, request *v1.CreateCourseRequest) (*v1.CreateCourseResponse, error) {
	poolID, err := authx.GetAuthbasePoolID(ctx)
	if err != nil {
		return nil, err
	}

	userID, err := authx.GetAuthbaseAccountID(ctx)
	if err != nil {
		return nil, err
	}

	res, err := c.docClient.CreateDocument(c.cfg.IntoContext(), &docv1.CreateDocumentRequest{
		ProjectId: poolID.String(),
	})
	if err != nil {
		return nil, err
	}

	course := &model.Course{
		ID:          uuid.New().String(),
		DocumentID:  res.Document.Id,
		CreatedByID: userID.String(),
		Status:      model.PostStatusDraft,
	}

	if err := c.store.CreateCourse(ctx, course); err != nil {
		return nil, err
	}

	return &v1.CreateCourseResponse{
		Course: &v1.Course{
			Id: course.ID,
		},
	}, nil
}

func (c *CourseService) GetCourse(ctx context.Context, request *v1.GetCourseRequest) (*v1.GetCourseResponse, error) {
	course, err := c.store.GetCourse(ctx, uuid.MustParse(request.GetId()))
	if err != nil {
		return nil, err
	}

	res, err := c.docClient.GetDocument(c.cfg.IntoContext(), &docv1.GetDocumentRequest{
		DocumentId: course.DocumentID,
	})
	if err != nil {
		return nil, err
	}

	doc := res.GetDocument()
	page := &v1.Page{
		Id: course.ID,
		//Title:     doc.GetTitle(),
		Content: doc.GetContent(),
		//Summary:   doc.GetSummary(),
		//Excerpt:   doc.GetExcerpt(),
		//Thumbnail: doc.GetThumbnail(),
		Tags:    make([]*v1.Tag, 0),
		Version: doc.GetVersion(),
		Status:  postStatusToProto(course.Status),
	}

	courseProto := &v1.Course{
		Id:          course.ID,
		CoverPage:   page,
		CreatedById: course.CreatedByID,
		Version:     page.Version,
	}

	for _, tag := range course.Tags {
		courseProto.Tags = append(courseProto.Tags, &v1.Tag{
			Id:   tag.ID,
			Name: tag.Name,
		})
	}

	return &v1.GetCourseResponse{
		Course: courseProto,
	}, nil
}

// ListCourse returns a list of courses.
// 1. List by owner
// 2. List by tag
// 3. List by status
// 4. List by reaction
func (c *CourseService) ListCourse(ctx context.Context, request *v1.ListCourseRequest) (*v1.ListCourseResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CourseService) UpdateCourse(ctx context.Context, request *v1.UpdateCourseRequest) (*v1.UpdateCourseResponse, error) {
	return &v1.UpdateCourseResponse{}, nil
}

func (c *CourseService) DeleteCourse(ctx context.Context, request *v1.DeleteCourseRequest) (*v1.DeleteCourseResponse, error) {
	courseID := uuid.MustParse(request.GetId())
	if err := c.store.DeleteCourse(ctx, courseID); err != nil {
		return nil, err
	}

	return &v1.DeleteCourseResponse{}, nil
}

func (c *CourseService) AddCourseTag(ctx context.Context, request *v1.AddCourseTagRequest) (*v1.AddCourseTagResponse, error) {
	courseID := uuid.MustParse(request.GetCourseId())
	tagID := uuid.MustParse(request.GetTagId())

	err := c.store.Transaction(ctx, func(ctx context.Context, tx store.UnstakStore) error {
		course, err := tx.GetCourse(ctx, courseID)
		if err != nil {
			return err
		}

		tag, err := tx.GetTag(ctx, tagID)
		if err != nil {
			return err
		}

		course.Tags = append(course.Tags, tag)

		err = tx.UpdateCourse(ctx, course)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.AddCourseTagResponse{}, nil
}

func (c *CourseService) RemoveCourseTag(ctx context.Context, request *v1.RemoveCourseTagRequest) (*v1.RemoveCourseTagResponse, error) {
	courseID := uuid.MustParse(request.GetCourseId())
	tagID := uuid.MustParse(request.GetTagId())

	err := c.store.Transaction(ctx, func(ctx context.Context, tx store.UnstakStore) error {
		course, err := tx.GetCourse(ctx, courseID)
		if err != nil {
			return err
		}

		tag, err := tx.GetTag(ctx, tagID)
		if err != nil {
			return err
		}

		for i, t := range course.Tags {
			if t.ID == tag.ID {
				course.Tags = append(course.Tags[:i], course.Tags[i+1:]...)
				break
			}
		}

		err = tx.UpdateCourse(ctx, course)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.RemoveCourseTagResponse{}, nil
}
