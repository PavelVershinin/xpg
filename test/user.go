package test

import (
	"github.com/PavelVershinin/xpg"
	"github.com/jackc/pgx"
)

type User struct {
	xpg.Model
	FirstName  string `xpg:"first_name VARCHAR(50) NOT NULL DEFAULT ''"`
	SecondName string `xpg:"second_name VARCHAR(50) NOT NULL DEFAULT ''"`
	LastName   string `xpg:"last_name VARCHAR(50) NOT NULL DEFAULT ''"`
	Email      string `xpg:"email VARCHAR(254) NOT NULL DEFAULT ''"`
	Phone      string `xpg:"phone VARCHAR(18) NOT NULL DEFAULT ''"`
	RoleID     int64  `xpg:"role_id BIGINT NOT NULL DEFAULT 0"`
	Balance    int64  `xpg:"balance BIGINT NOT NULL DEFAULT 0"`
}

// Table Возвращает название таблицы в базе данных
func (u User) Table() string {
	return "test_users"
}

// Columns Список полей, которые необходимо получать запросом SELECT
func (u User) Columns() string {
	return `
		"test_users"."id",
		"test_users"."first_name",
		"test_users"."second_name",
		"test_users"."last_name",
		"test_users"."email",
		"test_users"."phone",
		"test_users"."role_id",     
		"test_users"."balance",     
		"test_users"."created_at",
		"test_users"."updated_at"
	`
}

// Connection Возвращает название подключения к БД
func (u *User) Connection() (name string) {
	return "test"
}

// Scan Реализация чтения строки из результата запроса
func (u *User) Scan(rows pgx.Rows) (tabler xpg.Tabler, err error) {
	row := &User{}
	err = rows.Scan(
		&row.ID,
		&row.FirstName,
		&row.SecondName,
		&row.LastName,
		&row.Email,
		&row.Phone,
		&row.RoleID,
		&row.Balance,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	return row, err
}

// Save Сохранение новой/измененной структуры в БД
func (u *User) Save() (err error) {
	data := map[string]interface{}{
		"id":          u.ID,
		"first_name":  u.FirstName,
		"second_name": u.SecondName,
		"last_name":   u.LastName,
		"email":       u.Email,
		"phone":       u.Phone,
		"role_id":     u.RoleID,
		"balance":     u.Balance,
	}
	u.ID, err = xpg.New(u).Write(data)
	return err
}

// Delete Удаление записи из БД
func (u *User) Delete() (err error) {
	return xpg.New(u).Where("id", "=", u.ID).Delete()
}
