package ircbot

import (
	"strings"
)

type ActionFunc func(b *IrcBot, m *IrcMsg)

type Action struct {
	command string
	do      ActionFunc
}

func Pong(b *IrcBot, m *IrcMsg) {
	b.writer.PrintfLine("PONG %s", strings.Join(m.args, " "))
}
