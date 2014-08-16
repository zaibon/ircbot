package ircbot

import (
	"io"
	"log"
	"os"

	"code.google.com/p/go-sqlite/go1/sqlite3"
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

// func (d *DB) init() error {
// 	if d.db == nil {
// 		d.log.Printf("should open connection before initialize database\n")
// 	}

// 	if err := d.Exec(`CREATE TABLE IF NOT EXISTS logs(
// 		id INTEGER CONSTRAINT line_PK PRIMARY KEY,
// 		nick STRING,
// 		message TEXT,
// 		channel STRING,
// 		timestamp INTEGER)`); err != nil {

// 		return err
// 	}
// 	return nil
// }

func (d *DB) Close() error {
	return d.conn.Close()
}

func (d *DB) Query(sql string, args ...interface{}) (*sqlite3.Stmt, error) {
	logger.Printf("QUERY %s", sql)
	stmt, err := d.conn.Query(sql, args...)
	if err != nil && err != io.EOF {
		logger.Printf("error query : %s : %s\n", sql, err.Error())
	}
	return stmt, err
}

func (d *DB) Exec(sql string, args ...interface{}) error {
	logger.Printf("EXEC %s", sql)
	err := d.conn.Exec(sql, args...)
	if err != nil {
		logger.Printf("error exec : %s : %s", sql, err.Error())
	}
	return err
}

// func logMsg(m *IrcMsg, db *DB) error {
// 	sql := "INSERT INTO logs (nick,message,channel,timestamp) VALUES ($nick,$message,$channel,$timestamp)"

// 	msg := strings.Join(m.Trailing, " ")
// 	if err := db.Exec(sql, m.Nick(), msg, m.Channel(), time.Now()); err != nil {
// 		db.log.Printf("error inserting logs : %s", err.Error())
// 		return err
// 	}
// 	return nil
// }
