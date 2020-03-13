package xpg

import (
	"github.com/PavelVershinin/xpg/xpgtypes"
	"github.com/jackc/pgx"
)

// Model базовая модель соответствующая минимально требуемой структуре Tabler
type Model struct {
	ID int64 `json:"id"`
	//...
	CreatedAt xpgtypes.NullTime `json:"created_at"`
	UpdatedAt xpgtypes.NullTime `json:"updated_at"`
}

// Table Возвращает название таблицы в базе данных
func (m Model) Table() string {
	return "xpg_table"
}

// Columns Список полей, которые необходимо получать запросом SELECT
func (m Model) Columns() string {
	return "*"
}

// Connection Возвращает название подключения к БД
func (m *Model) Connection() (name string) {
	return ""
}

// Scan Реализация чтения строки из результата запроса
func (m *Model) Scan(rows pgx.Rows) (tabler Tabler, err error) {
	row := &Model{}
	err = rows.Scan(
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
func (m *Model) Delete() (err error) {
	return New(m).Where("id", "=", m.ID).Delete()
}
