package xpg

import "context"

// Column Свойства колонки таблицы
type Column struct {
	Name    string
	Type    string
	NotNull bool
	HasDef  bool
	Num     int
}

// Enums Вернёт доступные перечисления
func (p *Pool) Enums(ctx context.Context) (map[string][]string, error) {
	var enumValues = make(map[string][]string)

	rows, err := p.Query(ctx, `
		SELECT
			t.typname,  
			e.enumlabel
		FROM 
			pg_type t 
   		JOIN 
			pg_enum e ON t.oid = e.enumtypid  
   		JOIN 
			pg_catalog.pg_namespace n ON n.oid = t.typnamespace
	`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var name, value string
		if err := rows.Scan(&name, &value); err == nil {
			enumValues[name] = append(enumValues[name], value)
		}
	}

	rows.Close()

	return enumValues, nil
}

// EnumValues Вернёт доступные значения для типа ENUM
func (p *Pool) EnumValues(ctx context.Context, name string) ([]string, error) {
	var enums, err = p.Enums(ctx)
	if err != nil {
		return nil, err
	}
	return enums[name], nil
}

// Databases Список баз данных
func (p *Pool) Databases(ctx context.Context) ([]string, error) {
	var list []string
	rows, err := p.Query(ctx, "SELECT datname FROM pg_database WHERE datistemplate = false")
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var buff string
		if err := rows.Scan(&buff); err != nil {
			return list, err
		}
		list = append(list, buff)
	}
	return list, nil
}

// Columns Вернёт список колонок текущей таблицы
func (p *Pool) Columns(ctx context.Context) ([]Column, error) {
	var columns []Column

	rows, err := p.Query(ctx, `
                SELECT
                    a.attname,
                    pg_catalog.format_type(a.atttypid, a.atttypmod),
                    a.attnotnull,
                    a.atthasdef,
                    a.attnum
                FROM
                    pg_catalog.pg_attribute a
                WHERE
                    a.attrelid = (
                        SELECT
                            p.oid
                        FROM
                            pg_catalog.pg_class c
                        LEFT JOIN
                            pg_catalog.pg_namespace n ON n.oid = p.relnamespace
                        WHERE
                            pg_catalog.pg_table_is_visible(p.oid) AND
                            p.relname = $1
                    ) AND
                    a.attnum > 0 AND
                    NOT a.attisdropped
                ORDER BY a.attnum
    `, p.model.Table())
	if err != nil {
		return []Column{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var buff Column
		if err := rows.Scan(&buff.Name, &buff.Type, &buff.NotNull, &buff.HasDef, &buff.Num); err != nil {
			return columns, err
		}
		columns = append(columns, buff)
	}
	return columns, nil
}

// Tables Вернёт список таблиц в базе данных
func (p *Pool) Tables(ctx context.Context) ([]string, error) {
	var tables []string

	rows, err := p.Query(ctx, `SELECT "tablename" FROM pg_catalog.pg_tables WHERE "schemaname"='public'`)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var buff string
		if err := rows.Scan(&buff); err != nil {
			return tables, err
		}
		tables = append(tables, buff)
	}

	return tables, nil
}
