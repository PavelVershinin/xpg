package xpg

import (
	"bytes"
	"context"
	"errors"
	"strconv"
	"strings"
	"time"
)

// Write Запись в БД
func (p *Pool) Write(ctx context.Context, data map[string]interface{}) (int64, error) {
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
		p.Where("id", "=", updateID)
		return updateID, p.Update(ctx, data)
	}
	if len(p.wheres) > 0 {
		return -1, p.Update(ctx, data)
	}
	return p.Insert(ctx, data)
}

// Insert Вставка записи в БД
func (p *Pool) Insert(ctx context.Context, data map[string]interface{}) (id int64, err error) {
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
	sql.WriteString(p.model.Table())
	sql.WriteString(`"`)
	sql.WriteString(" (")
	sql.WriteString(`"`)
	sql.WriteString(strings.Join(columns, `","`))
	sql.WriteString(`"`)
	sql.WriteString(") VALUES (")
	sql.WriteString(strings.Join(values, ","))
	sql.WriteString(")")

	sql.WriteString(` RETURNING "id"`)
	res, err := p.Query(ctx, sql.String(), args...)
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
func (p *Pool) Update(ctx context.Context, data map[string]interface{}) error {
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

	where, args := p.buildWhere(args)

	sql.WriteString("UPDATE ")
	sql.WriteString(`"`)
	sql.WriteString(p.model.Table())
	sql.WriteString(`"`)
	sql.WriteString(" SET ")
	sql.Write(columns.Bytes())
	sql.WriteString(where)

	_, err := p.pool.Exec(ctx, sql.String(), args...)
	return err
}

// Delete Удаление записи из БД
func (p *Pool) Delete(ctx context.Context) error {
	var where, args = p.buildWhere(nil)
	var query bytes.Buffer

	query.WriteString("DELETE FROM ")
	query.WriteString(`"`)
	query.WriteString(p.model.Table())
	query.WriteString(`"`)
	query.WriteString(where)

	_, err := p.pool.Exec(ctx, query.String(), args...)
	return err
}

// Select Получить записи
func (p *Pool) Select(ctx context.Context) (*Rows, error) {
	query, args := p.BuildSelect()
	return p.Query(ctx, query, args...)
}

// Query запрос к БД
func (p *Pool) Query(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	var rows *Rows
	r, err := p.pool.Query(ctx, query, args...)
	if err == nil {
		rows = &Rows{}
		rows.pool = p
		rows.Rows = r
	}
	return rows, err
}

// First Получить первую запись
func (p *Pool) First(ctx context.Context) (Modeler, error) {
	p.limit = 1
	query, args := p.BuildSelect()
	rows, err := p.Query(ctx, query, args...)
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
func (p *Pool) Exists(ctx context.Context) (bool, error) {
	p.limit = 1
	query, args := p.BuildSelect()
	rows, err := p.Query(ctx, "SELECT EXISTS("+query+")", args...)
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
func (p *Pool) Count(ctx context.Context) (int64, error) {
	var count int64
	query, args := p.BuildCount()
	rows, err := p.Query(ctx, query, args...)
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
func (p *Pool) Sum(ctx context.Context, column string) (float64, error) {
	var sum float64
	query, args := p.BuildSum(column)
	rows, err := p.Query(ctx, query, args...)
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
