package models

import "time"

type Task struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"` // Для TEXT в Postgres используется string
	Completed   bool      `json:"completed" db:"completed"`     // Для BOOLEAN используется bool
	CreatedAt   time.Time `json:"created_at" db:"created_at"`   // Для TIMESTAMPTZ используется time.Time
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`   // Для TIMESTAMPTZ используется time.Time
}

type CreateTaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type UpdateTaskInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}
