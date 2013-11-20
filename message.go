package ircbot

import (
	"strings"
)

type IrcMsg struct {
	raw    string // :prefix commad :args
	prefix string // Nick!user@host

	nick string

	command string
	args    []string
}

func NewIrcMsg() *IrcMsg {
	return &IrcMsg{}
}

func (m *IrcMsg) Parseline(line string) {
	m.raw = line

	fields := strings.Fields(line)

	if strings.HasPrefix(line, ":") {
		//message send from a user
		m.prefix = fields[0]
		i := strings.Index(m.prefix, "!")
		if i > 1 {
			m.nick = m.prefix[1:i]
		}
		m.command = fields[1]
		m.args = fields[2:]
	} else {
		//message send from the server
		m.prefix = ""
		m.command = fields[0]
		m.args = fields[1:]
	}

}

func Parseline(line string) *IrcMsg {
	msg := NewIrcMsg()
	msg.Parseline(line)
	return msg
}
