package task

import "context"

// Repository is the persistence abstraction for task aggregates.
// Implement it in infra/database.
type Repository interface {
	Create(ctx context.Context, input CreateTaskInput) (Task, error)
	GetByID(ctx context.Context, id int64) (Task, error)
}

