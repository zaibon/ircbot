package actions

import (
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/mxk/go-sqlite/sqlite3"

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
		hit INTEGER,
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
	if m.Nick() == b.Nick {
		//don't listen to itself
		return
	}

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
	sql := "SELECT url FROM urls WHERE url=$url"
	q, err := db.Query(sql, url)
	if err != nil && err != io.EOF {
		fmt.Printf("ERROR: query url failed, %s\n", err)
		return
	}

	if err == io.EOF {
		//the url is not yet in the db
		sql = "INSERT INTO urls(nick,url,hit,timestamp) VALUES ($nick,$url,1,$timestamp)"
		if err := db.Exec(sql, nick, url, time.Now()); err != nil {
			fmt.Printf("ERROR: insert url failed, %s\n", err)
			return
		}

		fmt.Printf("INFO: insert url(%s) succeed\n", url)
		return
	}

	q.Close()
	//the url already exists, update hit counter
	sql = "UPDATE urls SET hit=hit+1 WHERE url=$url"
	if err := db.Exec(sql, url); err != nil {
		fmt.Printf("ERROR update url falied, %s\n", err)
		return
	}
	fmt.Printf("INFO: update url(%s) hit succeed\n", url)
}

var (
	stmtUpdate *sqlite3.Stmt
	stmtCount  *sqlite3.Stmt
)

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

func (u *URL) Command() []string {
	return []string{
		".url",
	}
}

func (u *URL) Usage() string {
	return ".url :args"
}

func (u *URL) Do(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
	sql := "SELECT url,hit,nick FROM urls "
	limit := 5 //hardcoded for now, maybe let the user choose a limit

	if len(m.Trailing) > 1 {
		q := strings.Join(m.Trailing[1:], " ")
		sql = sql + " WHERE url LIKE '%" + q + "%' "
		limit = 10
	}

	sql = sql + fmt.Sprintf(" ORDER BY timestamp DESC LIMIT %d ", limit)

	stmt, err := u.db.Query(sql)
	if err != nil && err != io.EOF {
		fmt.Printf("ERROR query db :%s", err)
		return
	}

	if err == io.EOF {
		b.Say(m.Channel(), "no results")
		return
	}

	for ; err == nil; err = stmt.Next() {
		var (
			url  string
			hit  int
			nick string
		)
		stmt.Scan(&url, &hit, &nick)
		output := fmt.Sprintf("%s (hit %d times - %s)", url, hit, nick)
		b.Say(m.Channel(), output)
	}

	if err := stmt.Close(); err != nil {
		fmt.Printf("ERROR commit : %s\n", err)
	}
}