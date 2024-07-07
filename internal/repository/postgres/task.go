package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"medium/internal/domain/task"
	"medium/pkg/store"
	"strings"
)

type TaskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) List(ctx context.Context, status string) (dest []task.Entity, err error) {
	query := `
		SELECT id, title, active_at
		FROM tasks
		WHERE status = $1 AND active_at <= CURRENT_DATE
		ORDER BY created_at`
	err = r.db.SelectContext(ctx, &dest, query, status)

	return
}

func (r *TaskRepository) Add(ctx context.Context, data task.Entity) (id string, err error) {
	query := `INSERT INTO tasks (title, active_at)
VALUES ($1, $2)
RETURNING id;`
	args := []any{data.Title, data.ActiveAt}

	if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}
	return
}

func (r *TaskRepository) Get(ctx context.Context, id string) (dest task.Entity, err error) {
	query := `SELECT id, title, active_at FROM tasks WHERE id = $1`

	args := []any{id}
	if err = r.db.GetContext(ctx, &dest, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}
	return
}
func (r *TaskRepository) Update(ctx context.Context, id string, data task.Entity) (err error) {
	sets, args := r.prepareArgs(data)
	if len(args) > 0 {

		args = append(args, id)
		sets = append(sets, "updated_at=CURRENT_TIMESTAMP")
		query := fmt.Sprintf("UPDATE tasks SET %s WHERE id=$%d RETURNING id", strings.Join(sets, ", "), len(args))

		if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = store.ErrorNotFound
			}
		}
	}

	return
}

func (r *TaskRepository) prepareArgs(data task.Entity) (sets []string, args []any) {
	if data.Title != nil {
		args = append(args, data.Title)
		sets = append(sets, fmt.Sprintf("title=$%d", len(args)))
	}

	if data.ActiveAt != nil {
		args = append(args, data.ActiveAt)
		sets = append(sets, fmt.Sprintf("active_at=$%d", len(args)))
	}
	return
}

func (r *TaskRepository) Delete(ctx context.Context, id string) (err error) {
	query := `
		DELETE FROM tasks
		WHERE id=$1
		RETURNING id`

	args := []any{id}

	if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		return
	}
	return
}

func (r *TaskRepository) Mark(ctx context.Context, id string) (err error) {
	query := `
			UPDATE tasks
			SET status = 'done', 
			updated_at = CURRENT_TIMESTAMP
			WHERE id = $1 
			RETURNING id;
`
	args := []any{id}

	if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}
	return
}
