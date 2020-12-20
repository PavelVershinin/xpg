package xpg_test

import (
	"context"
	"testing"

	"github.com/PavelVershinin/xpg"
	"github.com/stretchr/testify/require"

	"github.com/PavelVershinin/xpg/test"
	"github.com/stretchr/testify/assert"
)

func TestModel_Table(t *testing.T) {
	assert.Equal(t, "test_users", (&test.User{}).Table())
}

func TestModel_Columns(t *testing.T) {
	assert.Equal(t, `
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
	`, (&test.User{}).Columns())
}

func TestModel_Connection(t *testing.T) {
	assert.Equal(t, "test", (&test.User{}).Connection())
}

func TestModel_ScanRow(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	rows, err := xpg.DB("test").Query(context.Background(), `SELECT `+user.Columns()+` FROM `+user.Table()+` LIMIT 1`)
	require.NoError(t, err)
	defer rows.Close()

	if rows.Next() {
		row, err := user.ScanRow(rows)
		require.NoError(t, err)
		user = row.(*test.User)
	}

	assert.Greater(t, user.ID, int64(0))
}

func TestModel_Save(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{
		FirstName:  "Pavel",
		SecondName: "Vershinin",
		LastName:   "Nikolaevich",
		Email:      "xr.pavel@yandex.ru",
		Phone:      "secret!",
		RoleID:     1,
		Balance:    200,
	}

	assert.NoError(t, user.Save())
	assert.Greater(t, user.ID, int64(0))
}

func TestModel_Delete(t *testing.T) {
	defer test.Connect()()

	assert.NoError(t, xpg.New(&test.User{}).WhereNotIn("id", (&xpg.WhereInValues{}).Int64(1, 2)).Delete())
}
