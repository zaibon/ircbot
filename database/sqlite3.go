package ircbot

import (
	"io"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/mxk/go-sqlite/sqlite3"
)

type DB struct {
	conn *sqlite3.Conn
}

func Open(path string) (*DB, error) {
	conn, err := sqlite3.Open(path)
	if err != nil {
		log.WithFields(log.Fields{
			"path":  path,
			"error": err,
		}).Errorln("open database")
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
	log.WithFields(log.Fields{
		"query": flat(sql),
	}).Debugln("query")

	stmt, err := d.conn.Query(sql, args...)
	if err != nil && err != io.EOF {
		log.WithFields(log.Fields{
			"query": flat(sql),
			"error": err,
		}).Errorln("error query")
		return nil, err
	}
	return stmt, err
}

func (d *DB) Exec(sql string, args ...interface{}) error {
	log.WithFields(log.Fields{
		"query": flat(sql),
	}).Debugln("exec")

	if err := d.conn.Begin(); err != nil {
		log.WithFields(log.Fields{
			"query": flat(sql),
			"error": err,
		}).Errorln("begin exec error")
		return err
	}

	err := d.conn.Exec(sql, args...)
	if err != nil {
		log.WithFields(log.Fields{
			"query": flat(sql),
			"error": err,
		}).Errorln("error exec")

		if err := d.conn.Rollback(); err != nil {
			log.WithFields(log.Fields{
				"query": flat(sql),
				"error": err,
			}).Errorln("error rollback exec")
			return err
		}
	}

	if err := d.conn.Commit(); err != nil {
		log.WithFields(log.Fields{
			"query": flat(sql),
			"error": err,
		}).Errorln("error commit exec")
	}

	return err
}

func (d *DB) Prepare(sql string) (*sqlite3.Stmt, error) {
	log.WithField("query", flat(sql)).Debugln("prepare")
	return d.conn.Prepare(sql)
}

func flat(s string) string {
	return strings.Replace(s, "\n", " ", -1)
}
