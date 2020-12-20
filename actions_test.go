package xpg_test

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/PavelVershinin/xpg"
	"github.com/PavelVershinin/xpg/test"
	"github.com/stretchr/testify/assert"
)

func TestConnection_Write(t *testing.T) {
	defer test.Connect()()

	id, err := xpg.New(&test.User{}).Write(map[string]interface{}{
		"first_name":  "Pavel",
		"second_name": "Vershinin",
		"last_name":   "Nikolaevich",
		"email":       "xr.pavel@yandex.ru",
		"phone":       "secret!",
		"role":        1,
		"balance":     100,
	})

	assert.NoError(t, err)
	assert.Greater(t, id, int64(0))
}

func TestConnection_Insert(t *testing.T) {
	defer test.Connect()()

	var role = &test.Role{
		Name: "Test 1",
	}

	assert.NoError(t, role.Save())
	assert.Greater(t, role.ID, int64(0))

	var user = &test.User{
		FirstName:  "Pavel",
		SecondName: "Vershinin",
		LastName:   "Nikolaevich",
		Email:      "xr.pavel@yandex.ru",
		Phone:      "secret!",
		Role:       *role,
		Balance:    100,
	}

	assert.NoError(t, user.Save())
	assert.Greater(t, user.ID, int64(0))
}

func TestConnection_Update(t *testing.T) {
	defer test.Connect()()

	err := xpg.New(&test.User{}).Where("id", "=", 1).Update(map[string]interface{}{
		"role":    2,
		"balance": 120,
	})

	assert.NoError(t, err)
}

func TestConnection_Delete(t *testing.T) {
	defer test.Connect()()

	err := xpg.New(&test.User{}).WhereNotIn("id", (&xpg.WhereInValues{}).Int64(1, 2)).Delete()

	assert.NoError(t, err)
}

func TestConnection_Select(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	rows, err := xpg.New(user).Where("id", "=", 2).Select()
	require.NoError(t, err)

	if rows.Next() {
		row, err := rows.Get()
		require.NoError(t, err)
		user = row.(*test.User)
	}
	rows.Close()

	require.NoError(t, user.Role.DbTake())

	log.Printf("%+v\n", user.Role)

	assert.Equal(t, int64(2), user.ID)
}

func TestConnection_Query(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	rows, err := xpg.New(user).Query("SELECT "+user.Columns()+" FROM "+user.Table()+" WHERE id = $1", 2)
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
	defer test.Connect()()

	row, err := xpg.New(&test.User{}).Where("id", "=", 2).First()
	require.NoError(t, err)
	assert.Equal(t, int64(2), row.(*test.User).ID)
}

func TestConnection_Exists(t *testing.T) {
	defer test.Connect()()

	exists, err := xpg.New(&test.User{}).Exists()
	require.NoError(t, err)

	assert.Equal(t, true, exists)
}

func TestConnection_Count(t *testing.T) {
	defer test.Connect()()

	count, err := xpg.New(&test.User{}).Count()
	require.NoError(t, err)

	assert.Greater(t, count, int64(0))
}

func TestConnection_Sum(t *testing.T) {
	defer test.Connect()()

	sum, err := xpg.New(&test.User{}).Sum("balance")
	require.NoError(t, err)

	assert.Greater(t, sum, float64(0))
}
