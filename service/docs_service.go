package service

import (
	"context"
	"rtdocs/model"
	"rtdocs/repository"

	"github.com/google/uuid"
)

type DocumentService interface {
	GetDocument(ctx context.Context, id string) (*model.Document, error)
	GetAllDocuments(ctx context.Context) ([]*model.Document, error)
	UpdateDocumentContent(ctx context.Context, updatedDoc *model.Document) error
}

type documentService struct {
	repo repository.DocumentRepository
}

func NewDocumentService(repo repository.DocumentRepository) DocumentService {
	return &documentService{repo: repo}
}

func (s *documentService) GetDocument(ctx context.Context, id string) (*model.Document, error) {
	return s.repo.GetDocument(ctx, id)
}

func (s *documentService) GetAllDocuments(ctx context.Context) ([]*model.Document, error) {
	return s.repo.GetAllDocuments(ctx)
}

func (s *documentService) UpdateDocumentContent(ctx context.Context, updatedDoc *model.Document) error {
	updatedDoc.ID = uuid.New().String()
	return s.repo.SaveDocument(ctx, updatedDoc)
}
