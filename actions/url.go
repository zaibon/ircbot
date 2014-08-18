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

// type URLActionner struct {
// 	db *ircbot.DB
// }

// func (u *URLActionner) Command() []string {
// 	return []string{
// 		".url",
// 	}
// }

// func (u *URLActionner) Usage() string {
// 	return ""
// }

// func (u *URLActionner) Do(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
// 	sql := "SELECT url FROM urls"

// 	if len(m.Trailing) > 2 {
// 		q := strings.Join(m.Trailing[1:], " ")
// 		sql = sql + " WHERE url LIKE %" + q + "%"
// 	}

// 	for s, err := b.DB.Query(sql); err == nil; err = s.Next() {
// 		var url string
// 		s.Scan(&url)
// 		b.Say(m.Channel(), url)
// 	}

// }
