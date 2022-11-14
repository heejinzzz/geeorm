package session

import (
	"database/sql"
	"github.com/heejinzzz/geeorm/dialect"
	"testing"
)

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

func (*User) TableName() string {
	return "user"
}

func TestSession_CreateTable(t *testing.T) {
	db, _ := sql.Open("sqlite3", "gee.db")
	d, _ := dialect.GetDialect("sqlite3")
	user := &User{}
	s := New(db, d).Model(user)
	_ = s.DropTable()
	_ = s.CreateTable()
	if !s.HasTable() {
		t.Fatal("Failed to create table user")
	}
}
