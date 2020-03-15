package xpg

import "github.com/jackc/pgx"

// Tabler интерфейс модели
type Tabler interface {
	Table() string
	Columns() string
	Connection() string
	Scan(rows pgx.Rows) (Tabler, error)
	Save() error
	Delete() error
}
