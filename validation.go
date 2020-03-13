package xpg

import (
	"github.com/PavelVershinin/xpg/xpgtypes"

	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Validation Валидация данных перед записью
func (c *Connection) Validation(data map[string]interface{}) (validData map[string]interface{}, err error) {
	validData = make(map[string]interface{})
	var columns []Column

	if columns, err = c.Columns(); err != nil {
		return nil, err
	}

	for _, column := range columns {
		if value, ok := data[column.Name]; ok {
			var typeName, maxLength, isArray = parseColumnType(column.Type)

			var validValuesInt []int64
			var validValuesFloat []float64
			var validValuesString []string
			var validValuesBool []bool
			var validValuesTime []xpgtypes.NullTime

			for _, value := range toArray(value) {
				switch typeName {
				case "bigint":
					num := toInteger(value)
					if num > 9223372036854775807 {
						num = 9223372036854775807
					}
					if num < -9223372036854775808 {
						num = -9223372036854775808
					}
					validValuesInt = append(validValuesInt, num)
				case "integer":
					num := toInteger(value)
					if num > 2147483647 {
						num = 2147483647
					}
					if num < -2147483648 {
						num = -2147483648
					}
					validValuesInt = append(validValuesInt, num)
				case "smallint":
					num := toInteger(value)
					if num > 32767 {
						num = 32767
					}
					if num < -32768 {
						num = -32768
					}
					validValuesInt = append(validValuesInt, num)
				case "boolean":
					validValuesBool = append(validValuesBool, toBoolean(value))
				case "character",
					"text":
					validValuesString = append(validValuesString, toString(value, maxLength))
				case "double",
					"numeric",
					"money":
					validValuesFloat = append(validValuesFloat, toFloat(value))
				case "cidr":
					str := toString(value, 0)
					str = strings.Replace(str, ":", "/", 1)
					re := regexp.MustCompile(`^([0-9]|[0-9][0-9]|[01][0-9][0-9]|2[0-4][0-9]|25[0-5])(\.([0-9]|[0-9][0-9]|[01][0-9][0-9]|2[0-4][0-9]|25[0-5])){3}/([0-9]|[0-9][0-9]|[01][0-9][0-9]|2[0-4][0-9]|25[0-5])$`)
					if re.MatchString(str) {
						validValuesString = append(validValuesString, str)
					}
				case "timestamp",
					"date",
					"time":
					validValuesTime = append(validValuesTime, toTime(value))
				default:
					var enums, _ = c.Enums()
					for name, enumValues := range enums {
						if strings.ToLower(name) == typeName {
							str := toString(value, 0)
							if inArray(enumValues, str) {
								validValuesString = append(validValuesString, str)
							}
							break
						}
					}
				}
			}

			switch typeName {
			case "bigint",
				"integer",
				"smallint":
				if !isArray && len(validValuesInt) > 0 {
					validData[column.Name] = validValuesInt[0]
				} else if isArray {
					validData[column.Name] = validValuesInt
				}
			case "boolean":
				if !isArray && len(validValuesBool) > 0 {
					validData[column.Name] = validValuesBool[0]
				} else if isArray {
					validData[column.Name] = validValuesBool
				}
			case "double",
				"numeric",
				"money":
				if !isArray && len(validValuesFloat) > 0 {
					validData[column.Name] = validValuesFloat[0]
				} else if isArray {
					validData[column.Name] = validValuesFloat
				}
			case "timestamp",
				"date",
				"time":
				if !isArray && len(validValuesTime) > 0 {
					validData[column.Name] = validValuesTime[0]
				} else if isArray {
					validData[column.Name] = validValuesTime
				}
			default:
				if !isArray && len(validValuesString) > 0 {
					validData[column.Name] = validValuesString[0]
				} else if isArray {
					validData[column.Name] = validValuesString
				}
			}
		}
	}

	return
}

func parseColumnType(typeName string) (varTypeName string, maxLength int, isArray bool) {
	typeName = strings.ToLower(typeName)

	varTypeName = strings.Split(strings.Split(strings.Fields(typeName)[0], "(")[0], "[")[0]
	isArray = strings.Contains(typeName, "[]")

	if arr := strings.SplitN(typeName, "(", 2); len(arr) == 2 {
		maxLength, _ = strconv.Atoi(strings.TrimSpace(strings.SplitN(arr[1], ")", 2)[0]))
	}

	return
}

func toInteger(value interface{}) int64 {
	var valueOf = reflect.ValueOf(value)
	switch valueOf.Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return int64(valueOf.Int())
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr:
		return int64(valueOf.Uint())
	case reflect.Float32,
		reflect.Float64:
		return int64(valueOf.Float())
	case reflect.String:
		num, _ := strconv.ParseInt(valueOf.String(), 10, 64)
		return num
	case reflect.Bool:
		if valueOf.Bool() {
			return 1
		} else {
			return 0
		}
	case reflect.Array,
		reflect.Slice:
		if valueOf.Len() > 0 {
			return toInteger(valueOf.Index(0).Interface())
		}
	case reflect.Map:
		for _, key := range valueOf.MapKeys() {
			return toInteger(valueOf.MapIndex(key).Interface())
		}
	case reflect.Struct:
		if valueOf.IsValid() {
			field := valueOf.Elem().FieldByName("ID")
			if field.IsValid() {
				return field.Int()
			}
		}
	case reflect.Ptr:
		if valueOf.Elem().IsValid() {
			field := valueOf.FieldByName("ID")
			if field.IsValid() {
				return field.Int()
			}
		}
	}
	return 0
}

func toBoolean(value interface{}) bool {
	var valueOf = reflect.ValueOf(value)
	switch valueOf.Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return valueOf.Int() != 0
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr:
		return valueOf.Uint() != 0
	case reflect.Float32,
		reflect.Float64:
		return valueOf.Float() != 0
	case reflect.String:
		s := valueOf.String()
		return s == "b" || s == "Y"
	case reflect.Bool:
		return valueOf.Bool()
	case reflect.Array,
		reflect.Slice:
		if valueOf.Len() > 0 {
			return toBoolean(valueOf.Index(0).Interface())
		}
	case reflect.Map:
		for _, key := range valueOf.MapKeys() {
			return toBoolean(valueOf.MapIndex(key).Interface())
		}
	}
	return false
}

func toString(value interface{}, maxLength int) (result string) {
	var valueOf = reflect.ValueOf(value)
	switch valueOf.Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		result = strconv.FormatInt(valueOf.Int(), 10)
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr:
		result = strconv.FormatUint(valueOf.Uint(), 10)
	case reflect.Float32,
		reflect.Float64:
		result = strconv.FormatFloat(valueOf.Float(), 'f', -1, 64)
	case reflect.String:
		result = valueOf.String()
	case reflect.Bool:
		result = strconv.FormatBool(valueOf.Bool())
	case reflect.Array,
		reflect.Slice:
		var buff []string
		for i := 0; i < valueOf.Len(); i++ {
			buff = append(buff, toString(valueOf.Index(i).Interface(), 0))
		}
		result = strings.Join(buff, ",")
	case reflect.Map:
		var buff []string
		for _, key := range valueOf.MapKeys() {
			buff = append(buff, toString(valueOf.MapIndex(key).Interface(), 0))
		}
		result = strings.Join(buff, ",")
	}
	if maxLength > 0 && maxLength < len([]rune(result)) {
		return string([]rune(result)[:maxLength])
	}
	return result
}

func toFloat(value interface{}) float64 {
	var valueOf = reflect.ValueOf(value)
	switch valueOf.Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return float64(valueOf.Int())
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr:
		return float64(valueOf.Uint())
	case reflect.Float32,
		reflect.Float64:
		return valueOf.Float()
	case reflect.String:
		num, _ := strconv.ParseFloat(valueOf.String(), 64)
		return num
	case reflect.Bool:
		if valueOf.Bool() {
			return 1
		} else {
			return 0
		}
	case reflect.Array,
		reflect.Slice:
		if valueOf.Len() > 0 {
			return toFloat(valueOf.Index(0).Interface())
		}
	case reflect.Map:
		for _, key := range valueOf.MapKeys() {
			return toFloat(valueOf.MapIndex(key).Interface())
		}
	}
	return 0
}

func toTime(value interface{}) (pqt xpgtypes.NullTime) {
	if t, ok := value.(time.Time); ok {
		return xpgtypes.NullTime{
			Valid: true,
			Time:  t,
		}
	}

	var str = strings.TrimSpace(toString(value, 0))

	if str == "" {
		return
	}

	var err error
	// Unix time
	if regexp.MustCompile(`^[0-9]+$`).MatchString(str) {
		timestamp, err := strconv.ParseInt(str, 10, 64)
		pqt.Time = time.Unix(timestamp, 0)
		pqt.Valid = err == nil
		pqt.Error = err
		// 2006-01-02 15:04:05 -0700
	} else if regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} [+|-]{1}[0-9]{1,4}$`).MatchString(str) {
		pqt.Time, err = time.Parse("2006-01-02 15:04:05 -0700", str)
		pqt.Valid = err == nil
		pqt.Error = err
		// 2006-01-02 15:04:05
	} else if regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}$`).MatchString(str) {
		pqt.Time, err = time.Parse("2006-01-02 15:04:05", str)
		pqt.Valid = err == nil
		pqt.Error = err
		// 2006-01-02 15:04 -0700
	} else if regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2} [+|-]{1}[0-9]{1,4}$`).MatchString(str) {
		pqt.Time, err = time.Parse("2006-01-02 15:04 -0700", str)
		pqt.Valid = err == nil
		pqt.Error = err
		// 2006-01-02 15:04
	} else if regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}$`).MatchString(str) {
		pqt.Time, err = time.Parse("2006-01-02 15:04", str)
		pqt.Valid = err == nil
		pqt.Error = err
		// 2006-01-02
	} else if regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}$`).MatchString(str) {
		pqt.Time, err = time.Parse("2006-01-02", str)
		pqt.Valid = err == nil
		pqt.Error = err
		// 02.01.2006 15:04:05
	} else if regexp.MustCompile(`^[0-9]{2}\.[0-9]{2}\.[0-9]{4} [0-9]{2}:[0-9]{2}:[0-9]{2}$`).MatchString(str) {
		pqt.Time, err = time.Parse("02.01.2006 15:04:05", str)
		pqt.Valid = err == nil
		pqt.Error = err
		// 02.01.2006 15:04
	} else if regexp.MustCompile(`^[0-9]{2}\.[0-9]{2}\.[0-9]{4} [0-9]{2}:[0-9]{2}$`).MatchString(str) {
		pqt.Time, err = time.Parse("02.01.2006 15:04", str)
		pqt.Valid = err == nil
		pqt.Error = err
		// 02.01.2006
	} else if regexp.MustCompile(`^[0-9]{2}\.[0-9]{2}\.[0-9]{4}$`).MatchString(str) {
		pqt.Time, err = time.Parse("02.01.2006", str)
		pqt.Valid = err == nil
		pqt.Error = err
		// 2006.01.02
	} else if regexp.MustCompile(`^[0-9]{4}\.[0-9]{2}\.[0-9]{2}$`).MatchString(str) {
		pqt.Time, err = time.Parse("2006.01.02", str)
		pqt.Valid = err == nil
		pqt.Error = err
	}

	return
}

func toArray(value interface{}) []interface{} {
	var result []interface{}
	var valueOf = reflect.ValueOf(value)

	switch valueOf.Kind() {
	case reflect.Array,
		reflect.Slice:
		for i := 0; i < valueOf.Len(); i++ {
			result = append(result, valueOf.Index(i).Interface())
		}
	case reflect.Map:
		for _, key := range valueOf.MapKeys() {
			result = append(result, valueOf.MapIndex(key).Interface())
		}
	case reflect.String:
		result = append(result, valueOf.String())
	default:
		result = append(result, value)
	}
	return result
}

func inArray(arr []string, needle string) bool {
	for _, item := range arr {
		if item == needle {
			return true
		}
	}
	return false
}
