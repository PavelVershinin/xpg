package test

import (
	"context"
	"strconv"

	"github.com/PavelVershinin/xpg"
	"github.com/PavelVershinin/xpg/migrations"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	connPort = 5435
	connStr  = "postgres://postgres:postgres@localhost:" + strconv.Itoa(connPort) + "/postgres"
)

func Start(ctx context.Context) (*embeddedpostgres.EmbeddedPostgres, error) {
	pg := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().Port(uint32(connPort)))
	if err := pg.Start(); err != nil {
		return nil, err
	}
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}
	if err := xpg.NewConnectionPool(ctx, "test", config, ""); err != nil {
		return nil, err
	}

	return pg, nil
}

func Restore(ctx context.Context, makeRowsNum int) error {
	if err := migrations.Restore(ctx, &User{}); err != nil {
		return err
	}
	if err := migrations.Restore(ctx, &Role{}); err != nil {
		return err
	}

	var roles = make([]*Role, makeRowsNum)
	for i := 1; i <= makeRowsNum; i++ {
		roles[i-1] = &Role{
			Name: "Test " + strconv.Itoa(i),
		}
		if err := roles[i-1].Save(ctx); err != nil {
			return err
		}
	}
	for i := 1; i <= makeRowsNum; i++ {
		user := &User{
			FirstName:  "FirstName " + strconv.Itoa(i),
			SecondName: "SecondName " + strconv.Itoa(i),
			LastName:   "LastName " + strconv.Itoa(i),
			Email:      "my@email.ru",
			Phone:      "secret!",
			Role:       roles[i-1],
			Balance:    100,
		}
		if err := user.Save(ctx); err != nil {
			return err
		}
	}
	return nil
}

func Stop(pg *embeddedpostgres.EmbeddedPostgres) error {
	xpg.Close()
	if err := pg.Stop(); err != nil {
		return err
	}
	return nil
}
