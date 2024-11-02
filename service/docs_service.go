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
	CreateDocument(ctx context.Context, newDoc *model.Document) (*model.Document, error)
	UpdateDocument(ctx context.Context, updatedDoc *model.Document) (*model.Document, error)
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

func (s *documentService) CreateDocument(ctx context.Context, newDoc *model.Document) (*model.Document, error) {
	newDoc.ID = uuid.New().String()
	if newDoc.Title == "" {
		newDoc.Title = "Untitled Document"
	}

	return s.repo.CreateDocument(ctx, newDoc)
}

func (s *documentService) UpdateDocument(ctx context.Context, updatedDoc *model.Document) (*model.Document, error) {
	if updatedDoc.Title == "" {
		updatedDoc.Title = "Untitled Document"
	}

	return s.repo.UpdateDocument(ctx, updatedDoc)
}
