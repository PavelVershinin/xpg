package xpg_test

import (
	"testing"

	"github.com/PavelVershinin/xpg"

	"github.com/stretchr/testify/require"
)

func TestWhereInValues_Int(t *testing.T) {
	wv := (&xpg.WhereInValues{}).Int(1, 5, 9, 6, 8)
	sql, args := wv.Sql([]interface{}{0, 1, 2})
	require.Equal(t, " IN($4,$5,$6,$7,$8)", sql)
	require.Equal(t, []interface{}{0, 1, 2, 1, 5, 9, 6, 8}, args)
}

func TestWhereInValues_Int64(t *testing.T) {
	wv := (&xpg.WhereInValues{}).Int64(1, 5, 9, 6, 8)
	sql, args := wv.Sql([]interface{}{0, 1, 2})
	require.Equal(t, " IN($4,$5,$6,$7,$8)", sql)
	require.Equal(t, []interface{}{0, 1, 2, int64(1), int64(5), int64(9), int64(6), int64(8)}, args)
}

func TestWhereInValues_String(t *testing.T) {
	wv := (&xpg.WhereInValues{}).String("1", "5", "9", "6", "8")
	sql, args := wv.Sql([]interface{}{0, 1, 2})
	require.Equal(t, " IN($4,$5,$6,$7,$8)", sql)
	require.Equal(t, []interface{}{0, 1, 2, "1", "5", "9", "6", "8"}, args)
}

func TestWhereInValues_Interface(t *testing.T) {
	wv := (&xpg.WhereInValues{}).Interface("1", "5", "9", "6", "8")
	sql, args := wv.Sql([]interface{}{0, 1, 2})
	require.Equal(t, " IN($4,$5,$6,$7,$8)", sql)
	require.Equal(t, []interface{}{0, 1, 2, "1", "5", "9", "6", "8"}, args)
}
