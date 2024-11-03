package domain

type Document struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	OwnerID   string `json:"owner_id"` // ID of user who created it
	IsPublic  bool   `json:"is_public"`
	CanEdit   bool   `json:"can_edit"` // For access control
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}
