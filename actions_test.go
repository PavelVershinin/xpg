package xpg

import (
	"testing"
)

func TestConnection_Insert(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	for i := 0; i < 10; i++ {
		_, err := New(&testModel{}).Insert(map[string]interface{}{
			"column_one":   "insert",
			"column_three": i,
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestConnection_Update(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	for i := 0; i < 10; i++ {
		err := New(&testModel{}).Where("id", "=", i).Update(map[string]interface{}{
			"column_two":   "update",
			"column_three": i + 1,
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestConnection_Write(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	for i := 0; i < 10; i++ {
		_, err := New(&testModel{}).Write(map[string]interface{}{
			"id":           i,
			"column_one":   "text",
			"column_two":   "write",
			"column_three": i,
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestConnection_Delete(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()

	for i := 0; i < 10; i++ {
		err := New(&testModel{}).Where("id", "=", i).Delete()
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestConnection_Select(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	if rows, err := New(&testModel{}).WhereBetween("column_three", 0, 10).OrderBy("column_three", "DESC").Select(); err != nil {
		t.Error(err)
	} else {
		for row := range rows.Fetch() {
			t.Log(row.(*testModel).ColumnThree)
		}
	}
}

func TestConnection_Query(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	if rows, err := New(&testModel{}).Query(`SELECT * FROM "test_model_table"`); err != nil {
		t.Error(err)
	} else {
		for row := range rows.Fetch() {
			t.Log(row.(*testModel).ColumnThree)
		}
	}
}

func TestConnection_First(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	if res, err := New(&testModel{}).WhereBetween("column_three", 0, 10).OrderBy("column_three", "DESC").First(); err != nil {
		t.Error(err)
	} else {
		t.Log(res.(*testModel).ColumnThree)
	}
}

func TestConnection_Exists(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	if ok, err := New(&testModel{}).WhereBetween("column_three", 0, 10).OrderBy("column_three", "DESC").Exists(); err != nil {
		t.Error(err)
	} else {
		t.Log(ok)
	}
}

func TestConnection_Count(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	if cnt, err := New(&testModel{}).WhereBetween("column_three", 0, 10).OrderBy("column_three", "DESC").Count(); err != nil {
		t.Error(err)
	} else {
		t.Log(cnt)
	}
}

func TestConnection_Sum(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	if sum, err := New(&testModel{}).WhereBetween("column_three", 0, 10).OrderBy("column_three", "DESC").Sum("column_three"); err != nil {
		t.Error(err)
	} else {
		t.Log(sum)
	}
}

func TestConnection_WhereIn2(t *testing.T) {
	if err := testConnect(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Error(err)
		}
	}()
	var in = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if sum, err := New(&testModel{}).WhereIn("column_three", (&WhereInValues{}).Int(in...)).OrderBy("column_three", "DESC").Sum("column_three"); err != nil {
		t.Error(err)
	} else {
		t.Log(sum)
	}
}
