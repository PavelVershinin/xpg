package xpg_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PavelVershinin/xpg"
	"github.com/PavelVershinin/xpg/test"
)

func TestConnection_Join(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var role = &test.Role{}
	var query = xpg.New(user).
		Where(`"`+user.Table()+`"."id"`, "=", 1).
		Join(role.Table(), role.Table(), `"`+role.Table()+`"."id" = "`+user.Table()+`"."role"`)
	var expectSql = `SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users" INNER JOIN test_roles AS test_roles ON("test_roles"."id" = "test_users"."role") WHERE ("test_users"."id"=$1)`
	var expectArgs = []interface{}{1}
	var sql, args = query.BuildSelect()

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}

func TestConnection_LeftJoin(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var role = &test.Role{}
	var query = xpg.New(user).
		Where(`"`+user.Table()+`"."id"`, "=", 1).
		LeftJoin(role.Table(), role.Table(), `"`+role.Table()+`"."id" = "`+user.Table()+`"."role"`)
	var expectSql = `SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users" LEFT JOIN test_roles AS test_roles ON("test_roles"."id" = "test_users"."role") WHERE ("test_users"."id"=$1)`
	var expectArgs = []interface{}{1}
	var sql, args = query.BuildSelect()

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}

func TestConnection_RightJoin(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var role = &test.Role{}
	var query = xpg.New(user).
		Where(`"`+user.Table()+`"."id"`, "=", 1).
		RightJoin(role.Table(), role.Table(), `"`+role.Table()+`"."id" = "`+user.Table()+`"."role"`)
	var expectSql = `SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users" RIGHT JOIN test_roles AS test_roles ON("test_roles"."id" = "test_users"."role") WHERE ("test_users"."id"=$1)`
	var expectArgs = []interface{}{1}
	var sql, args = query.BuildSelect()

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}

func TestConnection_FullJoin(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var role = &test.Role{}
	var query = xpg.New(user).
		Where(`"`+user.Table()+`"."id"`, "=", 1).
		FullJoin(role.Table(), role.Table(), `"`+role.Table()+`"."id" = "`+user.Table()+`"."role"`)
	var expectSql = `SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users" FULL JOIN test_roles AS test_roles ON("test_roles"."id" = "test_users"."role") WHERE ("test_users"."id"=$1)`
	var expectArgs = []interface{}{1}
	var sql, args = query.BuildSelect()

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}

func TestConnection_Union(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var subQuery1 = xpg.New(user).Where("id", "=", 1)
	var subQuery2 = xpg.New(user).Where("id", "=", 2)
	var subQuery3 = xpg.New(user).Where("id", "=", 3)
	var query = xpg.New(user).Union(true, subQuery1, subQuery2, subQuery3)

	var expectSql = `SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM ( SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users" WHERE ("id"=$1) UNION ALL SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users" WHERE ("id"=$2) UNION ALL SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users" WHERE ("id"=$3)) AS "test_users"`
	var expectArgs = []interface{}{1, 2, 3}
	var sql, args = query.BuildSelect()

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}

func TestConnection_Limit(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var query = xpg.New(user).
		Where("id", "=", 1).
		Limit(1)
	var expectSql = `SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users" WHERE ("id"=$1) LIMIT 1`
	var expectArgs = []interface{}{1}
	var sql, args = query.BuildSelect()

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}

func TestConnection_Offset(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var query = xpg.New(user).
		Where("id", "=", 1).
		Offset(1)
	var expectSql = `SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users" WHERE ("id"=$1) OFFSET 1`
	var expectArgs = []interface{}{1}
	var sql, args = query.BuildSelect()

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}

func TestConnection_Where(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var query = xpg.New(user).
		Where("id", "=", 1).
		Where("role", "=", 2)
	var expectSql = `SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users" WHERE ("id"=$1 AND "role"=$2)`
	var expectArgs = []interface{}{1, 2}
	var sql, args = query.BuildSelect()

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}

func TestConnection_OrWhere(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var query = xpg.New(user).
		Where("id", "=", 1).
		OrWhere("role", "=", 2)
	var expectSql = `SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users" WHERE ("id"=$1 OR "role"=$2)`
	var expectArgs = []interface{}{1, 2}
	var sql, args = query.BuildSelect()

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}

func TestConnection_WhereBetween(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var query = xpg.New(user).
		WhereBetween("id", 1, 10)
	var expectSql = `SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users" WHERE ("id">=$1 AND "id"<=$2)`
	var expectArgs = []interface{}{1, 10}
	var sql, args = query.BuildSelect()

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}

func TestConnection_OrWhereBetween(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var query = xpg.New(user).
		WhereBetween("id", 1, 10).
		OrWhereBetween("id", 20, 30)
	var expectSql = `SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users" WHERE ("id">=$1 AND "id"<=$2) OR ("id">=$3 AND "id"<=$4)`
	var expectArgs = []interface{}{1, 10, 20, 30}
	var sql, args = query.BuildSelect()

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}

func TestConnection_GroupWhere(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var query = xpg.New(user).
		Where("role", "=", 1).
		GroupWhere(func(c *xpg.Connection) {
			c.Where("id", "=", 1)
			c.OrWhere("id", "=", 2)
			c.OrWhere("id", "=", 3)
		})
	var expectSql = `SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users" WHERE ("role"=$1) AND ("id"=$2 OR "id"=$3 OR "id"=$4)`
	var expectArgs = []interface{}{1, 1, 2, 3}
	var sql, args = query.BuildSelect()

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}

func TestConnection_OrGroupWhere(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var query = xpg.New(user).
		Where("role", "=", 1).
		OrGroupWhere(func(c *xpg.Connection) {
			c.Where("id", "=", 1)
			c.OrWhere("id", "=", 2)
			c.OrWhere("id", "=", 3)
		})
	var expectSql = `SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users" WHERE ("role"=$1) OR ("id"=$2 OR "id"=$3 OR "id"=$4)`
	var expectArgs = []interface{}{1, 1, 2, 3}
	var sql, args = query.BuildSelect()

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}

func TestConnection_WhereRaw(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var query = xpg.New(user).
		WhereRaw(`("id" = $1 AND "role" = $2) OR "role" = $3`, 1, 2, 3)
	var expectSql = `SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users" WHERE ((("id" = $1 AND "role" = $2) OR "role" = $3))`
	var expectArgs = []interface{}{1, 2, 3}
	var sql, args = query.BuildSelect()

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}

func TestConnection_OrWhereRaw(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var query = xpg.New(user).
		WhereRaw(`("id" = $1 AND "role" = $2) OR "role" = $3`, 1, 2, 3).
		OrWhereRaw(`("id" = $1 AND "role" = $2) OR "role" = $3`, 4, 5, 6)
	var expectSql = `SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users" WHERE ((("id" = $1 AND "role" = $2) OR "role" = $3) OR (("id" = $4 AND "role" = $5) OR "role" = $6))`
	var expectArgs = []interface{}{1, 2, 3, 4, 5, 6}
	var sql, args = query.BuildSelect()

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}

func TestConnection_WhereIn(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var query = xpg.New(user).
		WhereIn("id", (&xpg.WhereInValues{}).Int(1, 2, 3, 4))
	var expectSql = `SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users" WHERE ("id" IN($1,$2,$3,$4))`
	var expectArgs = []interface{}{1, 2, 3, 4}
	var sql, args = query.BuildSelect()

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}

func TestConnection_OrWhereIn(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var query = xpg.New(user).
		WhereIn("id", (&xpg.WhereInValues{}).Int(1, 2, 3, 4)).
		OrWhereIn("id", (&xpg.WhereInValues{}).Int(5, 6, 7, 8))
	var expectSql = `SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users" WHERE ("id" IN($1,$2,$3,$4) OR "id" IN($5,$6,$7,$8))`
	var expectArgs = []interface{}{1, 2, 3, 4, 5, 6, 7, 8}
	var sql, args = query.BuildSelect()

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}

func TestConnection_WhereNotIn(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var query = xpg.New(user).
		WhereNotIn("id", (&xpg.WhereInValues{}).Int(1, 2, 3, 4))
	var expectSql = `SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users" WHERE ("id" NOT IN($1,$2,$3,$4))`
	var expectArgs = []interface{}{1, 2, 3, 4}
	var sql, args = query.BuildSelect()

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}

func TestConnection_OrWhereNotIn(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var query = xpg.New(user).
		WhereNotIn("id", (&xpg.WhereInValues{}).Int(1, 2, 3, 4)).
		OrWhereNotIn("id", (&xpg.WhereInValues{}).Int(5, 6, 7, 8))
	var expectSql = `SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users" WHERE ("id" NOT IN($1,$2,$3,$4) OR "id" NOT IN($5,$6,$7,$8))`
	var expectArgs = []interface{}{1, 2, 3, 4, 5, 6, 7, 8}
	var sql, args = query.BuildSelect()

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}

func TestConnection_GroupBy(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var query = xpg.New(user).
		GroupBy("email", "phone")
	var expectSql = `SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users" GROUP BY email, phone`
	var expectArgs []interface{}
	var sql, args = query.BuildSelect()

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}

func TestConnection_Distinct(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var query = xpg.New(user).
		Distinct("email", "phone")
	var expectSql = `SELECT DISTINCT ON(email,phone) "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users"`
	var expectArgs []interface{}
	var sql, args = query.BuildSelect()

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}

func TestConnection_OrderBy(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var query = xpg.New(user).
		OrderBy("email", "ASC")
	var expectSql = `SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users" ORDER BY email ASC`
	var expectArgs []interface{}
	var sql, args = query.BuildSelect()

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}

func TestConnection_OrderByRaw(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var query = xpg.New(user).
		OrderByRaw(`email ASC`)
	var expectSql = `SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users" ORDER BY email ASC`
	var expectArgs []interface{}
	var sql, args = query.BuildSelect()

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}

func TestConnection_OrderByRand(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var query = xpg.New(user).
		OrderByRand()
	var expectSql = `SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users" ORDER BY RANDOM()`
	var expectArgs []interface{}
	var sql, args = query.BuildSelect()

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}

func TestConnection_BuildSelect(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var query = xpg.New(user)
	var expectSql = `SELECT "test_users"."id","test_users"."first_name","test_users"."second_name","test_users"."last_name","test_users"."email","test_users"."phone","test_users"."role","test_users"."balance","test_users"."created_at","test_users"."updated_at" FROM "test_users"`
	var expectArgs []interface{}
	var sql, args = query.BuildSelect()

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}

func TestConnection_BuildSum(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var query = xpg.New(user)
	var expectSql = `SELECT COALESCE(SUM("id"), 0) FROM "test_users"`
	var expectArgs []interface{}
	var sql, args = query.BuildSum("id")

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}

func TestConnection_BuildCount(t *testing.T) {
	defer test.Connect()()

	var user = &test.User{}
	var query = xpg.New(user)
	var expectSql = `SELECT COUNT(*) AS "count" FROM "test_users"`
	var expectArgs []interface{}
	var sql, args = query.BuildCount()

	assert.Equal(t, expectSql, test.ClearQuery(sql))
	assert.Equal(t, expectArgs, args)
}
