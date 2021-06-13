package xpg

import (
	"context"
	"database/sql"
	"database/sql/driver"

	"github.com/jackc/pgx/v4"
)

// Modeler интерфейс модели
type Modeler interface {
	sql.Scanner
	driver.Valuer
	Table() string
	Columns() string
	PoolName() string
	ScanRow(rows pgx.Rows) (Modeler, error)
	Save(context.Context) error
	Delete(context.Context) error
}
