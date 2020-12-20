package migrations

import (
	"github.com/PavelVershinin/xpg"
	"github.com/jackc/pgx/v4"
)

type migration struct {
	xpg.Model
	File       string `xpg:"file VARCHAR(25) NOT NULL DEFAULT ''"`
	connection string
}

// Table Возвращает название таблицы в базе данных
func (m migration) Table() string {
	return "xpg_migrations"
}

// Columns Список полей, которые необходимо получать запросом SELECT
func (m migration) Columns() string {
	return `
		"id",
		"file",
		"created_at",
		"updated_at"
	`
}

// Connection Возвращает название подключения к БД
func (m *migration) Connection() (name string) {
	return m.connection
}

func (m *migration) SetConnection(name string) {
	m.connection = name
}

// ScanRow Реализация чтения строки из результата запроса
func (m *migration) ScanRow(rows pgx.Rows) (tabler xpg.Tabler, err error) {
	row := &migration{}
	err = rows.Scan(
		&row.ID,
		&row.File,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	return row, err
}

// Save Сохранение новой/измененной структуры в БД
func (m *migration) Save() (err error) {
	data := map[string]interface{}{
		"id":   m.ID,
		"file": m.File,
	}
	m.ID, err = xpg.New(m).Write(data)
	return err
}
