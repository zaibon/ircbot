package actions

import (
	"fmt"
	"strings"
	"time"

	"github.com/mxk/go-sqlite/sqlite3"

	"github.com/zaibon/ircbot"
	db "github.com/zaibon/ircbot/database"
)

type Logger struct {
	db *db.DB
}

var insertStmt *sqlite3.Stmt

func NewLogger(b *ircbot.IrcBot) *Logger {
	conn, err := b.DBConnection()
	if err != nil {
		panic(err)
	}

	if err := conn.Exec(`CREATE TABLE IF NOT EXISTS logs(
		id INTEGER CONSTRAINT line_PK PRIMARY KEY,
		nick STRING,
		message TEXT,
		channel STRING,
		timestamp INTEGER)`); err != nil {
		panic(err)
	}

	insertStmt, err = conn.Prepare("INSERT INTO logs (nick,message,channel,timestamp) VALUES ($nick,$message,$channel,$timestamp)")
	if err != nil {
		panic(err)
	}

	return &Logger{
		db: conn,
	}
}

func (l *Logger) Command() []string {
	return []string{
		"PRIVMSG",
	}
}

func (l *Logger) Usage() string {
	return ""
}

func (l *Logger) Do(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
	msg := strings.Join(m.Trailing, " ")
	if err := insertStmt.Exec(m.Nick(), msg, m.Channel(), time.Now()); err != nil {
		fmt.Printf("DATABASE error execute stmt : %s\n", err)
	}
}
