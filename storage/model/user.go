package model

import "time"

// @openapi:schema
type User struct {
	UserEditable
	ID        string     `json:"id" bson:"_id"`
	CreatedAt time.Time  `json:"createdAt" bson:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt" bson:"updatedAt"`
}

// @openapi:schema
type UserEditable struct {
	Email     string `json:"email" bson:"email" validate:"required"`
	FirstName string `json:"firstName" bson:"firstName" validate:"required"`
	LastName  string `json:"lastName" bson:"lastName" validate:"required"`
}
