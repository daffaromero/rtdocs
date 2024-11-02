package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"rtdocs/model"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DocumentRepository interface {
	GetDocument(ctx context.Context, id string) (*model.Document, error)
	GetAllDocuments(ctx context.Context) ([]*model.Document, error)
	CreateDocument(ctx context.Context, document *model.Document) (*model.Document, error)
	UpdateDocument(ctx context.Context, document *model.Document) (*model.Document, error)
	ShareDocument(ctx context.Context, document *model.Document) (*model.Document, error)
}

type documentRepository struct {
	db *pgxpool.Pool
}

func NewDocumentRepository(db *pgxpool.Pool) DocumentRepository {
	return &documentRepository{db: db}
}

func (q *documentRepository) GetDocument(ctx context.Context, id string) (*model.Document, error) {
	if id == "" {
		log.Println("Document ID is required")
		return nil, nil
	}
	query := "SELECT * FROM docs WHERE id = $1"

	var document model.Document
	row := q.db.QueryRow(ctx, query, id)

	if err := row.Scan(&document.ID, &document.Title, &document.Content, &document.OwnerID, &document.IsPublic, &document.CanEdit, &document.CreatedAt, &document.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("document not found: %w", err)
		}
		return nil, err
	}

	return &document, nil
}

func (q *documentRepository) GetAllDocuments(ctx context.Context) ([]*model.Document, error) {
	query := "SELECT * FROM docs"
	rows, err := q.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []*model.Document
	for rows.Next() {
		var document model.Document
		if err := rows.Scan(&document.ID, &document.Title, &document.Content, &document.OwnerID, &document.IsPublic, &document.CanEdit, &document.CreatedAt, &document.UpdatedAt); err != nil {
			return nil, err
		}
		documents = append(documents, &document)
	}

	return documents, nil
}

func (q *documentRepository) CreateDocument(ctx context.Context, document *model.Document) (*model.Document, error) {
	query := "INSERT INTO docs (id, title, content, owner_id, is_public, can_edit, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, title, content, owner_id, is_public, can_edit, created_at, updated_at"
	if err := q.db.QueryRow(ctx, query, document.ID, document.Title, document.Content, document.OwnerID, document.IsPublic, document.CanEdit, document.CreatedAt, document.UpdatedAt).Scan(&document.ID, &document.Title, &document.Content, &document.OwnerID, &document.IsPublic, &document.CanEdit, &document.CreatedAt, &document.UpdatedAt); err != nil {
		return nil, err
	}

	return document, nil
}

func (q *documentRepository) UpdateDocument(ctx context.Context, document *model.Document) (*model.Document, error) {
	if document.ID == "" {
		log.Println("Document ID is required")
		return nil, errors.New("document ID is required")
	}

	query := "UPDATE docs SET title = $1, content = $2, owner_id = $3, is_public = $4, can_edit = $5, updated_at = $6 WHERE id = $7 RETURNING id, title, content, owner_id, is_public, can_edit, created_at, updated_at"

	if err := q.db.QueryRow(ctx, query, document.Title, document.Content, document.OwnerID, document.IsPublic, document.CanEdit, document.UpdatedAt, document.ID).Scan(&document.ID, &document.Title, &document.Content, &document.OwnerID, &document.IsPublic, &document.CanEdit, &document.CreatedAt, &document.UpdatedAt); err != nil {
		return nil, err
	}

	return document, nil
}

func (q *documentRepository) ShareDocument(ctx context.Context, document *model.Document) (*model.Document, error) {
	if document.ID == "" {
		log.Println("Document ID is required")
		return nil, errors.New("document ID is required")
	}

	query := "UPDATE docs SET is_public = $1, can_edit = $2, updated_at = $3 WHERE id = $4 RETURNING id, title, content, owner_id, is_public, can_edit, created_at, updated_at"

	if err := q.db.QueryRow(ctx, query, document.IsPublic, document.CanEdit, document.UpdatedAt, document.ID).Scan(&document.ID, &document.Title, &document.Content, &document.OwnerID, &document.IsPublic, &document.CanEdit, &document.CreatedAt, &document.UpdatedAt); err != nil {
		return nil, err
	}

	return document, nil
}
