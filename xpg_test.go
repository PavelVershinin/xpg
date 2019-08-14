package xpg

import (
	"github.com/jackc/pgx"
	"testing"
	"time"
)

func testConnect() error {
	return NewConnection("xpg_connection", pgx.ConnConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "xpg_test",
		User:     "postgres",
		Password: "123456",
	}, 1, "./migrations/test")
}

func TestNewConnection(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Error(err)
	}
}

func TestSetTimezone(t *testing.T) {
	if loc, err := time.LoadLocation("Europe/Moscow"); err != nil {
		t.Error(err)
	} else if err := SetTimezone("xpg_connection", loc); err != nil {
		t.Error(err)
	}
}

func TestConn(t *testing.T) {
	conn := DB("xpg_connection")
	if conn == nil {
		t.Errorf("Connection %s is nil\n", "test")
	}
	var dest int
	if err := conn.QueryRow(`SELECT 2 * 2`).Scan(&dest); err != nil {
		t.Error(err)
	}
}

func TestClose(t *testing.T) {
	if err := Close(); err != nil {
		t.Error(err)
	}
}
