package test

import (
	"github.com/PavelVershinin/xpg"
	"github.com/jackc/pgx"
)

type Role struct {
	xpg.Model
	Name string `xpg:"name VARCHAR(50) NOT NULL DEFAULT ''"`
}

// Table Возвращает название таблицы в базе данных
func (Role) Table() string {
	return "test_roles"
}

// Columns Список полей, которые необходимо получать запросом SELECT
func (Role) Columns() string {
	return `
		"test_roles"."name"
	`
}

// Connection Возвращает название подключения к БД
func (Role) Connection() (name string) {
	return "main"
}

// Scan Реализация чтения строки из результата запроса
func (Role) Scan(rows pgx.Rows) (tabler xpg.Tabler, err error) {
	row := &Role{}
	err = rows.Scan(
		&row.ID,
		&row.Name,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	return row, err
}

// Save Сохранение новой/измененной структуры в БД
func (r *Role) Save() (err error) {
	data := map[string]interface{}{
		"id":   r.ID,
		"name": r.Name,
	}
	r.ID, err = xpg.New(r).Write(data)
	return err
}

// Delete Удаление записи из БД
func (r *Role) Delete() (err error) {
	return xpg.New(r).Where("id", "=", r.ID).Delete()
}
