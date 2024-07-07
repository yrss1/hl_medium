package task

import "time"

type Entity struct {
	ID       string     `db:"id"`
	Title    *string    `db:"title"`
	ActiveAt *time.Time `db:"active_at"`
}
