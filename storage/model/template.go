package model

import "time"

type Template struct {
	TemplateEditable
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

type TemplateEditable struct {
	// Add here your model properties, and don't forget to modify SQL request in corresponding DAO file
	Code string `json:"code" validate:"required"`
}
