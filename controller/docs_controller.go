package controller

import (
	"encoding/json"
	"net/http"
	"rtdocs/model/domain"
	"rtdocs/service"
)

type DocumentController interface {
	GetDocument(w http.ResponseWriter, r *http.Request)
	GetAllDocuments(w http.ResponseWriter, r *http.Request)
	CreateDocument(w http.ResponseWriter, r *http.Request)
	UpdateDocument(w http.ResponseWriter, r *http.Request)
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

	id := r.PathValue("id")
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

func (c *documentController) CreateDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var document domain.Document
	if err := json.NewDecoder(r.Body).Decode(&document); err != nil {
		http.Error(w, "Invalid document data", http.StatusBadRequest)
		return
	}

	createdDoc, err := c.docService.CreateDocument(ctx, &document)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdDoc)
}

// UpdateDocument updates the content of a document
func (c *documentController) UpdateDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var document domain.Document
	if err := json.NewDecoder(r.Body).Decode(&document); err != nil {
		http.Error(w, "Invalid document data", http.StatusBadRequest)
		return
	}

	updatedDoc, err := c.docService.UpdateDocument(ctx, &document)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedDoc)
}
