package actions

import (
	"fmt"
	"strings"
	"testing"

	"github.com/zaibon/ircbot"
)

func TestURLLog(t *testing.T) {
	rd, wr, ln, _ := initActionnerTest(t)
	// defer bot.Disconnect()
	defer ln.Close()
	// fmt.Println(b.GetActionnersCmds())

	wr.PrintfLine(":user!~test@test.com PRIVMSG #test :http://google.com")
	wr.PrintfLine(":user!~test@test.com PRIVMSG #test :.url google")
	line, err := rd.ReadLine()
	if err != nil {
		t.Errorf("error reading server: %s", err)
	}
	msg := ircbot.ParseLine(line)
	if msg.Command != "PRIVMSG" {
		t.Errorf("command expected PRIVMSG actual %s", msg.Command)
	}
	if msg.Channel() != "#test" {
		t.Errorf("channel expected #test actual %s", msg.CmdParams[0])
	}

	args := strings.Join(msg.Trailing, " ")
	if !strings.Contains(args, "http://google.com") {
		t.Errorf("msg expected: http://google.com\nactual: %s", args)

	}

	wr.PrintfLine(":user!~test@test.com PRIVMSG #test :http://google.com")
	wr.PrintfLine(":user!~test@test.com PRIVMSG #test :.url")
	line, err = rd.ReadLine()
	if err != nil {
		t.Errorf("error reading server: %s", err)
	}
	msg = ircbot.ParseLine(line)
	args = strings.Join(msg.Trailing, " ")
	if !strings.Contains(args, "hit 2 times") {
		t.Errorf("msg expected: hit 2 times\nactual: %s", args)

	}

	wr.PrintfLine(":user!~test@test.com PRIVMSG #test-no-the-same :.url google")
	line, err = rd.ReadLine()
	if err != nil {
		t.Errorf("error reading server: %s", err)
	}
	msg = ircbot.ParseLine(line)
	args = strings.Join(msg.Trailing, " ")
	fmt.Println(args)
	if args != "no results" {
		t.Errorf("should not find url from other channel")
	}
}
