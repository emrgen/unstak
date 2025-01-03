package service

import (
	"context"
	v1 "github.com/emrgen/unpost/apis/v1"
)

func NewBookService() *BookService {
	return &BookService{}
}

var _ v1.BookServiceServer = new(BookService)

type BookService struct {
	v1.UnimplementedBookServiceServer
}

func (b *BookService) CreateBook(ctx context.Context, request *v1.CreateBookRequest) (*v1.CreateBookResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BookService) GetBook(ctx context.Context, request *v1.GetBookRequest) (*v1.GetBookResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BookService) ListBook(ctx context.Context, request *v1.ListBookRequest) (*v1.ListBookResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BookService) UpdateBook(ctx context.Context, request *v1.UpdateBookRequest) (*v1.UpdateBookResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BookService) DeleteBook(ctx context.Context, request *v1.DeleteBookRequest) (*v1.DeleteBookResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BookService) AddBookTag(ctx context.Context, request *v1.AddBookTagRequest) (*v1.AddBookTagResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BookService) RemoveBookTag(ctx context.Context, request *v1.RemoveBookTagRequest) (*v1.RemoveBookTagResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BookService) mustEmbedUnimplementedBookServiceServer() {
	//TODO implement me
	panic("implement me")
}
