package xpg_test

import (
	"context"
	"testing"

	"github.com/PavelVershinin/xpg"
	"github.com/stretchr/testify/require"

	"github.com/PavelVershinin/xpg/test"
)

func TestModel_Table(t *testing.T) {
	require.Equal(t, "test_users", (&test.User{}).Table())
}

func TestModel_Columns(t *testing.T) {
	require.Equal(t, test.ClearQuery(`
		"test_users"."id",
		"test_users"."first_name",
		"test_users"."second_name",
		"test_users"."last_name",
		"test_users"."email",
		"test_users"."phone",
		"test_users"."role_id",     
		"test_users"."balance",     
		"test_users"."created_at",
		"test_users"."updated_at"
	`), test.ClearQuery((&test.User{}).Columns()))
}

func TestModel_Connection(t *testing.T) {
	require.Equal(t, "test", (&test.User{}).PoolName())
}

func TestModel_ScanRow(t *testing.T) {
	ctx := context.Background()
	pg, err := test.Start(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, test.Stop(pg))
	}()
	require.NoError(t, test.Restore(ctx, 0, 1))

	user := &test.User{}
	rows, err := xpg.DB("test").Query(context.Background(), `SELECT `+user.Columns()+` FROM `+user.Table()+` LIMIT 1`)
	require.NoError(t, err)
	defer rows.Close()

	if rows.Next() {
		row, err := user.ScanRow(rows)
		require.NoError(t, err)
		user = row.(*test.User)
	}

	require.Equal(t, user.ID, int64(1))
}

func TestModel_Save(t *testing.T) {
	ctx := context.Background()
	pg, err := test.Start(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, test.Stop(pg))
	}()
	require.NoError(t, test.Restore(ctx, 0, 0))

	role := &test.Role{
		Name: "Test",
	}

	require.NoError(t, role.Save(ctx))
	require.Equal(t, role.ID, int64(1))

	user := &test.User{
		FirstName:  "Pavel",
		SecondName: "Vershinin",
		LastName:   "Nikolaevich",
		Email:      "xr.pavel@yandex.ru",
		Phone:      "secret!",
		RoleID:     role.ID,
		Balance:    200,
	}

	require.NoError(t, user.Save(ctx))
	require.Equal(t, user.ID, int64(1))
}

func TestModel_Delete(t *testing.T) {
	ctx := context.Background()
	pg, err := test.Start(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, test.Stop(pg))
	}()
	require.NoError(t, test.Restore(ctx, 0, 5))

	require.NoError(t, xpg.New(&test.User{}).WhereNotIn("id", (&xpg.WhereInValues{}).Int64(1, 2)).Delete(ctx))
}
