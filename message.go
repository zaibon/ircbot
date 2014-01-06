package ircbot

import (
	"strings"
)

type IrcMsg struct {
	Raw    string // [':' <préfixe> <espace> ] <command> <params> <crlf>
	Prefix string // <nom de serveur> | <pseudo> [ '!' <utilisateur> ] [ '@' <hôte> ]
	Nick   string

	Command string // <lettre> { <lettre> } | <nombre> <nombre> <nombre>

	Args []string // <espace> [ ':' <fin> | <milieu> <params> ]

	Channel string
}

func NewIrcMsg() *IrcMsg {
	return &IrcMsg{}
}

func (m *IrcMsg) Parseline(line string) {
	m.Raw = line

	fields := strings.Fields(line)

	if strings.HasPrefix(line, ":") {
		//action of a user

		m.Prefix = fields[0]

		i := strings.Index(m.Prefix, "!")
		if i > 1 {
			m.Nick = m.Prefix[1:i]
		}
		m.Command = fields[1]
		if len(fields) >= 2 {
			m.Channel = strings.TrimPrefix(fields[2], ":")
			m.Args = fields[3:]
		}
	} else {
		//message send from the server
		m.Prefix = ""
		m.Command = fields[0]
		m.Args = fields[1:]
	}

}

func Parseline(line string) *IrcMsg {
	msg := NewIrcMsg()
	msg.Parseline(line)
	return msg
}
