package xpg

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx"
)

var (
	mu          sync.RWMutex
	connections map[string]*Connection
)

func init() {
	connections = make(map[string]*Connection)
}

// NewConnection Создаст новое подключение к БД
func NewConnection(connectionName string, connConfig *pgx.ConnConfig, migrationsPath string) error {
	ctx := context.Background()
	conn, err := pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		return fmt.Errorf("xpg: Unable to connection to database: %v\n", err)
	}
	return AddConnection(connectionName, conn, ctx, migrationsPath)
}

// AddConnection Добавит существующее подключение в коллекцию
func AddConnection(connectionName string, conn *pgx.Conn, ctx context.Context, migrationsPath string) error {
	mu.Lock()
	defer mu.Unlock()
	if c, ok := connections[connectionName]; ok && c != nil {
		if err := c.Close(); err != nil {
			return fmt.Errorf("xpg: Unable to close database connection: %v\n", err)
		}
	}
	connections[connectionName] = newConn(conn, ctx, migrationsPath)
	return nil
}

// New Вернёт подключение для работы с моделью
func New(tabler Tabler) *Connection {
	conn, err := connection(tabler.Connection())
	if err != nil {
		panic(err)
	}
	return conn.new(tabler)
}

// Conn Вернёт соединение
func Conn(connectionName string) *Connection {
	conn, err := connection(connectionName)
	if err != nil {
		panic(err)
	}
	return conn
}

// DB Вернёт нативное подключение к БД
func DB(connectionName string) *pgx.Conn {
	conn, err := connection(connectionName)
	if err != nil {
		panic(err)
	}
	return conn.conn
}

// MigrationsPath Вернёт путь к директории с миграциями для
func MigrationsPath(connectionName string) string {
	conn, err := connection(connectionName)
	if err != nil {
		panic(err)
	}
	return conn.migrationsPath
}

// SetTimezone Задаст часовой пояс
func SetTimezone(connectionName string, location *time.Location) error {
	conn, err := connection(connectionName)
	if err != nil {
		return err
	}
	_, err = conn.conn.Exec(conn.ctx, `SET TIMEZONE='`+location.String()+`'`)
	return err
}

// Close Закроет все подключения к БД
func Close() error {
	mu.Lock()
	defer mu.Unlock()
	for _, c := range connections {
		if err := c.Close(); err != nil {
			return fmt.Errorf("xpg: Unable to close database connection: %v\n", err)
		}
	}
	connections = make(map[string]*Connection)
	return nil
}

func connection(connectionName string) (c *Connection, err error) {
	mu.RLock()
	c, ok := connections[connectionName]
	if !ok {
		err = fmt.Errorf("xpg: Connection `%s` not found", connectionName)
	}
	mu.RUnlock()
	return
}
