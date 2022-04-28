package models

import (
	"time"
)

type User struct {
	ID                uint      `json:"id,omitempty"`
	GoogleSub         string    `json:"-"`
	Name              string    `json:"name"`
	Birthdate         time.Time `json:"birthdate" example:"2006-01-02T15:04:05Z"`
	Gender            string    `json:"gender"`
	ProfilePictureURL string    `json:"profilePictureUrl"`
	Description       string    `json:"description"`
	ShortDescription  string    `json:"shortDescription"`
	Role              Role      `json:"role" example:"Reader"`
}

type Role string

const (
	RoleAdministrator Role = "Administrator"
	RoleWriter        Role = "Writer"
	RoleReader        Role = "Reader"
)
