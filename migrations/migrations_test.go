package migrations

import (
	"github.com/PavelVershinin/xpg"
	"github.com/jackc/pgx"
	"testing"
)

type testRestore struct {
	xpg.Model
	ColumnOne   string `xpg:"column_one CHAR(50) NOT NULL DEFAULT ''"`
	ColumnTwo   string `xpg:"column_two CHAR(50) NOT NULL DEFAULT ''"`
	ColumnThree int64  `xpg:"column_three BIGINT NOT NULL DEFAULT 0"`
}

func (tr testRestore) Table() string {
	return "test_restore"
}

func (tr testRestore) Connection() string {
	return "xpg_connection"
}

func testConnect() error {
	return xpg.NewConnection("xpg_connection", pgx.ConnConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "xpg_test",
		User:     "postgres",
		Password: "123456",
	}, 1, "./test")
}

func TestRestore(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := xpg.Close(); err != nil {
			t.Error(err)
		}
	}()
	// Грохнем таблицу если есть
	if _, err := xpg.DB("xpg_connection").Exec(`DROP TABLE "test_restore"`); err != nil {
		t.Log(err)
	}
	// Попробуем восстановить таблицу
	if err := Restore(&testRestore{}); err != nil {
		t.Fatal(err)
	}
	// Грохнем в таблице несколько колонок
	if _, err := xpg.DB("xpg_connection").Exec(`ALTER TABLE "test_restore" DROP COLUMN "column_one", DROP COLUMN "column_three" CASCADE`); err != nil {
		t.Log(err)
	}
	// Попробуем восстановить колонки
	if err := Restore(&testRestore{}); err != nil {
		t.Error(err)
	}
}

func TestUp(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := xpg.Close(); err != nil {
			t.Error(err)
		}
	}()
	if err := Up("xpg_connection", -1); err != nil {
		t.Fatal(err)
	}
}

//func TestDown(t *testing.T) {
//	if err := testConnect(); err != nil {
//		t.Fatal(err)
//	}
//	defer func() {
//		if err := xpg.Close(); err != nil {
//			t.Error(err)
//		}
//	}()
//	if err := Down("xpg_connection", 0); err != nil {
//		t.Fatal(err)
//	}
//}
