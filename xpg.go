package xpg

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	mu    sync.RWMutex
	pools = make(map[string]*Pool)
)

// NewConnectionPool Создаст новое подключение к БД
func NewConnectionPool(ctx context.Context, poolName string, connConfig *pgxpool.Config, migrationsPath string) error {
	pool, err := pgxpool.ConnectConfig(ctx, connConfig)
	if err != nil {
		return fmt.Errorf("xpg: Unable to connection to database: %w\n", err)
	}
	return AddConnectionsPool(poolName, pool, migrationsPath)
}

// AddConnectionsPool Добавит существующее подключение в коллекцию
func AddConnectionsPool(poolName string, pool *pgxpool.Pool, migrationsPath string) error {
	mu.Lock()
	defer mu.Unlock()
	if c, ok := pools[poolName]; ok && c != nil {
		c.Close()
	}
	pools[poolName] = addPool(pool, migrationsPath)
	return nil
}

// New Вернёт подключение для работы с моделью
func New(model Modeler) *Pool {
	p, err := pool(model.PoolName())
	if err != nil {
		panic(err)
	}
	return p.new(model)
}

// DB Вернёт нативное подключение к БД
func DB(poolName string) *pgxpool.Pool {
	p, err := pool(poolName)
	if err != nil {
		panic(err)
	}
	return p.pool
}

// MigrationsPath Вернёт путь к директории с миграциями
func MigrationsPath(poolName string) string {
	p, err := pool(poolName)
	if err != nil {
		panic(err)
	}
	return p.migrationsPath
}

// SetTimezone Задаст часовой пояс
func SetTimezone(ctx context.Context, poolName string, location *time.Location) error {
	p, err := pool(poolName)
	if err != nil {
		return err
	}
	_, err = p.pool.Exec(ctx, `SET TIMEZONE='`+location.String()+`'`)
	return err
}

// Close Закроет все подключения к БД
func Close() {
	mu.Lock()
	defer mu.Unlock()
	for _, c := range pools {
		c.Close()
	}
	pools = make(map[string]*Pool)
}

func pool(connectionName string) (*Pool, error) {
	mu.RLock()
	defer mu.RUnlock()
	p, ok := pools[connectionName]
	if !ok {
		return nil, fmt.Errorf("xpg: Pool of connections `%s` not found", connectionName)
	}
	return p, nil
}
