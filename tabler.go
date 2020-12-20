package xpg

import (
	"database/sql"
	"database/sql/driver"

	"github.com/jackc/pgx/v4"
)

// Tabler интерфейс модели
type Tabler interface {
	sql.Scanner
	driver.Valuer
	Table() string
	Columns() string
	Connection() string
	ScanRow(rows pgx.Rows) (Tabler, error)
	Save() error
	Delete() error
}
