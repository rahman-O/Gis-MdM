package domain

// Icon mirrors Java Icon for API.
type Icon struct {
	ID         *int   `json:"id,omitempty"`
	CustomerID int    `json:"customerId,omitempty"`
	Name       string `json:"name"`
	FileID     int    `json:"fileId"`
	FileName   string `json:"fileName,omitempty"`
}
