package xpg

import (
	"fmt"
	"github.com/jackc/pgx"
	"log"
	"sync"
	"time"
)

var (
	mu          sync.RWMutex
	connections map[string]*Connection
)

func init() {
	connections = make(map[string]*Connection)
}

// NewConnection Создаст новое подключение к БД
func NewConnection(connectionName string, conf pgx.ConnConfig, maxConnections int, migrationsPath string) error {
	conn, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     conf,
		MaxConnections: maxConnections,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	})
	if err != nil {
		log.Fatal(err)
	}
	return AddConnection(connectionName, conn, migrationsPath)
}

// AddConnection Добавит существующее подключение в коллекцию
func AddConnection(connectionName string, conn *pgx.ConnPool, migrationsPath string) error {
	mu.Lock()
	defer mu.Unlock()
	if c, ok := connections[connectionName]; ok && c != nil {
		c.Close()
	}
	connections[connectionName] = newConn(conn, migrationsPath)
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
func DB(connectionName string) *pgx.ConnPool {
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
	_, err = conn.conn.Exec(`SET TIMEZONE='` + location.String() + `'`)
	return err
}

// Close Закроет все подключения к БД
func Close() error {
	mu.Lock()
	defer mu.Unlock()
	for _, c := range connections {
		c.Close()
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
