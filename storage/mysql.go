//
//  mysql.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package storage

import (
	"log"
	"fmt"
	"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ortuman/mercury/config"
	"github.com/ortuman/mercury/logger"
)

const maxMySqlTransactionRetries = 10

// create table statements

const mySqlCreateUsersStatement = "" +
	"CREATE TABLE IF NOT EXISTS users (" +
	"  user_id VARCHAR(128) NOT NULL," +
	"  sender_id VARCHAR(16) CHARSET ascii NOT NULL," +
	"  token TEXT NOT NULL," +
	"  badge INT DEFAULT 0," +
	"  status VARCHAR(32) NOT NULL," +
	"  created_at DATETIME NOT NULL," +
	"  updated_at DATETIME NOT NULL," +
	" PRIMARY KEY (user_id, sender_id)) DEFAULT CHARSET = utf8mb4"

// create index statements
const createUsersIdIndexStatement = "CREATE INDEX users_id_idx ON users (user_id)"

type MySql struct {
	db *sql.DB
}

func NewMySql() *MySql {
	s := new(MySql)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/", config.MySql.User, config.MySql.Password, config.MySql.Host)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("mysql: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("mysql: %v", err)
	}
	s.db = db

	if _, err := s.db.Exec("SET NAMES utf8mb4 COLLATE utf8mb4_unicode_ci"); err != nil {
		log.Fatalf("mysql: %v", err)
	}

	if err := s.createDatabase(); err != nil {
		log.Fatalf("mysql: %v", err)
	}
	if err := s.createTables(); err != nil {
		log.Fatalf("mysql: %v", err)
	}

	go s.performAdditionalTasks()
	return s
}

func (s *MySql) InsertSenderInfo(senderInfo *SenderInfo) error {
	return nil
}

func (s *MySql) FetchSenderInfo(userID string, senderID string) (*SenderInfo, error) {
	return nil, nil
}

func (s *MySql) FetchSenderInfoArray(userID string) ([]*SenderInfo, error) {
	return []*SenderInfo{}, nil
}

func (s *MySql) DeleteSenderInfo(userID string, senderID string) error {
	return nil
}

func (s *MySql) createDatabase() error {
	if _, err := s.db.Exec("CREATE DATABASE IF NOT EXISTS mercury"); err != nil {
		return err
	}
	if _, err := s.db.Exec("USE mercury"); err != nil {
		return err
	}
	return nil
}

func (s *MySql) inTransaction(function func(*sql.Tx) error) error {
	var err error
	for i := 0; i < maxMySqlTransactionRetries; i++ {
		tx, txErr := s.db.Begin()
		if txErr != nil {
			return txErr
		}
		tx.Exec("SET NAMES utf8mb4 COLLATE utf8mb4_unicode_ci")
		err = function(tx)
		if err == nil {
			tx.Commit()
			return nil
		} else {
			tx.Rollback()
			continue
		}
	}
	return err
}

func (s *MySql) performAdditionalTasks() {
	t1 := time.NewTicker(time.Hour * 24)
	for {
		select {
		case <-t1.C:
			startT := time.Now().UnixNano()
			err := s.optimizeTables()
			if err != nil {
				logger.Errorf("mysql: couldn't optimize tables: %v", err)
			} else {
				endT := time.Duration((time.Now().UnixNano() - startT)) / time.Millisecond
				logger.Infof("mysql: table optimization completed (%dms)", endT)
			}
		}
	}
}

func (s *MySql) optimizeTables() error {

	// analyze tables
	s.db.Exec("ANALYZE TABLE users")

	return nil
}

func (s *MySql) createTables() error {

	// create tables
	if _, err := s.db.Exec(mySqlCreateUsersStatement); err != nil 	{ return err }

	// create additional indexes
	s.db.Exec(createUsersIdIndexStatement)

	return nil
}
