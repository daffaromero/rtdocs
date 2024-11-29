package service

import (
	"context"
	"rtdocs/model/domain"
	"rtdocs/model/web"
	"rtdocs/repository"
	"time"

	"github.com/google/uuid"
)

type DocumentService interface {
	GetDocument(ctx context.Context, id string) (*domain.Document, error)
	GetAllDocuments(ctx context.Context) ([]*domain.Document, error)
	CreateDocument(ctx context.Context, newDoc *web.CreateDocument) (*domain.Document, error)
	UpdateDocument(ctx context.Context, updatedDoc *domain.Document) (*domain.Document, error)
}

type documentService struct {
	repo repository.DocumentRepository
}

func NewDocumentService(repo repository.DocumentRepository) DocumentService {
	return &documentService{repo: repo}
}

func (s *documentService) GetDocument(ctx context.Context, id string) (*domain.Document, error) {
	return s.repo.GetDocument(ctx, id)
}

func (s *documentService) GetAllDocuments(ctx context.Context) ([]*domain.Document, error) {
	return s.repo.GetAllDocuments(ctx)
}

func (s *documentService) CreateDocument(ctx context.Context, request *web.CreateDocument) (*domain.Document, error) {
	var newDoc domain.Document
	newDoc.ID = uuid.New().String()
	if request.Title == "" {
		newDoc.Title = "Untitled Document"
	} else {
		newDoc.Title = request.Title
	}
	newDoc.Content = request.Content
	newDoc.OwnerID = request.OwnerID
	newDoc.IsPublic = false
	newDoc.CanEdit = true
	newDoc.CreatedAt = time.Now()
	newDoc.UpdatedAt = time.Now()

	return s.repo.CreateDocument(ctx, &newDoc)
}

func (s *documentService) UpdateDocument(ctx context.Context, updatedDoc *domain.Document) (*domain.Document, error) {
	if updatedDoc.Title == "" {
		updatedDoc.Title = "Untitled Document"
	}

	return s.repo.UpdateDocument(ctx, updatedDoc)
}
