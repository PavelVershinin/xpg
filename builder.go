package xpg

import (
	"bytes"
	"strconv"
	"strings"
)

// Join Присоединит таблицу INNER JOIN
func (p *Pool) Join(table, alias, condition string) *Pool {
	p.joins = append(p.joins, join{
		joinType:  "INNER",
		table:     table,
		alias:     alias,
		condition: condition,
	})
	return p
}

// LeftJoin Присоединит таблицу LEFT OUTER JOIN
func (p *Pool) LeftJoin(table, alias, condition string) *Pool {
	p.joins = append(p.joins, join{
		joinType:  "LEFT",
		table:     table,
		alias:     alias,
		condition: condition,
	})
	return p
}

// RightJoin Присоединит таблицу RIGHT OUTER JOIN
func (p *Pool) RightJoin(table, alias, condition string) *Pool {
	p.joins = append(p.joins, join{
		joinType:  "RIGHT",
		table:     table,
		alias:     alias,
		condition: condition,
	})
	return p
}

// FullJoin Присоединит таблицу FULL OUTER JOIN
func (p *Pool) FullJoin(table, alias, condition string) *Pool {
	p.joins = append(p.joins, join{
		joinType:  "FULL",
		table:     table,
		alias:     alias,
		condition: condition,
	})
	return p
}

// Union Объединение запросов
func (p *Pool) Union(all bool, queries ...*Pool) *Pool {
	for _, query := range queries {
		p.unions = append(p.unions, union{
			all:  all,
			pool: query,
		})
	}
	return p
}

// Limit Выбрать limit записей
func (p *Pool) Limit(limit int) *Pool {
	p.limit = limit
	return p
}

// Offset Пропустить offset записей
func (p *Pool) Offset(offset int) *Pool {
	p.offset = offset
	return p
}

// Where Добавит условие WHERE через AND
func (p *Pool) Where(column, operator string, value interface{}) *Pool {
	p.where(" AND ", column, operator, value)
	return p
}

// OrWhere Добавит условие WHERE через OR
func (p *Pool) OrWhere(column, operator string, value interface{}) *Pool {
	p.where(" OR ", column, operator, value)
	return p
}

// WhereBetween Добавит условие WHERE BETWEEN через AND
func (p *Pool) WhereBetween(column string, from, to interface{}) *Pool {
	p.GroupWhere(func(p *Pool) {
		p.Where(column, ">=", from)
		p.Where(column, "<=", to)
	})
	return p
}

// OrWhereBetween Добавит условие WHERE BETWEEN через OR
func (p *Pool) OrWhereBetween(column string, from, to interface{}) *Pool {
	p.OrGroupWhere(func(p *Pool) {
		p.Where(column, ">=", from)
		p.Where(column, "<=", to)
	})
	return p
}

// GroupWhere Добавит групповое условие WHERE через AND
func (p *Pool) GroupWhere(f func(p *Pool)) *Pool {
	var group = p.openedGroupWhere()
	if len(group.wheres) > 0 {
		group.closed = true
		group = p.openedGroupWhere()
	}
	f(p)
	group.closed = true
	return p
}

// OrGroupWhere Добавит групповое условие WHERE через OR
func (p *Pool) OrGroupWhere(f func(p *Pool)) *Pool {
	var group = p.openedGroupWhere()
	if len(group.wheres) > 0 {
		group.closed = true
		group = p.openedGroupWhere()
	}
	f(p)
	group.logic = " OR "
	group.closed = true
	return p
}

// WhereRaw Произвольное условие WHERE через AND
func (p *Pool) WhereRaw(sql string, bindings ...interface{}) *Pool {
	p.whereRaw(" AND ", whereRaw{
		sql:      sql,
		bindings: bindings,
	})
	return p
}

// OrWhereRaw Произвольное условие WHERE через OR
func (p *Pool) OrWhereRaw(sql string, bindings ...interface{}) *Pool {
	p.whereRaw(" OR ", whereRaw{
		sql:      sql,
		bindings: bindings,
	})
	return p
}

// WhereIn Добавит условие WHERE IN через AND
func (p *Pool) WhereIn(column string, values *WhereInValues) *Pool {
	p.where(" AND ", column, "IN", values)
	return p
}

// OrWhereIn Добавит условие WHERE IN через OR
func (p *Pool) OrWhereIn(column string, values *WhereInValues) *Pool {
	p.where(" OR ", column, "IN", values)
	return p
}

// WhereNotIn Добавит условие WHERE NOT IN через AND
func (p *Pool) WhereNotIn(column string, values *WhereInValues) *Pool {
	p.where(" AND ", column, "NOT IN", values)
	return p
}

// OrWhereNotIn Добавит условие WHERE NOT IN через OR
func (p *Pool) OrWhereNotIn(column string, values *WhereInValues) *Pool {
	p.where(" OR ", column, "NOT IN", values)
	return p
}

// GroupBy Группировка по колонкам
func (p *Pool) GroupBy(column string, columns ...string) *Pool {
	p.groupBy = append(p.groupBy, append([]string{column}, columns...)...)
	return p
}

// Distinct Удаление дублей
func (p *Pool) Distinct(on ...string) *Pool {
	p.distinct.active = true
	p.distinct.on = on
	return p
}

// OrderBy Отсортировать по
func (p *Pool) OrderBy(column, order string) *Pool {
	var sql bytes.Buffer
	sql.WriteString(column)
	sql.WriteString(" ")
	sql.WriteString(order)
	p.orderBy = append(p.orderBy, sql.String())
	return p
}

// OrderByRaw Произвольная сортировка
func (p *Pool) OrderByRaw(orderRaw string) *Pool {
	p.orderBy = append(p.orderBy, orderRaw)
	return p
}

// OrderByRand Отсортировать в случайном порядке
func (p *Pool) OrderByRand() *Pool {
	p.orderBy = append(p.orderBy, "RANDOM()")
	return p
}

// BuildSelect Вернёт строку запроса и аргументы
func (p *Pool) BuildSelect() (string, []interface{}) {
	var query bytes.Buffer
	from, args := p.buildFrom(nil)
	where, args := p.buildWhere(args)

	query.WriteString(p.buildSelect())
	query.WriteString(from)
	query.WriteString(p.buildJoin())
	query.WriteString(where)
	query.WriteString(p.buildGroupBy())
	query.WriteString(p.buildOrderBy())
	query.WriteString(p.buildOffset())
	query.WriteString(p.buildLimit())

	return query.String(), args
}

// BuildSum Вернёт строку запроса и аргументы
func (p *Pool) BuildSum(column string) (string, []interface{}) {
	var query bytes.Buffer
	from, args := p.buildFrom(nil)
	where, args := p.buildWhere(args)

	query.WriteString(`SELECT COALESCE(SUM("`)
	query.WriteString(column)
	query.WriteString(`"), 0)`)
	query.WriteString(from)
	query.WriteString(p.buildJoin())
	query.WriteString(where)

	return query.String(), args
}

// BuildCount Вернёт строку запроса и аргументы
func (p *Pool) BuildCount() (string, []interface{}) {
	var query bytes.Buffer
	from, args := p.buildFrom(nil)
	where, args := p.buildWhere(args)

	if p.distinct.active {
		if len(p.distinct.on) == 1 {
			query.WriteString(`SELECT COUNT(DISTINCT ` + p.distinct.on[0] + `) AS "count"`)
		} else if len(p.distinct.on) > 1 {
			query.WriteString(`SELECT COUNT(DISTINCT CONCAT(` + strings.Join(p.distinct.on, ",") + `)) AS "count"`)
		} else {
			query.WriteString(`SELECT COUNT(DISTINCT *) AS "count"`)
		}
	} else {
		query.WriteString(`SELECT COUNT(*) AS "count"`)
	}
	query.WriteString(from)
	query.WriteString(p.buildJoin())
	query.WriteString(where)

	return query.String(), args
}

func (p *Pool) buildWhere(args []interface{}) (string, []interface{}) {
	var query bytes.Buffer
	var groupNum uint
	for _, group := range p.wheres {
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
					if !strings.HasPrefix(where.column, `"`) {
						query.WriteString(`"`)
					}
					query.WriteString(where.column)
					if !strings.HasSuffix(where.column, `"`) {
						query.WriteString(`"`)
					}
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

func (p *Pool) buildSelect() string {
	var query bytes.Buffer
	query.WriteString("SELECT ")
	if p.distinct.active {
		query.WriteString("DISTINCT ")
		if len(p.distinct.on) > 0 {
			query.WriteString(`ON(`)
			query.WriteString(strings.Join(p.distinct.on, `,`))
			query.WriteString(`) `)
		}
	}
	query.WriteString(p.model.Columns())
	return query.String()
}

func (p *Pool) buildFrom(args []interface{}) (string, []interface{}) {
	var query bytes.Buffer
	query.WriteString(` FROM `)
	if len(p.unions) > 0 {
		query.WriteString("(\n")
		for i, union := range p.unions {
			sql, subArgs := union.pool.BuildSelect()
			for j, arg := range subArgs {
				args = append(args, arg)
				sql = strings.Replace(sql, "$"+strconv.Itoa(j+1), "$"+strconv.Itoa(len(args)), -1)
			}
			query.WriteString("\t ")
			if i > 0 {
				query.WriteString("UNION ")
				if union.all {
					query.WriteString("ALL ")
				}
			}
			query.WriteString(sql)
			query.WriteString("\n")
		}
		query.WriteString(`) AS "`)
	} else {
		query.WriteString(`"`)
	}
	query.WriteString(p.model.Table())
	query.WriteString(`"`)
	return query.String(), args
}

func (p *Pool) buildJoin() string {
	var query bytes.Buffer
	for _, join := range p.joins {
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

func (p *Pool) buildOffset() string {
	var query bytes.Buffer
	if p.offset > 0 {
		query.WriteString(" OFFSET ")
		query.WriteString(strconv.Itoa(p.offset))
	}
	return query.String()
}

func (p *Pool) buildLimit() string {
	var query bytes.Buffer
	if p.limit > 0 {
		query.WriteString(" LIMIT ")
		query.WriteString(strconv.Itoa(p.limit))
	}
	return query.String()
}

func (p *Pool) buildOrderBy() string {
	var query bytes.Buffer
	if len(p.orderBy) > 0 {
		query.WriteString(" ORDER BY ")
		query.WriteString(strings.Join(p.orderBy, ", "))
	}
	return query.String()
}

func (p *Pool) buildGroupBy() string {
	var query bytes.Buffer
	if len(p.groupBy) > 0 {
		query.WriteString(" GROUP BY ")
		query.WriteString(strings.Join(p.groupBy, ", "))
	}
	return query.String()
}

func (p *Pool) whereRaw(logic string, raw whereRaw) {
	var group = p.openedGroupWhere()
	group.wheres = append(group.wheres, where{
		logic: logic,
		raw:   raw,
	})
}

func (p *Pool) where(logic, column, operator string, value interface{}) {
	var group = p.openedGroupWhere()
	group.wheres = append(group.wheres, where{
		logic:    logic,
		column:   column,
		operator: operator,
		value:    value,
	})
}

func (p *Pool) openedGroupWhere() *groupWhere {
	var index = len(p.wheres) - 1
	if index == -1 || p.wheres[index].closed {
		p.wheres = append(p.wheres, groupWhere{
			logic: " AND ",
		})
		index++
	}
	return &p.wheres[index]
}
