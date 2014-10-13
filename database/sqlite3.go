package ircbot

import (
	"io"
	"log"
	"os"

	"github.com/mxk/go-sqlite/sqlite3"
)

var logger = log.New(os.Stdout, "DATABASE :", log.Ldate|log.Ltime)

type DB struct {
	conn *sqlite3.Conn
}

func Open(path string) (*DB, error) {
	conn, err := sqlite3.Open(path)
	if err != nil {
		logger.Println("error opening database %s : %s\n", path, err.Error())
		return nil, err
	}

	DB := &DB{}
	DB.conn = conn

	return DB, err
}

func (d *DB) Path(name ...string) string {
	n := ""
	if name == nil {
		n = "main"
	} else {
		n = name[0]
	}
	return d.conn.Path(n)
}

func (d *DB) Close() error {
	return d.conn.Close()
}

func (d *DB) Begin() error {
	return d.conn.Begin()
}

func (d *DB) Commit() error {
	return d.conn.Commit()
}

func (d *DB) Query(sql string, args ...interface{}) (*sqlite3.Stmt, error) {
	logger.Printf("QUERY %s", sql)
	stmt, err := d.conn.Query(sql, args...)
	if err != nil && err != io.EOF {
		logger.Printf("error query : %s : %s\n", sql, err.Error())
		return nil, err
	}
	return stmt, err
}

func (d *DB) Exec(sql string, args ...interface{}) error {
	logger.Printf("EXEC %s", sql)

	if err := d.conn.Begin(); err != nil {
		logger.Printf("error begin exec :%s\n", err)
		return err
	}

	err := d.conn.Exec(sql, args...)
	if err != nil {
		logger.Printf("error exec : %s : %s", sql, err.Error())
		if err := d.conn.Rollback(); err != nil {
			logger.Printf("error rollback exec : %s", err.Error())
			return err
		}
	}

	if err := d.conn.Commit(); err != nil {
		logger.Printf("error commit exec : %s", err.Error())
	}

	return err
}

func (d *DB) Prepare(sql string) (*sqlite3.Stmt, error) {
	logger.Println("PREPARE %s", sql)
	return d.conn.Prepare(sql)
}
