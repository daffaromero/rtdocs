package repository

import (
	"context"
	"rtdocs/model"
	"rtdocs/repository/query"
)

type DocumentRepository interface {
	GetDocument(ctx context.Context, id string) (*model.Document, error)
	GetAllDocuments(ctx context.Context) ([]*model.Document, error)
	SaveDocument(ctx context.Context, document *model.Document) error
}

type documentRepository struct {
	docQuery query.DocumentQuery
}

func NewDocumentRepository(docQuery query.DocumentQuery) DocumentRepository {
	return &documentRepository{
		docQuery: docQuery,
	}
}

func (r *documentRepository) GetDocument(ctx context.Context, id string) (*model.Document, error) {
	return r.docQuery.GetDocument(ctx, id)
}

func (r *documentRepository) GetAllDocuments(ctx context.Context) ([]*model.Document, error) {
	return r.docQuery.GetAllDocuments(ctx)
}

func (r *documentRepository) SaveDocument(ctx context.Context, document *model.Document) error {
	return r.docQuery.SaveDocument(ctx, document)
}
