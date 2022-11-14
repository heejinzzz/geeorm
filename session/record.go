package session

import (
	"errors"
	"github.com/heejinzzz/geeorm/clause"
	"reflect"
)

// Insert record
func (s *Session) Insert(values ...interface{}) (int64, error) {
	// BeforeInsert hook
	for _, value := range values {
		s.CallMethod(BeforeInsert, reflect.ValueOf(value).Interface())
	}
	var recordValues []interface{}
	for _, value := range values {
		table := s.Model(value).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}

	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	// AfterInsert hook
	for _, value := range values {
		s.CallMethod(AfterInsert, reflect.ValueOf(value).Interface())
	}

	return result.RowsAffected()
}

// Find records
func (s *Session) Find(values interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()
	s.CallMethod(BeforeQuery, nil) // BeforeQuery hook

	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).Query()
	if err != nil {
		return err
	}

	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var fieldValues []interface{}
		for _, field := range table.FieldNames {
			fieldValues = append(fieldValues, dest.FieldByName(field).Addr().Interface())
		}
		if err := rows.Scan(fieldValues...); err != nil {
			return err
		}
		s.CallMethod(AfterQuery, dest.Addr().Interface()) // AfterQuery hook
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}

// Update records with Where clause
// support map[string]interface{}
// also support kv list: "Name", "Tom", "Age", 18, ....
func (s *Session) Update(kv ...interface{}) (int64, error) {
	s.CallMethod(BeforeUpdate, nil) // BeforeUpdate hook
	m, ok := kv[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}
	s.clause.Set(clause.UPDATE, s.refTable.Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterUpdate, nil) // AfterUpdate hook
	return result.RowsAffected()
}

// Delete records with Where clause
func (s *Session) Delete() (int64, error) {
	s.CallMethod(BeforeDelete, nil) // BeforeDelete hook
	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterDelete, nil) // AfterDelete hook
	return result.RowsAffected()
}

// Count records with Where clause
func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var count int64
	if err := row.Scan(&count); err != nil {
		return -1, err
	}
	return count, nil
}

// Limit adds LIMIT condition to clause
func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

// Where adds WHERE condition to clause
func (s *Session) Where(desc string, args ...interface{}) *Session {
	s.clause.Set(clause.WHERE, append([]interface{}{desc}, args...)...)
	return s
}

// OrderBy adds ORDER BY condition to clause
func (s *Session) OrderBy(fields ...string) *Session {
	var vars []interface{}
	for _, field := range fields {
		vars = append(vars, field)
	}
	s.clause.Set(clause.ORDERBY, vars...)
	return s
}

// First gets the first record that matches
func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	dest.Set(destSlice.Index(0))
	return nil
}
