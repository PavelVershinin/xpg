package test

import (
	"context"

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

func (Role) PoolName() (name string) {
	return "test"
}

// ScanRow Реализация чтения строки из результата запроса
func (Role) ScanRow(rows pgx.Rows) (xpg.Modeler, error) {
	row := &Role{}
	err := rows.Scan(
		&row.ID,
		&row.Name,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	return row, err
}

func (r *Role) Save(ctx context.Context) error {
	var err error
	data := map[string]interface{}{
		"id":   r.ID,
		"name": r.Name,
	}
	r.ID, err = xpg.New(r).Write(ctx, data)
	return err
}

// Delete Удаление записи из БД
func (r *Role) Delete(ctx context.Context) error {
	return xpg.New(r).Where("id", "=", r.ID).Delete(ctx)
}
