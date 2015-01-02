package actions

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"os"
	"testing"

	"github.com/zaibon/ircbot"
)

func TestMain(m *testing.M) {
	exit := m.Run()
	if err := os.Remove("irc_test.db"); err != nil {
		fmt.Println("error deleting database")
	}
	os.Exit(exit)
}

func initActionnerTest(t *testing.T) (*textproto.Reader, *textproto.Writer, net.Listener, *ircbot.IrcBot) {
	ln, err := net.Listen("tcp", ":3333")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("test server listen on %s\n", ln.Addr().String())

	b := ircbot.NewIrcBot("testBot", "testBot", "", "127.0.0.1", 3333, []string{"#test"}, "irc_test.db")
	b.AddInternAction(&Greet{})
	b.AddInternAction(NewURLLog(b))
	b.AddInternAction(NewTitleExtract())
	b.AddUserAction(&Ping{})
	b.AddUserAction(NewURL(b))

	b.Connect()

	fmt.Printf("test server waiting connection...\n")
	conn, err := ln.Accept()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("accepted %s\n", conn.RemoteAddr())

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	rd := textproto.NewReader(r)
	wr := textproto.NewWriter(w)

	for {
		line, err := rd.ReadLine()
		if err != nil {
			t.Errorf("error reading server : %s\n", err)
		}
		msg := ircbot.ParseLine(line)

		if msg.Command == "NICK" {
			//NICK is the last command send by the bot when connecting
			//we can send him message now
			break
		}
	}
	return rd, wr, ln, b
}
