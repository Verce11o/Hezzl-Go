package models

import "time"

type Product struct {
	ID          int       `json:"id" db:"id"`
	ProjectID   int       `json:"projectId" db:"project_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Priority    int       `json:"priority" db:"priority"`
	Removed     bool      `json:"removed" db:"removed"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

type ProductList struct {
	Meta Meta `json:"meta"`

	Goods []Product `json:"goods"`
}

type Meta struct {
	Total   int `json:"total"`
	Removed int `json:"removed"`
	Limit   int `json:"limit"`
	Offset  int `json:"offset"`
}

type Priority struct {
	ID       int `json:"id" db:"id"`
	Priority int `json:"priority" db:"priority"`
}
