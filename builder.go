package xpg

import (
	"bytes"
	"strconv"
	"strings"
)

// Join Присоединит таблицу INNER JOIN
func (c *Connection) Join(table, alias, condition string) *Connection {
	c.joins = append(c.joins, join{
		joinType:  "INNER",
		table:     table,
		alias:     alias,
		condition: condition,
	})
	return c
}

// LeftJoin Присоединит таблицу LEFT OUTER JOIN
func (c *Connection) LeftJoin(table, alias, condition string) *Connection {
	c.joins = append(c.joins, join{
		joinType:  "LEFT",
		table:     table,
		alias:     alias,
		condition: condition,
	})
	return c
}

// RightJoin Присоединит таблицу RIGHT OUTER JOIN
func (c *Connection) RightJoin(table, alias, condition string) *Connection {
	c.joins = append(c.joins, join{
		joinType:  "RIGHT",
		table:     table,
		alias:     alias,
		condition: condition,
	})
	return c
}

// FullJoin Присоединит таблицу FULL OUTER JOIN
func (c *Connection) FullJoin(table, alias, condition string) *Connection {
	c.joins = append(c.joins, join{
		joinType:  "FULL",
		table:     table,
		alias:     alias,
		condition: condition,
	})
	return c
}

// Union Объединение запросов
func (c *Connection) Union(all bool, queries ...*Connection) *Connection {
	for _, query := range queries {
		c.unions = append(c.unions, union{
			all:  all,
			conn: query,
		})
	}
	return c
}

// Limit Выбрать limit записей
func (c *Connection) Limit(limit int) *Connection {
	c.limit = limit
	return c
}

// Offset Пропустить offset записей
func (c *Connection) Offset(offset int) *Connection {
	c.offset = offset
	return c
}

// Where Добавит условие WHERE через AND
func (c *Connection) Where(column, operator string, value interface{}) *Connection {
	c.where(" AND ", column, operator, value)
	return c
}

// OrWhere Добавит условие WHERE через OR
func (c *Connection) OrWhere(column, operator string, value interface{}) *Connection {
	c.where(" OR ", column, operator, value)
	return c
}

// WhereBetween Добавит условие WHERE BETWEEN через AND
func (c *Connection) WhereBetween(column string, from, to interface{}) *Connection {
	c.GroupWhere(func(c *Connection) {
		c.Where(column, ">=", from)
		c.Where(column, "<=", to)
	})
	return c
}

// OrWhereBetween Добавит условие WHERE BETWEEN через OR
func (c *Connection) OrWhereBetween(column string, from, to interface{}) *Connection {
	c.OrGroupWhere(func(c *Connection) {
		c.Where(column, ">=", from)
		c.Where(column, "<=", to)
	})
	return c
}

// GroupWhere Добавит групповое условие WHERE через AND
func (c *Connection) GroupWhere(f func(c *Connection)) *Connection {
	var group = c.openedGroupWhere()
	if len(group.wheres) > 0 {
		group.closed = true
		group = c.openedGroupWhere()
	}
	f(c)
	group.closed = true
	return c
}

// OrGroupWhere Добавит групповое условие WHERE через OR
func (c *Connection) OrGroupWhere(f func(c *Connection)) *Connection {
	var group = c.openedGroupWhere()
	if len(group.wheres) > 0 {
		group.closed = true
		group = c.openedGroupWhere()
	}
	f(c)
	group.logic = " OR "
	group.closed = true
	return c
}

// WhereRaw Произвольное условие WHERE через AND
func (c *Connection) WhereRaw(sql string, bindings ...interface{}) *Connection {
	c.whereRaw(" AND ", whereRaw{
		sql:      sql,
		bindings: bindings,
	})
	return c
}

// OrWhereRaw Произвольное условие WHERE через OR
func (c *Connection) OrWhereRaw(sql string, bindings ...interface{}) *Connection {
	c.whereRaw(" OR ", whereRaw{
		sql:      sql,
		bindings: bindings,
	})
	return c
}

// WhereIn Добавит условие WHERE IN через AND
func (c *Connection) WhereIn(column string, values *WhereInValues) *Connection {
	c.where(" AND ", column, "IN", values)
	return c
}

// OrWhereIn Добавит условие WHERE IN через OR
func (c *Connection) OrWhereIn(column string, values *WhereInValues) *Connection {
	c.where(" OR ", column, "IN", values)
	return c
}

// WhereNotIn Добавит условие WHERE NOT IN через AND
func (c *Connection) WhereNotIn(column string, values *WhereInValues) *Connection {
	c.where(" AND ", column, "NOT IN", values)
	return c
}

// OrWhereNotIn Добавит условие WHERE NOT IN через OR
func (c *Connection) OrWhereNotIn(column string, values *WhereInValues) *Connection {
	c.where(" OR ", column, "NOT IN", values)
	return c
}

// GroupBy Группировка по колонкам
func (c *Connection) GroupBy(column string, columns ...string) *Connection {
	for _, column := range append([]string{column}, columns...) {
		var sql bytes.Buffer
		sql.WriteString(column)
		c.groupBy = append(c.groupBy, sql.String())
	}
	return c
}

// Distinct Удаление дублей
func (c *Connection) Distinct(on ...string) *Connection {
	c.distinct.active = true
	c.distinct.on = on
	return c
}

// OrderBy Отсортировать по
func (c *Connection) OrderBy(column, order string) *Connection {
	var sql bytes.Buffer
	sql.WriteString(column)
	sql.WriteString(" ")
	sql.WriteString(order)
	c.orderBy = append(c.orderBy, sql.String())
	return c
}

// OrderByRaw Произвольная сортировка
func (c *Connection) OrderByRaw(orderRaw string) *Connection {
	c.orderBy = append(c.orderBy, orderRaw)
	return c
}

// OrderByRand Отсортировать в случайном порядке
func (c *Connection) OrderByRand() *Connection {
	c.orderBy = append(c.orderBy, "RANDOM()")
	return c
}

// BuildSelect Вернёт строку запроса и аргументы
func (c *Connection) BuildSelect() (string, []interface{}) {
	var query bytes.Buffer
	from, args := c.buildFrom(nil)
	where, args := c.buildWhere(args)

	query.WriteString(c.buildSelect())
	query.WriteString(from)
	query.WriteString(c.buildJoin())
	query.WriteString(where)
	query.WriteString(c.buildGroupBy())
	query.WriteString(c.buildOrderBy())
	query.WriteString(c.buildOffset())
	query.WriteString(c.buildLimit())

	return query.String(), args
}

// BuildSum Вернёт строку запроса и аргументы
func (c *Connection) BuildSum(column string) (string, []interface{}) {
	var query bytes.Buffer
	from, args := c.buildFrom(nil)
	where, args := c.buildWhere(args)

	query.WriteString(`SELECT COALESCE(SUM("`)
	query.WriteString(column)
	query.WriteString(`"), 0)`)
	query.WriteString(from)
	query.WriteString(c.buildJoin())
	query.WriteString(where)

	return query.String(), args
}

// BuildCount Вернёт строку запроса и аргументы
func (c *Connection) BuildCount() (string, []interface{}) {
	var query bytes.Buffer
	from, args := c.buildFrom(nil)
	where, args := c.buildWhere(args)

	if c.distinct.active {
		if len(c.distinct.on) == 1 {
			query.WriteString(`SELECT COUNT(DISTINCT ` + c.distinct.on[0] + `) AS "count"`)
		} else if len(c.distinct.on) > 1 {
			query.WriteString(`SELECT COUNT(DISTINCT CONCAT(` + strings.Join(c.distinct.on, ",") + `)) AS "count"`)
		} else {
			query.WriteString(`SELECT COUNT(DISTINCT *) AS "count"`)
		}
	} else {
		query.WriteString(`SELECT COUNT(*) AS "count"`)
	}
	query.WriteString(from)
	query.WriteString(c.buildJoin())
	query.WriteString(where)

	return query.String(), args
}

func (c *Connection) buildWhere(args []interface{}) (string, []interface{}) {
	var query bytes.Buffer
	var groupNum uint
	for _, group := range c.wheres {
		if len(group.wheres) > 0 {
			if groupNum > 0 {
				query.WriteString(group.logic)
			}
			query.WriteString("(")
			for i, where := range group.wheres {
				if i > 0 {
					query.WriteString(where.logic)
				}
				if where.raw.sql == "" {
					query.WriteString(`"`)
					query.WriteString(where.column)
					query.WriteString(`"`)
					switch where.operator {
					case "IN":
						var sql string
						sql, args = where.value.(*WhereInValues).Sql(args)
						query.WriteString(sql)
					case "NOT IN":
						var sql string
						sql, args = where.value.(*WhereInValues).Sql(args)
						query.WriteString(" NOT")
						query.WriteString(sql)
					default:
						args = append(args, where.value)
						query.WriteString(where.operator)
						query.WriteString("$")
						query.WriteString(strconv.Itoa(len(args)))
					}
				} else {
					sql := where.raw.sql
					for i, arg := range where.raw.bindings {
						args = append(args, arg)
						sql = strings.Replace(sql, "$"+strconv.Itoa(i+1), "$"+strconv.Itoa(len(args)), -1)
					}
					query.WriteString("(")
					query.WriteString(sql)
					query.WriteString(")")
				}
			}
			query.WriteString(")")
			groupNum++
		}
	}

	if query.Len() > 0 {
		return " WHERE " + query.String(), args
	}

	return "", args
}

func (c *Connection) buildSelect() string {
	var query bytes.Buffer
	query.WriteString("SELECT ")
	if c.distinct.active {
		query.WriteString("DISTINCT ")
		if len(c.distinct.on) > 0 {
			query.WriteString(`ON(`)
			query.WriteString(strings.Join(c.distinct.on, `,`))
			query.WriteString(`) `)
		}
	}
	query.WriteString(c.tabler.Columns())
	return query.String()
}

func (c *Connection) buildFrom(args []interface{}) (string, []interface{}) {
	var query bytes.Buffer
	query.WriteString(` FROM `)
	if len(c.unions) > 0 {
		query.WriteString("(\n")
		for i, union := range c.unions {
			sql, subArgs := union.conn.BuildSelect()
			for j, arg := range subArgs {
				args = append(args, arg)
				sql = strings.Replace(sql, "$"+strconv.Itoa(j+1), "$"+strconv.Itoa(len(args)), -1)
			}
			query.WriteString("\t")
			if i > 0 {
				query.WriteString("UNION ")
				if union.all {
					query.WriteString("ALL ")
				}
			}
			query.WriteString(sql)
			query.WriteString("\n")
		}
		query.WriteString(`) AS "xpg_union_`)
	} else {
		query.WriteString(`"`)
	}
	query.WriteString(c.tabler.Table())
	query.WriteString(`"`)
	return query.String(), args
}

func (c *Connection) buildJoin() string {
	var query bytes.Buffer
	for _, join := range c.joins {
		query.WriteString(" ")
		query.WriteString(join.joinType)
		query.WriteString(" JOIN ")
		query.WriteString(join.table)
		if join.alias != "" {
			query.WriteString(" AS ")
			query.WriteString(join.alias)
		}
		if join.condition != "" {
			query.WriteString(" ON(")
			query.WriteString(join.condition)
			query.WriteString(")")
		}
	}
	return query.String()
}

func (c *Connection) buildOffset() string {
	var query bytes.Buffer
	if c.offset > 0 {
		query.WriteString(" OFFSET ")
		query.WriteString(strconv.Itoa(c.offset))
	}
	return query.String()
}

func (c *Connection) buildLimit() string {
	var query bytes.Buffer
	if c.limit > 0 {
		query.WriteString(" LIMIT ")
		query.WriteString(strconv.Itoa(c.limit))
	}
	return query.String()
}

func (c *Connection) buildOrderBy() string {
	var query bytes.Buffer
	if len(c.orderBy) > 0 {
		query.WriteString(" ORDER BY ")
		query.WriteString(strings.Join(c.orderBy, ", "))
	}
	return query.String()
}

func (c *Connection) buildGroupBy() string {
	var query bytes.Buffer
	if len(c.groupBy) > 0 {
		query.WriteString(" GROUP BY ")
		query.WriteString(strings.Join(c.groupBy, ", "))
	}
	return query.String()
}

func (c *Connection) whereRaw(logic string, raw whereRaw) {
	var group = c.openedGroupWhere()
	group.wheres = append(group.wheres, where{
		logic: logic,
		raw:   raw,
	})
}

func (c *Connection) where(logic, column, operator string, value interface{}) {
	var group = c.openedGroupWhere()
	group.wheres = append(group.wheres, where{
		logic:    logic,
		column:   column,
		operator: operator,
		value:    value,
	})
}

func (c *Connection) openedGroupWhere() *groupWhere {
	var index = len(c.wheres) - 1
	if index == -1 || c.wheres[index].closed {
		c.wheres = append(c.wheres, groupWhere{
			logic: " AND ",
		})
		index++
	}
	return &c.wheres[index]
}
