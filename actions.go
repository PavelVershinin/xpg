package xpg

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"time"
)

// Write Запись в БД
func (c *Connection) Write(data map[string]interface{}) (id int64, err error) {
	var updateID int64
	if id, ok := data["id"]; ok {
		switch aId := id.(type) {
		case int64:
			updateID = aId
		case int32:
			updateID = int64(aId)
		case int:
			updateID = int64(aId)
		case uint64:
			updateID = int64(aId)
		case uint32:
			updateID = int64(aId)
		case uint:
			updateID = int64(aId)
		default:
			return 0, errors.New("xpg: Unsupported type of column id")
		}
		delete(data, "id")
	}
	if updateID > 0 {
		c.Where("id", "=", updateID)
		return updateID, c.Update(data)
	}
	if len(c.wheres) > 0 {
		return -1, c.Update(data)
	}
	return c.Insert(data)
}

// Insert Вставка записи в БД
func (c *Connection) Insert(data map[string]interface{}) (id int64, err error) {
	var columns = make([]string, 0, len(data)+1)
	var values = make([]string, 0, len(data)+1)
	var args = make([]interface{}, 0, len(data)+1)

	for column, val := range data {
		if column != "id" {
			args = append(args, val)
			columns = append(columns, column)
			values = append(values, `$`+strconv.Itoa(len(args)))
		}
	}

	if _, ok := data["created_at"]; !ok {
		args = append(args, time.Now())
		columns = append(columns, "created_at")
		values = append(values, `$`+strconv.Itoa(len(args)))
	}

	var sql bytes.Buffer
	sql.WriteString("INSERT INTO ")
	sql.WriteString(`"`)
	sql.WriteString(c.tabler.Table())
	sql.WriteString(`"`)
	sql.WriteString(" (")
	sql.WriteString(`"`)
	sql.WriteString(strings.Join(columns, `","`))
	sql.WriteString(`"`)
	sql.WriteString(") VALUES (")
	sql.WriteString(strings.Join(values, ","))
	sql.WriteString(")")

	sql.WriteString(` RETURNING "id"`)
	res, err := c.Query(sql.String(), args...)
	if err != nil {
		return 0, err
	}
	if res.Next() {
		err = res.Scan(&id)
	}
	res.Close()

	return id, err
}

// Update Изменение записи в БД
func (c *Connection) Update(data map[string]interface{}) (err error) {
	var columns, sql bytes.Buffer
	var args = make([]interface{}, 0, len(data)+1)

	for column, val := range data {
		args = append(args, val)
		columns.WriteString(`"`)
		columns.WriteString(column)
		columns.WriteString(`"`)
		columns.WriteString("=$")
		columns.WriteString(strconv.Itoa(len(args)))
		columns.WriteString(",")
	}

	args = append(args, time.Now())
	columns.WriteString(`"`)
	columns.WriteString("updated_at")
	columns.WriteString(`"`)
	columns.WriteString("=$")
	columns.WriteString(strconv.Itoa(len(args)))

	where, args := c.buildWhere(args)

	sql.WriteString("UPDATE ")
	sql.WriteString(`"`)
	sql.WriteString(c.tabler.Table())
	sql.WriteString(`"`)
	sql.WriteString(" SET ")
	sql.Write(columns.Bytes())
	sql.WriteString(where)

	_, err = c.conn.Exec(c.ctx, sql.String(), args...)
	return err
}

// Delete Удаление записи из БД
func (c *Connection) Delete() (err error) {
	var where, args = c.buildWhere(nil)
	var query bytes.Buffer

	query.WriteString("DELETE FROM ")
	query.WriteString(`"`)
	query.WriteString(c.tabler.Table())
	query.WriteString(`"`)
	query.WriteString(where)

	_, err = c.conn.Exec(c.ctx, query.String(), args...)
	return err
}

// Select Получить записи
func (c *Connection) Select() (rows *Rows, err error) {
	query, args := c.BuildSelect()
	return c.Query(query, args...)
}

// Query запрос к БД
func (c *Connection) Query(query string, args ...interface{}) (rows *Rows, err error) {
	r, err := c.conn.Query(c.ctx, query, args...)
	if err == nil {
		rows = &Rows{}
		rows.conn = c
		rows.Rows = r
	}
	return
}

// First Получить первую запись
func (c *Connection) First() (Tabler, error) {
	c.limit = 1
	query, args := c.BuildSelect()
	rows, err := c.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, errors.New("xpg: No records found")
	}
	row, err := rows.Get()
	return row, err
}

// Exists Проверка наличия записи в базе
func (c *Connection) Exists() (bool, error) {
	c.limit = 1
	query, args := c.BuildSelect()
	rows, err := c.Query("SELECT EXISTS("+query+")", args...)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	if rows.Next() {
		var buff bool
		err = rows.Scan(&buff)
		return buff, err
	}
	return false, nil
}

// Count Получить количество записей
func (c *Connection) Count() (int64, error) {
	var count int64
	query, args := c.BuildCount()
	rows, err := c.Query(query, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return 0, err
		}
	}
	return count, err
}

// Sum Получить сумму записей
func (c *Connection) Sum(column string) (float64, error) {
	var sum float64
	query, args := c.BuildSum(column)
	rows, err := c.Query(query, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if rows.Next() {
		if err = rows.Scan(&sum); err != nil {
			return 0, err
		}
	}
	return sum, err
}
