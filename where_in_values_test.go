package xpg

import "testing"

func TestWhereInValues_Int(t *testing.T) {
	var args = []interface{}{
		0, 1, 2,
	}
	var win = (&WhereInValues{}).Int(1, 5, 9, 6, 8)
	sql, args := win.Sql(args)
	if sql != ` IN($4,$5,$6,$7,$8)` {
		t.Log(sql)
		t.Error(`Wrong query`)
	}
	if len(args) != 8 {
		t.Errorf(`Wrong arguments number, expected 1, real %d`, len(args))
	}
}

func TestWhereInValues_Int64(t *testing.T) {
	var args = []interface{}{
		0, 1, 2,
	}
	var win = (&WhereInValues{}).Int64(1, 5, 9, 6, 8)
	sql, args := win.Sql(args)
	if sql != ` IN($4,$5,$6,$7,$8)` {
		t.Log(sql)
		t.Error(`Wrong query`)
	}
	if len(args) != 8 {
		t.Errorf(`Wrong arguments number, expected 1, real %d`, len(args))
	}
}

func TestWhereInValues_String(t *testing.T) {
	var args = []interface{}{
		0, 1, 2,
	}
	var win = (&WhereInValues{}).String("1", "5", "9", "6", "8")
	sql, args := win.Sql(args)
	if sql != ` IN($4,$5,$6,$7,$8)` {
		t.Log(sql)
		t.Error(`Wrong query`)
	}
	if len(args) != 8 {
		t.Errorf(`Wrong arguments number, expected 1, real %d`, len(args))
	}
}

func TestWhereInValues_Interface(t *testing.T) {
	var args = []interface{}{
		0, 1, 2,
	}
	var win = (&WhereInValues{}).Interface([]interface{}{"1", "5", "9", "6", "8"}...)
	sql, args := win.Sql(args)
	if sql != ` IN($4,$5,$6,$7,$8)` {
		t.Log(sql)
		t.Error(`Wrong query`)
	}
	if len(args) != 8 {
		t.Errorf(`Wrong arguments number, expected 1, real %d`, len(args))
	}
}
