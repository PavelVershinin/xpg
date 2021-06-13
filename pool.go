package xpg

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type whereRaw struct {
	sql      string
	bindings []interface{}
}

type where struct {
	logic    string
	column   string
	operator string
	value    interface{}
	raw      whereRaw
}

type groupWhere struct {
	logic  string
	closed bool
	wheres []where
}

type distinct struct {
	active bool
	on     []string
}

type union struct {
	all  bool
	pool *Pool
}

type join struct {
	joinType  string
	table     string
	alias     string
	condition string
}

// Pool пул соединений
type Pool struct {
	pool           *pgxpool.Pool
	model          Modeler
	wheres         []groupWhere
	limit          int
	offset         int
	distinct       distinct
	groupBy        []string
	orderBy        []string
	unions         []union
	joins          []join
	migrationsPath string
}

// Close Закроет подключение к БД
func (p *Pool) Close() {
	if p != nil && p.pool != nil {
		p.pool.Close()
	}
}

func addPool(pool *pgxpool.Pool, migrationsPath string) *Pool {
	p := &Pool{}
	p.pool = pool
	p.migrationsPath = migrationsPath
	return p
}

func (p *Pool) new(model Modeler) *Pool {
	np := &Pool{}
	np.pool = p.pool
	np.migrationsPath = p.migrationsPath
	np.model = model
	return np
}
