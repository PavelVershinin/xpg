package xpg

// Column Свойства колонки таблицы
type Column struct {
	Name    string
	Type    string
	NotNull bool
	HasDef  bool
	Num     int
}

// Enums Вернёт доступные перечисления
func (c *Connection) Enums() (map[string][]string, error) {
	var enumValues = make(map[string][]string)

	rows, err := c.Query(`
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
func (c *Connection) EnumValues(name string) ([]string, error) {
	var enums, err = c.Enums()
	if err != nil {
		return nil, err
	}
	return enums[name], nil
}

// Databases Список баз данных
func (c *Connection) Databases() ([]string, error) {
	var list []string
	rows, err := c.Query("SELECT datname FROM pg_database WHERE datistemplate = false")
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
func (c *Connection) Columns() ([]Column, error) {
	var columns []Column

	rows, err := c.Query(`
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
                            c.oid
                        FROM
                            pg_catalog.pg_class c
                        LEFT JOIN
                            pg_catalog.pg_namespace n ON n.oid = c.relnamespace
                        WHERE
                            pg_catalog.pg_table_is_visible(c.oid) AND
                            c.relname = $1
                    ) AND
                    a.attnum > 0 AND
                    NOT a.attisdropped
                ORDER BY a.attnum
    `, c.tabler.Table())
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
func (c *Connection) Tables() ([]string, error) {
	var tables []string

	rows, err := c.Query(`SELECT "tablename" FROM pg_catalog.pg_tables WHERE "schemaname"='public'`)
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
