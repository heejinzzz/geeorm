package dialect

import (
	"fmt"
	"reflect"
	"time"
)

// provide dialect for mysql
type mysql struct{}

func (*mysql) DataTypeOf(value reflect.Value) string {
	switch value.Kind() {
	case reflect.Bool:
		return "tinyint(1)"
	case reflect.Int, reflect.Int32:
		return "int"
	case reflect.Uint, reflect.Uint32, reflect.Uintptr:
		return "int unsigned"
	case reflect.Int8:
		return "tinyint"
	case reflect.Uint8:
		return "tinyint unsigned"
	case reflect.Int16:
		return "smallint"
	case reflect.Uint16:
		return "smallint unsigned"
	case reflect.Int64:
		return "bigint"
	case reflect.Uint64:
		return "bigint unsigned"
	case reflect.Float32:
		return "float"
	case reflect.Float64:
		return "double"
	case reflect.String:
		return "varchar(768)"
	case reflect.Struct:
		if _, ok := value.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("unsupported type: %s (%s)", value.Type().Name(), value.Kind()))
}

func (*mysql) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "select TABLE_NAME from information_schema.TABLES where TABLE_SCHEMA = (select database()) and TABLE_NAME = ?", args
}
