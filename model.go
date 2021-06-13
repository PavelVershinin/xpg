package xpg

import (
	"context"
	"database/sql/driver"
	"errors"

	"github.com/PavelVershinin/xpg/xpgtypes"
	"github.com/jackc/pgx/v4"
)

// Model базовая модель соответствующая минимально требуемой структуре Modeler
type Model struct {
	ID int64
	//...
	CreatedAt xpgtypes.NullTime
	UpdatedAt xpgtypes.NullTime
}

// Table Возвращает название таблицы в базе данных
func (Model) Table() string {
	return "xpg_table"
}

// Columns Список полей, которые необходимо получать запросом SELECT
func (Model) Columns() string {
	return "*"
}

// PoolName Возвращает название подключения к БД
func (Model) PoolName() (name string) {
	return ""
}

// ScanRow Реализация чтения строки из результата запроса
func (Model) ScanRow(rows pgx.Rows) (Modeler, error) {
	row := &Model{}
	err := rows.Scan(
		&row.ID,
		//...
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	return row, err
}

// Save Сохранение новой/измененной структуры в БД
func (m *Model) Save(ctx context.Context) (err error) {
	data := map[string]interface{}{
		"id": m.ID,
		//...
	}
	m.ID, err = New(m).Write(ctx, data)
	return err
}

// Delete Удаление записи из БД
func (m *Model) Delete(ctx context.Context) error {
	return New(m).Where("id", "=", m.ID).Delete(ctx)
}

// Scan Реализация интерфейса sql.Scanner
func (m *Model) Scan(src interface{}) error {
	var ok bool
	m.ID, ok = src.(int64)
	if !ok {
		return errors.New("can't assert interface to int64")
	}
	return nil
}

// Value Реализация интерфейса driver.Valuer
func (m Model) Value() (driver.Value, error) {
	return m.ID, nil
}
