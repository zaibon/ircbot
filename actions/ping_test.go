package actions

import (
	"strings"

	"github.com/zaibon/ircbot"

	"testing"
)

func TestPing(t *testing.T) {
	rd, wr, ln, _ := initActionnerTest(t)
	// defer bot.Disconnect()
	defer ln.Close()

	// wr.PrintfLine(":guestUser!guestUser PRIVMSG #test :.ping")
	wr.PrintfLine(":guestUser!guestUser PRIVMSG #test :.ping")
	line, err := rd.ReadLine()
	if err != nil {
		t.Errorf("error reading server: %s", err)
	}
	msg := ircbot.ParseLine(line)
	if msg.Command != "PRIVMSG" {
		t.Errorf("command expected PRIVMSG actual %s", msg.Command)
	}
	if msg.CmdParams[0] != "#test" {
		t.Errorf("channel expected #test actual %s", msg.CmdParams[0])
	}
	args := strings.Join(msg.Trailing, " ")
	if args != "pong" {
		t.Errorf("msg expected Salut guestUser actual %s", args)
	}
}
