package session

import "github.com/heejinzzz/geeorm/log"

// Begin opens a transaction
func (s *Session) Begin() error {
	var err error
	if s.tx, err = s.db.Begin(); err != nil {
		log.Error("open transaction failed: ", err)
		return err
	}
	log.Info("transaction begin")
	return nil
}

// Commit commits a transaction
func (s *Session) Commit() error {
	if err := s.tx.Commit(); err != nil {
		log.Error("commit transaction failed: ", err)
		return err
	}
	log.Info("transaction commit")
	return nil
}

// Rollback rollbacks a transaction
func (s *Session) Rollback() error {
	if err := s.tx.Rollback(); err != nil {
		log.Error("rollback transaction failed: ", err)
		return err
	}
	log.Info("transaction rollback")
	return nil
}
