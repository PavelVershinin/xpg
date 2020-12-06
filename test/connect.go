package test

import (
	"context"
	"log"
	"os"

	"github.com/PavelVershinin/xpg/migrations"

	"github.com/PavelVershinin/xpg"
	"github.com/jackc/pgx/v4"
)

func Connect() func() {
	connString := os.Getenv("XPG_CONN_STRING")
	if connString == "" {
		log.Fatal("missing environment variable XPG_CONN_STRING")
	}

	config, err := pgx.ParseConfig(connString)
	if err != nil {
		log.Fatal(err)
	}
	config.LogLevel = pgx.LogLevelDebug

	err = xpg.NewConnection(context.Background(), "test", config, "")
	if err != nil {
		log.Fatal(err)
	}

	if err := migrations.Restore(&User{}); err != nil {
		log.Fatal(err)
	}

	if err := migrations.Restore(&Role{}); err != nil {
		log.Fatal(err)
	}

	return func() {
		if err := xpg.Close(); err != nil {
			log.Println(err)
		}
	}
}
