package task

import "context"

type Repository interface {
	List(ctx context.Context, status string) (dest []Entity, err error)
	Add(ctx context.Context, data Entity) (id string, err error)
	Get(ctx context.Context, id string) (dest Entity, err error)
	Update(ctx context.Context, id string, data Entity) (err error)
	Delete(ctx context.Context, id string) (err error)
	Mark(ctx context.Context, id string) (err error)
}
