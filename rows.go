package xpg

import (
	"log"

	"github.com/jackc/pgx/v4"
)

// Rows Интерфейс для хранения результата запроса к БД
type Rows struct {
	pgx.Rows
	conn *Connection
}

// Get Получение очередной строки
func (r *Rows) Get() (Tabler, error) {
	return r.conn.tabler.Scan(r.Rows)
}

// Fetch Метод для перебора for row := range res.Fetch() {
func (r *Rows) Fetch() <-chan Tabler {
	var channel = make(chan Tabler)
	go func() {
		for r.Next() {
			if row, err := r.Get(); err != nil {
				log.Println(err.Error())
			} else {
				channel <- row
			}
		}
		r.Close()
		close(channel)
	}()
	return channel
}
