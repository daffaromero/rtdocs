package query

import (
	"context"
	"errors"
	"fmt"
	"log"
	"rtdocs/model"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DocumentQuery interface {
	GetDocument(ctx context.Context, id string) (*model.Document, error)
	GetAllDocuments(ctx context.Context) ([]*model.Document, error)
	SaveDocument(ctx context.Context, document *model.Document) error
}

type documentQuery struct {
	db *pgxpool.Pool
}

func NewDocumentQuery(db *pgxpool.Pool) DocumentQuery {
	return &documentQuery{db: db}
}

func (q *documentQuery) GetDocument(ctx context.Context, id string) (*model.Document, error) {
	if id == "" {
		log.Println("Document ID is required")
		return nil, nil
	}
	query := "SELECT id, title, content FROM docs WHERE id = $1"

	var document model.Document
	row := q.db.QueryRow(ctx, query, id)

	if err := row.Scan(&document.ID, &document.Title, &document.Content); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("document not found: %w", err)
		}
		return nil, err
	}

	return &document, nil
}

func (q *documentQuery) GetAllDocuments(ctx context.Context) ([]*model.Document, error) {
	query := "SELECT id, title, content FROM docs"
	rows, err := q.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []*model.Document
	for rows.Next() {
		var document model.Document
		if err := rows.Scan(&document.ID, &document.Title, &document.Content); err != nil {
			return nil, err
		}
		documents = append(documents, &document)
	}
	return documents, nil
}

func (q *documentQuery) SaveDocument(ctx context.Context, document *model.Document) error {
	if document.ID == "" {
		log.Println("Document ID is required")
		return nil
	}

	query := "INSERT INTO docs (id, title, content) VALUES ($1, $2, $3) ON CONFLICT(id) DO UPDATE SET content = $2 RETURNING id, title, content"
	if err := q.db.QueryRow(ctx, query, document.ID, document.Title, document.Content).Scan(&document.ID, &document.Title, &document.Content); err != nil {
		return err
	}
	return nil
}
