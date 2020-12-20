package test

import (
	"github.com/PavelVershinin/xpg"
	"github.com/jackc/pgx/v4"
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
		"test_roles"."id",
		"test_roles"."name",     
		"test_roles"."created_at",
		"test_roles"."updated_at"
	`
}

// Connection Возвращает название подключения к БД
func (Role) Connection() (name string) {
	return "test"
}

// ScanRow Реализация чтения строки из результата запроса
func (Role) ScanRow(rows pgx.Rows) (xpg.Tabler, error) {
	row := &Role{}
	err := rows.Scan(
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
func (r *Role) Delete() error {
	return xpg.New(r).Where("id", "=", r.ID).Delete()
}

// DbTake Получение записи из БД
func (r *Role) DbTake(force ...bool) error {
	if r.ID > 0 && (!r.Valid || (len(force) > 0 && force[0])) {
		row, err := xpg.New(&Role{}).Where("id", "=", r.ID).First()
		if err != nil {
			return err
		}
		*r = *row.(*Role)
		r.Valid = true
	}
	return nil
}
