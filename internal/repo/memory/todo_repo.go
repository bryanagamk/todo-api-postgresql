package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/bryanagamk/todo-app-postgresql/internal/domain"
)

var (
	ErrNotFound = errors.New("todo not found")
)

type TodoRepo struct {
	mu    sync.RWMutex
	seq   int64
	items map[int64]*domain.Todo
}

func NewTodoRepo() *TodoRepo {
	return &TodoRepo{items: make(map[int64]*domain.Todo)}
}

func (r *TodoRepo) Create(ctx context.Context, t *domain.Todo) (*domain.Todo, error) {
	if err := t.Validate(); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	cl := *t
	cl.ID = r.seq
	cl.CreatedAt = time.Now().UTC()
	cl.UpdatedAt = cl.CreatedAt
	r.items[cl.ID] = &cl
	return &cl, nil
}

func (r *TodoRepo) Get(ctx context.Context, id int64) (*domain.Todo, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	it, ok := r.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	cl := *it
	return &cl, nil
}
