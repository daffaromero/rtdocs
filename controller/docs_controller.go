package controller

import (
	"encoding/json"
	"net/http"
	"rtdocs/model"
	"rtdocs/service"
)

type DocumentController interface {
	GetDocument(w http.ResponseWriter, r *http.Request)
	GetAllDocuments(w http.ResponseWriter, r *http.Request)
	UpdateDocumentContent(w http.ResponseWriter, r *http.Request)
}

type documentController struct {
	docService service.DocumentService
}

func NewDocumentController(docService service.DocumentService) DocumentController {
	return &documentController{docService: docService}
}

// GetDocument retrieves a document by its ID
func (c *documentController) GetDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Document ID is required", http.StatusBadRequest)
		return
	}

	document, err := c.docService.GetDocument(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(document)
}

// GetAllDocuments retrieves all documents
func (c *documentController) GetAllDocuments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	documents, err := c.docService.GetAllDocuments(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(documents)
}

// UpdateDocumentContent updates the content of a document
func (c *documentController) UpdateDocumentContent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var document model.Document
	if err := json.NewDecoder(r.Body).Decode(&document); err != nil {
		http.Error(w, "Invalid document data", http.StatusBadRequest)
		return
	}

	createdDoc, err := c.docService.UpdateDocumentContent(ctx, &document)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdDoc)
}
