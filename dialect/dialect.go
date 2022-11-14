package dialect

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"reflect"
)

var dialectsMap = map[string]Dialect{}

type Dialect interface {
	DataTypeOf(reflect.Value) string
	TableExistSQL(tableName string) (string, []interface{})
}

func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}

func init() {
	RegisterDialect("sqlite3", &sqlite3{})
	RegisterDialect("mysql", &mysql{})
}
