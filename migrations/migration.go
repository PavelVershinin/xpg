package migrations

import (
	"context"

	"github.com/PavelVershinin/xpg"
	"github.com/jackc/pgx/v4"
)

type migration struct {
	xpg.Model
	File     string `xpg:"file VARCHAR(25) NOT NULL DEFAULT ''"`
	poolName string
}

// Table Возвращает название таблицы в базе данных
func (migration) Table() string {
	return "xpg_migrations"
}

// Columns Список полей, которые необходимо получать запросом SELECT
func (migration) Columns() string {
	return `
		"id",
		"file",
		"created_at",
		"updated_at"
	`
}

// PoolName Возвращает название подключения к БД
func (m migration) PoolName() string {
	return m.poolName
}

func (m *migration) SetPool(name string) {
	m.poolName = name
}

// ScanRow Реализация чтения строки из результата запроса
func (m *migration) ScanRow(rows pgx.Rows) (xpg.Modeler, error) {
	row := &migration{}
	err := rows.Scan(
		&row.ID,
		&row.File,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	return row, err
}

// Save Сохранение новой/измененной структуры в БД
func (m *migration) Save(ctx context.Context) error {
	var err error
	data := map[string]interface{}{
		"id":   m.ID,
		"file": m.File,
	}
	m.ID, err = xpg.New(m).Write(ctx, data)
	return err
}
