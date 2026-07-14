package database

import (
	"database/sql"
	"fmt"
	"project/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type TaskStor struct {
	db *sqlx.DB
}

func NewTaskStore(db *sqlx.DB) *TaskStor {
	return &TaskStor{db: db}
}

func (s *TaskStor) GetAll() ([]models.Task, error) {
	var tasks []models.Task

	query := `
 SELECT id, title, description, completed, created_at, updated_at 
 FROM tasks 
 ORDER BY created_at DESC;` // Исправили ORDER BY и DESC

	err := s.db.Select(&tasks, query)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *TaskStor) GetByID(id int) (*models.Task, error) {
	var task models.Task

	query := `
 SELECT id, title, description, completed, created_at, updated_at 
 FROM tasks 
 WHERE id = $1;`

	err := s.db.Get(&task, query, id)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task with id %d not found", id) // Добавили кавычки
	}

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *TaskStor) Create(input models.CreateTaskInput) (*models.Task, error) {
	var task models.Task

	query := `
 INSERT INTO tasks (title, description, completed, created_at, updated_at)
 VALUES ($1, $2, $3, $4, $5)
 RETURNING id, title, description, completed, created_at, updated_at;` // Убрали ";" перед RETURNING

	now := time.Now()

	err := s.db.QueryRowx(query, input.Title, input.Description, input.Completed, now, now).StructScan(&task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *TaskStor) Update(id int, input models.UpdateTaskInput) (*models.Task, error) {
	task, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		task.Title = *input.Title
	}

	if input.Description != nil {
		task.Description = *input.Description
	}

	if input.Completed != nil {
		task.Completed = *input.Completed
	}

	task.UpdatedAt = time.Now()

	query := `
 UPDATE tasks
 SET title = $1, description = $2, completed = $3, updated_at = $4
 WHERE id = $5
 RETURNING id, title, description, completed, created_at, updated_at;` // Убрали ";" и поправили плейсхолдеры

	var updateTask models.Task

	// Передаем task.ID пятым аргументом под $5
	err = s.db.QueryRowx(query, task.Title, task.Description, task.Completed, task.UpdatedAt, task.ID).StructScan(&updateTask)
	if err != nil {
		return nil, err
	}

	return &updateTask, nil
}

func (s *TaskStor) Delete(id int) error {
	query := "DELETE FROM tasks WHERE id = $1;" // Добавили кавычки

	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("task with id %d not found", id) // Добавили кавычки
	}

	return nil
}
