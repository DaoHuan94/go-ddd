package repository_impl

import (
	"context"
	"errors"
	"sync"

	"go-ddd/domain/task"

	"gorm.io/gorm"
)

var ErrTaskNotFound = errors.New("task not found")

// InMemoryTaskRepository is a simple repository implementation for development.
// Replace with a real DB-backed implementation (GORM/pg/sql) later.
type InMemoryTaskRepository struct {
	mu     sync.RWMutex
	nextID int64
	data   map[int64]task.Task
}

func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		nextID: 1,
		data:   make(map[int64]task.Task),
	}
}

func (r *InMemoryTaskRepository) Create(ctx context.Context, input task.CreateTaskInput) (task.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := r.nextID
	r.nextID++

	t := task.Task{
		ID:          id,
		ProjectID:   input.ProjectID,
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
	}
	r.data[id] = t

	return t, nil
}

func (r *InMemoryTaskRepository) GetByID(ctx context.Context, id int64) (task.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.data[id]
	if !ok {
		return task.Task{}, ErrTaskNotFound
	}
	return t, nil
}

// PostgresTaskRepository is a real DB implementation backed by PostgreSQL via GORM.
type PostgresTaskRepository struct {
	db *gorm.DB
}

type taskRow struct {
	ID          int64   `gorm:"column:id;primaryKey"`
	ProjectID   int64   `gorm:"column:project_id"`
	Title       string  `gorm:"column:title"`
	Description *string `gorm:"column:description"`
	Priority    *string `gorm:"column:priority"`
}

func NewPostgresTaskRepository(db *gorm.DB) *PostgresTaskRepository {
	return &PostgresTaskRepository{db: db}
}

func (r *PostgresTaskRepository) Create(ctx context.Context, input task.CreateTaskInput) (task.Task, error) {
	desc := nullableString(input.Description)
	priority := nullableString(input.Priority)

	row := taskRow{
		ProjectID:   input.ProjectID,
		Title:       input.Title,
		Description: desc,
		Priority:    priority,
	}

	// Select only columns we want so GORM doesn't accidentally write zero values to other columns.
	if err := r.db.WithContext(ctx).
		Table("tasks").
		Select("project_id", "title", "description", "priority").
		Create(&row).Error; err != nil {
		return task.Task{}, err
	}

	return task.Task{
		ID:          row.ID,
		ProjectID:   row.ProjectID,
		Title:       row.Title,
		Description: derefString(row.Description),
		Priority:    derefString(row.Priority),
	}, nil
}

func (r *PostgresTaskRepository) GetByID(ctx context.Context, id int64) (task.Task, error) {
	var row taskRow
	err := r.db.WithContext(ctx).
		Table("tasks").
		First(&row, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return task.Task{}, ErrTaskNotFound
		}
		return task.Task{}, err
	}

	return task.Task{
		ID:          row.ID,
		ProjectID:   row.ProjectID,
		Title:       row.Title,
		Description: derefString(row.Description),
		Priority:    derefString(row.Priority),
	}, nil
}

func nullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func derefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

