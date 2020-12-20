package xpg

import "github.com/jackc/pgx/v4"

// Tabler интерфейс модели
type Tabler interface {
	Table() string
	Columns() string
	Connection() string
	ScanRow(rows pgx.Rows) (Tabler, error)
	Save() error
	Delete() error
}
