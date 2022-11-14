package dialect

import (
	"fmt"
	"reflect"
	"time"
)

// provide dialect for sqlite3
type sqlite3 struct{}

func (*sqlite3) DataTypeOf(value reflect.Value) string {
	switch value.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := value.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("unsupported type: %s (%s)", value.Type().Name(), value.Kind()))
}

func (*sqlite3) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "SELECT name FROM sqlite_master WHERE type = 'table' and name = ?", args
}
