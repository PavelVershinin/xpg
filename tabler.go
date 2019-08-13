package xpg

import "github.com/jackc/pgx"

type Tabler interface {
	Table() string
	Columns() string
	Connection() (name string)
	Scan(rows *pgx.Rows) (tabler Tabler, err error)
	Save() (err error)
	Delete() (err error)
}
