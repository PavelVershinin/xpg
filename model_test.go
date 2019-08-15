package xpg

import (
	"github.com/jackc/pgx"
	"testing"
)

type testModel struct {
	Model
	ColumnOne   string `xpg:"column_one CHAR(50) NOT NULL DEFAULT ''"`
	ColumnTwo   string `xpg:"column_two CHAR(50) NOT NULL DEFAULT ''"`
	ColumnThree int64  `xpg:"column_three BIGINT NOT NULL DEFAULT 0"`
}

// Table Возвращает название таблицы в базе данных
func (m testModel) Table() string {
	return "test_model_table"
}

// Connection Возвращает название подключения к БД
func (m *testModel) Connection() (name string) {
	return "xpg_connection"
}

// Columns Список полей, которые необходимо получать запросом SELECT
func (m testModel) Columns() string {
	return `*`
}

// Scan Реализация чтения строки из результата запроса
func (m *testModel) Scan(rows *pgx.Rows) (tabler Tabler, err error) {
	row := &testModel{}
	err = rows.Scan(
		&row.ID,
		&row.ColumnOne,
		&row.ColumnTwo,
		&row.ColumnThree,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	return row, err
}

// Save Сохранение новой/измененной структуры в БД
func (m *testModel) Save() (err error) {
	data := map[string]interface{}{
		"id":           m.ID,
		"column_one":   m.ColumnOne,
		"column_two":   m.ColumnTwo,
		"column_three": m.ColumnThree,
	}
	m.ID, err = New(m).Write(data)
	return err
}

// Delete Удаление записи из БД
func (m *testModel) Delete() (err error) {
	return New(m).Where("id", "=", m.ID).Delete()
}

/* ************************************* */
func TestModel_Columns(t *testing.T) {
	if (&testModel{}).Columns() != "*" {
		t.Errorf(`Wrong Columns, expected *, real %s`, (&testModel{}).Columns())
	}
}
func TestModel_Connection(t *testing.T) {
	if (&testModel{}).Connection() != "xpg_connection" {
		t.Errorf(`Wrong Connection, expected xpg_connection, real %s`, (&testModel{}).Connection())
	}
}
func TestModel_Table(t *testing.T) {
	if (&testModel{}).Table() != "test_model_table" {
		t.Errorf(`Wrong Table, expected test_model_table, real %s`, (&testModel{}).Table())
	}
}
func TestModel_Save(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	for i := 0; i <= 10; i++ {
		var item = &testModel{}
		item.ColumnOne = "one"
		item.ColumnTwo = "two"
		item.ColumnThree = int64(i)
		if err := item.Save(); err != nil {
			t.Fatal(err)
		}
	}
}
func TestModel_Delete(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	if rows, err := New(&testModel{}).WhereBetween("column_three", 0, 10).OrderBy("column_three", "DESC").Select(); err != nil {
		t.Error(err)
	} else {
		var items []*testModel
		for row := range rows.Fetch() {
			items = append(items, row.(*testModel))
		}
		for _, item := range items {
			if err := item.Delete(); err != nil {
				t.Fatal(err)
			}
		}
	}
}
