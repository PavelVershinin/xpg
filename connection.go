package xpg

import "github.com/jackc/pgx"

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
	conn *Connection
}

type Connection struct {
	conn           *pgx.Conn
	tabler         Tabler
	wheres         []groupWhere
	limit          int
	offset         int
	distinct       distinct
	groupBy        []string
	orderBy        []string
	unions         []union
	migrationsPath string
}

// Close Закроет подключение к БД
func (c *Connection) Close() error {
	return c.conn.Close()
}

func newConn(conn *pgx.Conn, migrationsPath string) *Connection {
	connection := &Connection{}
	connection.conn = conn
	connection.migrationsPath = migrationsPath
	return connection
}

func (c *Connection) new(tabler Tabler) (conn *Connection) {
	conn = &Connection{}
	conn.conn = c.conn
	conn.migrationsPath = c.migrationsPath
	conn.tabler = tabler
	return
}
