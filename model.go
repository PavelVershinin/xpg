package xpg

import (
	"database/sql/driver"
	"errors"

	"github.com/PavelVershinin/xpg/xpgtypes"
	"github.com/jackc/pgx/v4"
)

// Model базовая модель соответствующая минимально требуемой структуре Tabler
type Model struct {
	ID int64 `json:"id"`
	//...
	CreatedAt xpgtypes.NullTime `json:"created_at"`
	UpdatedAt xpgtypes.NullTime `json:"updated_at"`
	Valid     bool              `json:"_"`
}

// Table Возвращает название таблицы в базе данных
func (Model) Table() string {
	return "xpg_table"
}

// Columns Список полей, которые необходимо получать запросом SELECT
func (Model) Columns() string {
	return "*"
}

// Connection Возвращает название подключения к БД
func (Model) Connection() (name string) {
	return ""
}

// ScanRow Реализация чтения строки из результата запроса
func (Model) ScanRow(rows pgx.Rows) (Tabler, error) {
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
func (m *Model) Save() (err error) {
	data := map[string]interface{}{
		"id": m.ID,
		//...
	}
	m.ID, err = New(m).Write(data)
	return err
}

// Delete Удаление записи из БД
func (m *Model) Delete() error {
	return New(m).Where("id", "=", m.ID).Delete()
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

// DbTake Получение записи из БД
func (m *Model) DbTake(force ...bool) error {
	if m.ID > 0 && (!m.Valid || (len(force) > 0 && force[0])) {
		t, err := New(m).Where("id", "=", m.ID).First()
		if err != nil {
			return err
		}
		*m = *t.(*Model)
		m.Valid = true
	}
	return nil
}
