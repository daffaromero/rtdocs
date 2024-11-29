package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"rtdocs/model/domain"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DocumentRepository interface {
	GetDocument(ctx context.Context, id string) (*domain.Document, error)
	GetAllDocuments(ctx context.Context) ([]*domain.Document, error)
	CreateDocument(ctx context.Context, document *domain.Document) (*domain.Document, error)
	UpdateDocument(ctx context.Context, document *domain.Document) (*domain.Document, error)
	ShareDocument(ctx context.Context, document *domain.Document) (*domain.Document, error)
}

type documentRepository struct {
	db *pgxpool.Pool
}

func NewDocumentRepository(db *pgxpool.Pool) DocumentRepository {
	return &documentRepository{db: db}
}

func (q *documentRepository) GetDocument(ctx context.Context, id string) (*domain.Document, error) {
	if id == "" {
		return nil, nil
	}
	query := "SELECT * FROM docs WHERE id = $1"

	var document domain.Document
	row := q.db.QueryRow(ctx, query, id)

	if err := row.Scan(&document.ID, &document.Title, &document.Content, &document.OwnerID, &document.IsPublic, &document.CanEdit, &document.CreatedAt, &document.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("document not found: %w", err)
		}
		return nil, err
	}

	return &document, nil
}

func (q *documentRepository) GetAllDocuments(ctx context.Context) ([]*domain.Document, error) {
	query := "SELECT * FROM docs"
	rows, err := q.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []*domain.Document
	for rows.Next() {
		var document domain.Document
		if err := rows.Scan(&document.ID, &document.Title, &document.Content, &document.OwnerID, &document.IsPublic, &document.CanEdit, &document.CreatedAt, &document.UpdatedAt); err != nil {
			return nil, err
		}
		documents = append(documents, &document)
	}

	return documents, nil
}

func (q *documentRepository) CreateDocument(ctx context.Context, document *domain.Document) (*domain.Document, error) {
	var newDoc domain.Document
	query := "INSERT INTO docs (id, title, content, owner_id, is_public, can_edit, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, title, content, owner_id, is_public, can_edit, created_at, updated_at"
	row := q.db.QueryRow(ctx, query, document.ID, document.Title, document.Content, document.OwnerID, document.IsPublic, document.CanEdit, document.CreatedAt, document.UpdatedAt)
	if err := row.Scan(&newDoc.ID, &newDoc.Title, &newDoc.Content, &newDoc.OwnerID, &newDoc.IsPublic, &newDoc.CanEdit, &newDoc.CreatedAt, &newDoc.UpdatedAt); err != nil {
		log.Println(err)
		// return nil, err
	}

	return &newDoc, nil
}

func (q *documentRepository) UpdateDocument(ctx context.Context, document *domain.Document) (*domain.Document, error) {
	var updatedDoc domain.Document
	log.Println("ID:", document.ID)
	log.Println("Title:", document.Title)
	log.Println("Content:", document.Content)
	log.Println("OwnerID:", document.OwnerID)
	log.Println("IsPublic:", document.IsPublic)
	log.Println("CanEdit:", document.CanEdit)
	log.Println("UpdatedAt:", document.UpdatedAt)

	if document.ID == "" {
		return nil, errors.New("document ID is required")
	}

	query := "UPDATE docs SET title = $1, content = $2, owner_id = $3, is_public = $4, can_edit = $5, updated_at = $6 WHERE id = $7 RETURNING id, title, content, owner_id, is_public, can_edit, created_at, updated_at"

	row := q.db.QueryRow(ctx, query, document.Title, document.Content, document.OwnerID, document.IsPublic, document.CanEdit, document.UpdatedAt, document.ID)
	if err := row.Scan(&updatedDoc.ID, &updatedDoc.Title, &updatedDoc.Content, &updatedDoc.OwnerID, &updatedDoc.IsPublic, &updatedDoc.CanEdit, &updatedDoc.CreatedAt, &updatedDoc.UpdatedAt); err != nil {
		return nil, err
	}

	return &updatedDoc, nil
}

func (q *documentRepository) ShareDocument(ctx context.Context, document *domain.Document) (*domain.Document, error) {
	var sharedDoc domain.Document
	if document.ID == "" {
		return nil, errors.New("document ID is required")
	}

	query := "UPDATE docs SET is_public = $1, can_edit = $2, updated_at = $3 WHERE id = $4 RETURNING id, title, content, owner_id, is_public, can_edit, created_at, updated_at"

	row := q.db.QueryRow(ctx, query, document.IsPublic, document.CanEdit, document.UpdatedAt, document.ID)
	if err := row.Scan(&sharedDoc.ID, &sharedDoc.Title, &sharedDoc.Content, &sharedDoc.OwnerID, &sharedDoc.IsPublic, &sharedDoc.CanEdit, &sharedDoc.CreatedAt, &sharedDoc.UpdatedAt); err != nil {
		return nil, err
	}

	return &sharedDoc, nil
}
