package xpg

import (
	"bytes"
	"strconv"
)

// WhereInValues Адаптор для передачи слайсов, в запрос WHERE IN
type WhereInValues struct {
	values []interface{}
}

// Int Загрузка ...int
func (w *WhereInValues) Int(in ...int) *WhereInValues {
	w.values = make([]interface{}, len(in))
	for i, v := range in {
		w.values[i] = v
	}
	return w
}

// Int64 Загрузка ...int64
func (w *WhereInValues) Int64(in ...int64) *WhereInValues {
	w.values = make([]interface{}, len(in))
	for i, v := range in {
		w.values[i] = v
	}
	return w
}

// String Загрузка ...string
func (w *WhereInValues) String(in ...string) *WhereInValues {
	w.values = make([]interface{}, len(in))
	for i, v := range in {
		w.values[i] = v
	}
	return w
}

// Interface Загрузка ...interface{}
func (w *WhereInValues) Interface(in ...interface{}) *WhereInValues {
	w.values = make([]interface{}, len(in))
	copy(w.values, in)
	return w
}

// Sql Вернёт подготовленную строку запроса и дополненный слайс аргументов
func (w *WhereInValues) Sql(args []interface{}) (string, []interface{}) {
	var buff bytes.Buffer
	var start = len(args) + 1
	args = append(args, w.values...)
	var end = len(args)
	buff.WriteString(" IN(")
	for i := start; i <= end; i++ {
		buff.WriteString("$")
		buff.WriteString(strconv.Itoa(i))
		if i < end {
			buff.WriteString(",")
		}
	}
	buff.WriteString(")")
	return buff.String(), args
}
