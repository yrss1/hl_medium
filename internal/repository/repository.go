package repository

import (
	"medium/internal/domain/task"
	"medium/internal/repository/postgres"
	"medium/pkg/store"
)

type Configuration func(r *Repository) error

type Repository struct {
	postgres store.SQLX
	Task     task.Repository
}

func New(configs ...Configuration) (s *Repository, err error) {
	s = &Repository{}

	for _, cfg := range configs {
		if err = cfg(s); err != nil {
			return
		}
	}

	return
}

func WithPostgresStore(dbName string) Configuration {
	return func(r *Repository) (err error) {
		r.postgres, err = store.New(dbName)
		if err != nil {
			return
		}
		if err = store.Migrate(dbName); err != nil {
			return
		}

		r.Task = postgres.NewTaskRepository(r.postgres.Client)

		return
	}
}
