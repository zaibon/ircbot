package ircbot

import (
	"fmt"
	"strings"
)

type ActionFunc func(b *IrcBot, m *IrcMsg)

type Action struct {
	command string
	do      ActionFunc
}

//Pong sends a pong response to ping
func Pong(b *IrcBot, m *IrcMsg) {

	s := fmt.Sprintf("PONG %s", strings.Join(m.Args, " "))
	fmt.Println("irc >> ", s)
	b.writer.PrintfLine(s)
}

func ValidConnect(b *IrcBot, m *IrcMsg) {
	//channel is not a good name in this case
	//MODE command put nick at channel place
	if strings.Contains(m.Channel, b.Nick) {
		fmt.Println("Info : connection terminated")
		b.Joined = true
	}
}
