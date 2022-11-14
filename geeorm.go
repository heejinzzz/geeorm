package geeorm

import (
	"database/sql"
	"fmt"
	"github.com/heejinzzz/geeorm/dialect"
	"github.com/heejinzzz/geeorm/log"
	"github.com/heejinzzz/geeorm/session"
	"strings"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver string, source string) (*Engine, error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// Send a ping to make sure the database connection is alive
	if err = db.Ping(); err != nil {
		log.Error(err)
		return nil, err
	}

	// make sure the specific dialect exists
	dlct, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return nil, nil
	}

	engine := &Engine{db: db, dialect: dlct}
	log.Info("Connect database success")
	return engine, nil
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Errorf("Failed to close database: ", err)
	} else {
		log.Info("Close database success")
	}
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}

type TxFunc func(*session.Session) (interface{}, error)

// Transaction create and begin a transaction.
// Any errors occur, rollback automatically,
// or commit if no errors occur.
func (e *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := e.NewSession()
	if err = s.Begin(); err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p) // re-throw panic after rollback
		} else if err != nil {
			_ = s.Rollback() // err is not nil, rollback and return err
		} else {
			err = s.Commit() // try to commit transaction
			if err != nil {
				_ = s.Rollback() // commit failed, rollback and return commit error
			}
		}
	}()

	return f(s)
	// "result, err = f(s)" => "defer func() {...}()" => "return"
}

// difference returns a - b
func difference(a, b []string) (diff []string) {
	mp := map[string]bool{}
	for _, v := range b {
		mp[v] = true
	}
	for _, v := range a {
		if !mp[v] {
			diff = append(diff, v)
		}
	}
	return
}

// Migrate add or delete columns of a table
func (e *Engine) Migrate(value interface{}) error {
	_, err := e.Transaction(func(s *session.Session) (interface{}, error) {
		if !s.Model(value).HasTable() {
			log.Errorf("table %s not exists", s.RefTable().Name)
			return nil, fmt.Errorf("table %s not exists", s.RefTable().Name)
		}
		table := s.RefTable()
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).Query()
		columns, _ := rows.Columns()
		_ = rows.Close()
		addCols := difference(table.FieldNames, columns)
		delCols := difference(columns, table.FieldNames)
		log.Infof("add cols %v, delete cols %v", addCols, delCols)

		for _, col := range addCols {
			field := table.GetField(col)
			sqlStr := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s %s;", table.Name, field.Name, field.Type, field.Tag)
			if _, err := s.Raw(sqlStr).Exec(); err != nil {
				return nil, err
			}
		}

		if len(delCols) == 0 {
			return nil, nil
		}
		temp := "temp_" + table.Name
		fields := strings.Join(table.FieldNames, ", ")
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s FROM %s;", temp, fields, table.Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", temp, table.Name))
		_, err := s.Exec()
		return nil, err
	})
	return err
}
