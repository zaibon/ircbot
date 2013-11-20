package ircbot

import (
	"fmt"
	"math/rand"
	"strings"
)

type ActionFunc func(b *IrcBot, m *IrcMsg)

type Action struct {
	command string
	do      ActionFunc
}

func Pong(b *IrcBot, m *IrcMsg) {

	b.Out <- &IrcMsg{
		command: "PONG",
		args:    m.args,
	}
}

func Join(b *IrcBot, m *IrcMsg) {
	if m.nick == b.Nick {
		b.joined = true
		return
	}

	s := fmt.Sprintf("%s :Salut %s", b.Channel[0], m.nick)
	b.Out <- &IrcMsg{
		command: "PRIVMSG",
		args:    []string{s},
	}
}

func Respond(b *IrcBot, m *IrcMsg) {
	response := []string{
		"oui ?",
		"on parle de moi ?",
		"Je suis pas là",
	}

	s := strings.Join(m.args, " ")

	if strings.Contains(s, b.Nick) {
		nbr := rand.Intn(len(response))
		line := fmt.Sprintf("%s :%s", b.Channel[0], response[nbr])
		b.Out <- &IrcMsg{
			command: "PRIVMSG",
			args:    []string{line},
		}
	}
}
