package model

import "time"

// @openapi:schema
type Template struct {
	TemplateEditable
	ID        string     `json:"id" bson:"_id"`
	CreatedAt time.Time  `json:"createdAt" bson:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt" bson:"updatedAt"`
}

// @openapi:schema
type TemplateEditable struct {
	// Add here your model properties, and don't forget to modify SQL request in corresponding DAO file
	Code string `json:"code" bson:"code" validate:"required"`
}
