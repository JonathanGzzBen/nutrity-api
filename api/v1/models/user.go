package models

type User struct {
	// Auth
	ID          uint   `json:"id,omitempty"`
	GoogleToken string `json:"-"`
	AccessToken string `json:"-"`
	// Data
	Username          string `json:"username"`
	Email             string `json:"email"`
	FirstName         string `json:"firstname"`
	LastName          string `json:"lastname"`
	UserProfileEdited bool   `json:"userProfileEdited"`
	Calories          uint   `json:"calories"`
	Carbs             uint   `json:"carbs"`
	Day               uint   `json:"day"`
	Fats              uint   `json:"fats"`
	Proteins          uint   `json:"proteins"`
	RecipesAdded      string `json:"recipesAdded"` // List of recipes divided by character '^'
}
