package ircbot

import (
	"io"
	"log"
	"os"
	"strings"
	"time"

	"code.google.com/p/go-sqlite/go1/sqlite3"
)

type DB struct {
	db  *sqlite3.Conn
	log *log.Logger
}

func newDB() *DB {
	l := log.New(os.Stdout, "DATABASE :", log.Ldate|log.Ltime|log.Lshortfile)
	return &DB{
		log: l,
	}
}

func (d *DB) init() error {
	if d.db == nil {
		d.log.Printf("should open connection before initialize database\n")
	}

	if err := d.exec(`CREATE TABLE IF NOT EXISTS logs(
		id INTEGER CONSTRAINT line_PK PRIMARY KEY,
		nick STRING,
		message TEXT,
		channel STRING,
		timestamp INTEGER)`); err != nil {

		return err
	}
	return nil
}

func (d *DB) open(path string) error {
	db, err := sqlite3.Open(path)
	if err != nil {
		d.log.Println("error opening database %s : %s\n", path, err.Error())
	}
	d.db = db
	return err
}

func (d *DB) close() error {
	return d.db.Close()
}

func (d *DB) query(sql string, args ...interface{}) (*sqlite3.Stmt, error) {
	stmt, err := d.db.Query(sql, args...)
	if err != nil && err != io.EOF {
		d.log.Printf("error query : %s : %s\n", sql, err.Error())
	}
	return stmt, err
}

func (d *DB) exec(sql string, args ...interface{}) error {
	err := d.db.Exec(sql, args...)
	if err != nil {
		d.log.Printf("error exec : %s : %s", sql, err.Error())
	}
	return err
}

func logMsg(m *IrcMsg, db *DB) error {
	sql := "INSERT INTO logs (nick,message,channel,timestamp) VALUES ($nick,$message,$channel,$timestamp)"

	msg := strings.Join(m.Args, " ")
	if err := db.exec(sql, m.Nick, msg, m.Channel, time.Now()); err != nil {
		db.log.Printf("error inserting logs : %s", err.Error())
		return err
	}
	return nil
}
