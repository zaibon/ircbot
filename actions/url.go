package actions

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Zaibon/ircbot"
	db "github.com/Zaibon/ircbot/database"
)

type URLLog struct {
	db *db.DB
}

func NewURLLog(bot *ircbot.IrcBot) *URLLog {
	conn, err := bot.DBConnection()
	if err != nil {
		panic(err)
	}

	initDB(conn)
	return &URLLog{
		db: conn,
	}
}

func initDB(db *db.DB) {
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS urls(
		id INTEGER CONSTRAINT url_PK PRIMARY KEY,
		nick STRING,
		url TEXT,
		timestamp INTEGER)`); err != nil {

		panic(err)
	}
}

func (u *URLLog) Command() []string {
	return []string{
		"PRIVMSG",
	}
}

func (u *URLLog) Usage() string {
	return ""
}

func (u *URLLog) Do(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
	for _, word := range m.Trailing {

		if !strings.Contains(word, "http") {
			continue
		}

		URL, err := url.Parse(word)
		if err != nil {
			fmt.Println("ERROR: URLLog parse url failed: ", err)
			continue
		}
		insertUrl(URL.String(), m.Nick(), u.db)
	}
}

func insertUrl(url, nick string, db *db.DB) {
	sql := "INSERT INTO urls(nick,url,timestamp) VALUES ($nick,$url,$timestamp)"
	if err := db.Exec(sql, nick, url, time.Now()); err != nil {
		fmt.Printf("ERROR: insert url failed, %s\n", err)
	}
	fmt.Printf("INFO: insert url(%s) succeed\n", url)
}

type URL struct {
	db *db.DB
}

func NewURL(bot *ircbot.IrcBot) *URL {
	conn, err := bot.DBConnection()
	if err != nil {
		panic(err)
	}

	initDB(conn)
	return &URL{
		db: conn,
	}
}

func (act *URL) initDB(db *db.DB) {
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS urls(
		id INTEGER CONSTRAINT url_PK PRIMARY KEY,
		nick STRING,
		url TEXT,
		timestamp INTEGER)`); err != nil {

		panic(err)
	}
}

func (u *URL) Command() []string {
	return []string{
		".url",
	}
}

func (u *URL) Usage() string {
	return ".url :args"
}

func (u *URL) Do(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
	sql := "SELECT url FROM urls "
	limit := 5

	if len(m.Trailing) > 1 {
		q := strings.Join(m.Trailing[1:], " ")
		sql = sql + " WHERE url LIKE '%" + q + "%' "
		limit = 10
	}

	sql = sql + fmt.Sprintf(" ORDER BY timestamp DESC LIMIT %d ", limit)

	stmt, err := u.db.Query(sql)
	if err != nil {
		fmt.Printf("ERROR query db :%s", err)
		return
	}

	for ; err == nil; err = stmt.Next() {
		var url string
		stmt.Scan(&url)
		b.Say(m.Channel(), url)
	}

	if err := stmt.Close(); err != nil {
		fmt.Printf("ERROR commit : %s\n", err)
	}
}
