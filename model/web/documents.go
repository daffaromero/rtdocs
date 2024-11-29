package web

type CreateDocument struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	OwnerID string `json:"owner_id"` // ID of user who created it
}
