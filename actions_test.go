package xpg_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PavelVershinin/xpg"
	"github.com/PavelVershinin/xpg/test"

	"github.com/stretchr/testify/require"
)

func TestConnection_Write(t *testing.T) {
	ctx := context.Background()
	pg, err := test.Start(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, test.Stop(pg))
	}()
	require.NoError(t, test.Restore(ctx, 0))

	id, err := xpg.New(&test.User{}).Write(ctx, map[string]interface{}{
		"first_name":  "Pavel",
		"second_name": "Vershinin",
		"last_name":   "Nikolaevich",
		"email":       "xr.pavel@yandex.ru",
		"phone":       "secret!",
		"role_id":     1,
		"balance":     100,
	})

	require.NoError(t, err)
	require.Equal(t, id, int64(1))
}

func TestConnection_Insert(t *testing.T) {
	ctx := context.Background()
	pg, err := test.Start(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, test.Stop(pg))
	}()
	require.NoError(t, test.Restore(ctx, 2))

	var role = &test.Role{
		Name: "Test 1",
	}

	assert.NoError(t, role.Save(ctx))
	assert.Equal(t, role.ID, int64(3))

	var user = &test.User{
		FirstName:  "Pavel",
		SecondName: "Vershinin",
		LastName:   "Nikolaevich",
		Email:      "xr.pavel@yandex.ru",
		Phone:      "secret!",
		Role:       role,
		Balance:    100,
	}

	assert.NoError(t, user.Save(ctx))
	assert.Equal(t, user.ID, int64(3))
	assert.Equal(t, user.Role.ID, int64(3))
}

func TestConnection_Update(t *testing.T) {
	ctx := context.Background()
	pg, err := test.Start(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, test.Stop(pg))
	}()
	require.NoError(t, test.Restore(ctx, 1))

	err = xpg.New(&test.User{}).Where("id", "=", 1).Update(ctx, map[string]interface{}{
		"role_id": 2,
		"balance": 120,
	})

	assert.NoError(t, err)
}

func TestConnection_Delete(t *testing.T) {
	ctx := context.Background()
	pg, err := test.Start(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, test.Stop(pg))
	}()
	require.NoError(t, test.Restore(ctx, 5))

	assert.NoError(t, xpg.New(&test.User{}).WhereNotIn("id", (&xpg.WhereInValues{}).Int64(1, 2)).Delete(ctx))
}

func TestConnection_Select(t *testing.T) {
	ctx := context.Background()
	pg, err := test.Start(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, test.Stop(pg))
	}()
	require.NoError(t, test.Restore(ctx, 2))

	var user = &test.User{}
	rows, err := xpg.New(user).Where("id", "=", 2).Select(ctx)
	require.NoError(t, err)

	if rows.Next() {
		row, err := rows.Get()
		require.NoError(t, err)
		user = row.(*test.User)
	}
	rows.Close()

	assert.Equal(t, int64(2), user.ID)
}

func TestConnection_Query(t *testing.T) {
	ctx := context.Background()
	pg, err := test.Start(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, test.Stop(pg))
	}()
	require.NoError(t, test.Restore(ctx, 2))

	var user = &test.User{}
	rows, err := xpg.New(user).Query(ctx, "SELECT "+user.Columns()+" FROM "+user.Table()+" WHERE id = $1", 2)
	require.NoError(t, err)
	defer rows.Close()

	if rows.Next() {
		row, err := rows.Get()
		require.NoError(t, err)
		user = row.(*test.User)
	}

	assert.Equal(t, int64(2), user.ID)
}

func TestConnection_First(t *testing.T) {
	ctx := context.Background()
	pg, err := test.Start(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, test.Stop(pg))
	}()
	require.NoError(t, test.Restore(ctx, 2))

	row, err := xpg.New(&test.User{}).Where("id", "=", 2).First(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(2), row.(*test.User).ID)
}

func TestConnection_Exists(t *testing.T) {
	ctx := context.Background()
	pg, err := test.Start(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, test.Stop(pg))
	}()
	require.NoError(t, test.Restore(ctx, 2))

	exists, err := xpg.New(&test.User{}).Where("id", "=", 2).Exists(ctx)
	require.NoError(t, err)

	assert.Equal(t, true, exists)
}

func TestConnection_Count(t *testing.T) {
	ctx := context.Background()
	pg, err := test.Start(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, test.Stop(pg))
	}()
	require.NoError(t, test.Restore(ctx, 4))

	count, err := xpg.New(&test.User{}).Count(ctx)
	require.NoError(t, err)

	assert.Equal(t, count, int64(4))
}

func TestConnection_Sum(t *testing.T) {
	ctx := context.Background()
	pg, err := test.Start(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, test.Stop(pg))
	}()
	require.NoError(t, test.Restore(ctx, 4))

	sum, err := xpg.New(&test.User{}).Sum(ctx, "balance")
	require.NoError(t, err)

	assert.Equal(t, sum, float64(400))
}
