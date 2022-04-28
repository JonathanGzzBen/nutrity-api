package models

type Category struct {
	ID       uint   `json:"id,omitempty"`
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl"`
}
