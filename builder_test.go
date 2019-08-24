package xpg

import (
	"strings"
	"testing"
)

func TestConnection_Limit(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	var query = New(&testModel{})
	query.Limit(10)
	var sql, args = query.BuildSelect()
	if sql != `SELECT * FROM "test_model_table" LIMIT 10` {
		t.Log(sql)
		t.Error(`Wrong query`)
	}
	if len(args) != 0 {
		t.Errorf(`Wrong arguments number, expected 0, real %d`, len(args))
	}
}

func TestConnection_Offset(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	var query = New(&testModel{})
	query.Offset(10)
	var sql, args = query.BuildSelect()
	if sql != `SELECT * FROM "test_model_table" OFFSET 10` {
		t.Log(sql)
		t.Error(`Wrong query`)
	}
	if len(args) != 0 {
		t.Errorf(`Wrong arguments number, expected 0, real %d`, len(args))
	}
}

func TestConnection_Where(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	var query = New(&testModel{})
	query.Where("column_one", "=", "test")
	var sql, args = query.BuildSelect()
	if sql != `SELECT * FROM "test_model_table" WHERE ("column_one"=$1)` {
		t.Log(sql)
		t.Error(`Wrong query`)
	}
	if len(args) != 1 {
		t.Errorf(`Wrong arguments number, expected 1, real %d`, len(args))
	}
}

func TestConnection_OrWhere(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	var query = New(&testModel{})
	query.Where("column_one", "=", "test")
	query.OrWhere("column_two", "=", "test")
	var sql, args = query.BuildSelect()
	if sql != `SELECT * FROM "test_model_table" WHERE ("column_one"=$1 OR "column_two"=$2)` {
		t.Log(sql)
		t.Error(`Wrong query`)
	}
	if len(args) != 2 {
		t.Errorf(`Wrong arguments number, expected 2, real %d`, len(args))
	}
}

func TestConnection_OrGroupWhere(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	var query = New(&testModel{})
	query.GroupWhere(func(c *Connection) {
		c.Where("column_one", "=", "test")
		c.OrWhere("column_two", "=", "test")
	})
	query.OrGroupWhere(func(c *Connection) {
		c.Where("column_one", "=", "test")
		c.OrWhere("column_two", "=", "test")
	})
	var sql, args = query.BuildSelect()
	if sql != `SELECT * FROM "test_model_table" WHERE ("column_one"=$1 OR "column_two"=$2) OR ("column_one"=$3 OR "column_two"=$4)` {
		t.Log(sql)
		t.Error(`Wrong query`)
	}
	if len(args) != 4 {
		t.Errorf(`Wrong arguments number, expected 4, real %d`, len(args))
	}
}

func TestConnection_GroupWhere(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	var query = New(&testModel{})
	query.GroupWhere(func(c *Connection) {
		c.Where("column_one", "=", "test")
		c.OrWhere("column_two", "=", "test")
	})
	query.GroupWhere(func(c *Connection) {
		c.Where("column_one", "=", "test")
		c.OrWhere("column_two", "=", "test")
	})
	var sql, args = query.BuildSelect()
	if sql != `SELECT * FROM "test_model_table" WHERE ("column_one"=$1 OR "column_two"=$2) AND ("column_one"=$3 OR "column_two"=$4)` {
		t.Log(sql)
		t.Error(`Wrong query`)
	}
	if len(args) != 4 {
		t.Errorf(`Wrong arguments number, expected 4, real %d`, len(args))
	}
}

func TestConnection_WhereRaw(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	var query = New(&testModel{})
	query.WhereRaw(`"column_one"=$1 AND "column_two"=$2`, "test", "test2")
	query.WhereRaw(`"column_one"=$1 AND "column_two"=$2`, "test", "test2")
	var sql, args = query.BuildSelect()
	if sql != `SELECT * FROM "test_model_table" WHERE (("column_one"=$1 AND "column_two"=$2) AND ("column_one"=$3 AND "column_two"=$4))` {
		t.Log(sql)
		t.Error(`Wrong query`)
	}
	if len(args) != 4 {
		t.Errorf(`Wrong arguments number, expected 4, real %d`, len(args))
	}
}

func TestConnection_OrWhereRaw(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	var query = New(&testModel{})
	query.WhereRaw(`"column_one"=$1 AND "column_two"=$2`, "test", "test2")
	query.OrWhereRaw(`"column_one"!=$1 AND "column_two"!=$2`, "test", "test2")
	var sql, args = query.BuildSelect()
	if sql != `SELECT * FROM "test_model_table" WHERE (("column_one"=$1 AND "column_two"=$2) OR ("column_one"!=$3 AND "column_two"!=$4))` {
		t.Log(sql)
		t.Error(`Wrong query`)
	}
	if len(args) != 4 {
		t.Errorf(`Wrong arguments number, expected 4, real %d`, len(args))
	}
}

func TestConnection_WhereIn(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	var in = []int{1, 2, 3, 4, 5, 6, 7}
	var query = New(&testModel{})
	query.Where("id", "=", 10)
	query.WhereIn("id", (&WhereInValues{}).Int(in...))
	query.Where("id", "=", 11)
	var sql, args = query.BuildSelect()
	if sql != `SELECT * FROM "test_model_table" WHERE ("id"=$1 AND "id" IN($2,$3,$4,$5,$6,$7,$8) AND "id"=$9)` {
		t.Log(sql)
		t.Error(`Wrong query`)
	}
	if len(args) != 9 {
		t.Errorf(`Wrong arguments number, expected 1, real %d`, len(args))
	}
}

func TestConnection_WhereNotIn(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	var in = []int{1, 2, 3, 4, 5, 6, 7}
	var query = New(&testModel{})
	query.Where("id", "=", 10)
	query.WhereNotIn("id", (&WhereInValues{}).Int(in...))
	query.Where("id", "=", 11)
	var sql, args = query.BuildSelect()
	if sql != `SELECT * FROM "test_model_table" WHERE ("id"=$1 AND "id" NOT IN($2,$3,$4,$5,$6,$7,$8) AND "id"=$9)` {
		t.Log(sql)
		t.Error(`Wrong query`)
	}
	if len(args) != 9 {
		t.Errorf(`Wrong arguments number, expected 1, real %d`, len(args))
	}
}

func TestConnection_WhereBetween(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	var query = New(&testModel{})
	query.WhereBetween("id", 1, 5)
	var sql, args = query.BuildSelect()
	if sql != `SELECT * FROM "test_model_table" WHERE ("id">=$1 AND "id"<=$2)` {
		t.Log(sql)
		t.Error(`Wrong query`)
	}
	if len(args) != 2 {
		t.Errorf(`Wrong arguments number, expected 2, real %d`, len(args))
	}
}

func TestConnection_OrWhereBetween(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	var query = New(&testModel{})
	query.WhereBetween("id", 1, 5)
	query.OrWhereBetween("id", 10, 15)
	var sql, args = query.BuildSelect()
	if sql != `SELECT * FROM "test_model_table" WHERE ("id">=$1 AND "id"<=$2) OR ("id">=$3 AND "id"<=$4)` {
		t.Log(sql)
		t.Error(`Wrong query`)
	}
	if len(args) != 4 {
		t.Errorf(`Wrong arguments number, expected 4, real %d`, len(args))
	}
}

func TestConnection_Distinct(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	var query = New(&testModel{})
	query.Distinct()
	var sql, args = query.BuildSelect()
	if sql != `SELECT DISTINCT * FROM "test_model_table"` {
		t.Log(sql)
		t.Error(`Wrong query`)
	}
	if len(args) != 0 {
		t.Errorf(`Wrong arguments number, expected 0, real %d`, len(args))
	}
}

func TestConnection_Distinct2(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	var query = New(&testModel{})
	query.Distinct("id", "column_one")
	var sql, args = query.BuildSelect()
	if sql != `SELECT DISTINCT ON("id","column_one") * FROM "test_model_table"` {
		t.Log(sql)
		t.Error(`Wrong query`)
	}
	if len(args) != 0 {
		t.Errorf(`Wrong arguments number, expected 0, real %d`, len(args))
	}
}

func TestConnection_GroupBy(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	var query = New(&testModel{})
	query.GroupBy("id", "column_one")
	var sql, args = query.BuildSelect()
	if sql != `SELECT * FROM "test_model_table" GROUP BY "id", "column_one"` {
		t.Log(sql)
		t.Error(`Wrong query`)
	}
	if len(args) != 0 {
		t.Errorf(`Wrong arguments number, expected 0, real %d`, len(args))
	}
}

func TestConnection_OrderBy(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	var query = New(&testModel{})
	query.OrderBy("id", "ASC")
	query.OrderBy("created_at", "DESC")
	var sql, args = query.BuildSelect()
	if sql != `SELECT * FROM "test_model_table" ORDER BY "id" ASC, "created_at" DESC` {
		t.Log(sql)
		t.Error(`Wrong query`)
	}
	if len(args) != 0 {
		t.Errorf(`Wrong arguments number, expected 0, real %d`, len(args))
	}
}

func TestConnection_OrderByRaw(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	var query = New(&testModel{})
	query.OrderByRaw("id, created_at ASC")
	var sql, args = query.BuildSelect()
	if sql != `SELECT * FROM "test_model_table" ORDER BY id, created_at ASC` {
		t.Log(sql)
		t.Error(`Wrong query`)
	}
	if len(args) != 0 {
		t.Errorf(`Wrong arguments number, expected 0, real %d`, len(args))
	}
}

func TestConnection_OrderByRand(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	var query = New(&testModel{})
	query.OrderByRand()
	var sql, args = query.BuildSelect()
	if sql != `SELECT * FROM "test_model_table" ORDER BY RANDOM()` {
		t.Log(sql)
		t.Error(`Wrong query`)
	}
	if len(args) != 0 {
		t.Errorf(`Wrong arguments number, expected 0, real %d`, len(args))
	}
}

func TestConnection_Union(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()

	var query1 = New(&testModel{}).WhereBetween("id", 1, 4).OrderBy("id", "DESC").Limit(20).Offset(1)
	var query2 = New(&testModel{}).WhereBetween("id", 6, 8).OrderBy("id", "ASC").Limit(20).Offset(1)
	var query3 = New(&testModel{}).WhereBetween("id", 10, 14).OrderBy("id", "ASC").Limit(20).Offset(1)
	var query4 = New(&testModel{}).WhereBetween("id", 20, 22).OrderBy("id", "DESC").Limit(20).Offset(1)
	var query = New(&testModel{})
	query.Union(true, query1, query2)
	query.Union(false, query3, query4)
	query.Where("id", ">", 5).OrWhere("id", "<", 20).Limit(20).Offset(1).OrderBy("created_at", "DESC")
	var sql, args = query.BuildSelect()
	var expected = strings.TrimSpace(`
SELECT * FROM (
	SELECT * FROM "test_model_table" WHERE ("id">=$1 AND "id"<=$2) ORDER BY "id" DESC OFFSET 1 LIMIT 20
	UNION ALL SELECT * FROM "test_model_table" WHERE ("id">=$3 AND "id"<=$4) ORDER BY "id" ASC OFFSET 1 LIMIT 20
	UNION SELECT * FROM "test_model_table" WHERE ("id">=$5 AND "id"<=$6) ORDER BY "id" ASC OFFSET 1 LIMIT 20
	UNION SELECT * FROM "test_model_table" WHERE ("id">=$7 AND "id"<=$8) ORDER BY "id" DESC OFFSET 1 LIMIT 20
) AS "xpg_union_test_model_table" WHERE ("id">$9 OR "id"<$10) ORDER BY "created_at" DESC OFFSET 1 LIMIT 20
`)
	if sql != expected {
		t.Log(sql)
		t.Error(`Wrong query`)
	}
	if len(args) != 10 {
		t.Errorf(`Wrong arguments number, expected 10, real %d`, len(args))
	}
}

func TestConnection_BuildCount(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	var query = New(&testModel{})
	query.Where("id", ">", 5)
	var sql, args = query.BuildCount()
	if sql != `SELECT COUNT(*) AS "count" FROM "test_model_table" WHERE ("id">$1)` {
		t.Log(sql)
		t.Error(`Wrong query`)
	}
	if len(args) != 1 {
		t.Errorf(`Wrong arguments number, expected 1, real %d`, len(args))
	}
}

func TestConnection_BuildCount2(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	var query = New(&testModel{})
	query.Where("id", ">", 5)
	query.Distinct()
	var sql, args = query.BuildCount()
	if sql != `SELECT COUNT(DISTINCT *) AS "count" FROM "test_model_table" WHERE ("id">$1)` {
		t.Log(sql)
		t.Error(`Wrong query`)
	}
	if len(args) != 1 {
		t.Errorf(`Wrong arguments number, expected 1, real %d`, len(args))
	}
}

func TestConnection_BuildSum(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	var query = New(&testModel{})
	query.Where("id", ">", 5)
	var sql, args = query.BuildSum("column_three")
	if sql != `SELECT SUM("column_three") FROM "test_model_table" WHERE ("id">$1)` {
		t.Log(sql)
		t.Error(`Wrong query`)
	}
	if len(args) != 1 {
		t.Errorf(`Wrong arguments number, expected 1, real %d`, len(args))
	}
}
