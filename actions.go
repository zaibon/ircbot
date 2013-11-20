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

	s := fmt.Sprintf("PONG %s", strings.Join(m.args, " "))
	fmt.Println("irc >> ", s)
	b.writer.PrintfLine(s)
}

func Join(b *IrcBot, m *IrcMsg) {
	if m.nick == b.Nick {
		b.joined = true
		return
	}

	s := fmt.Sprintf("%s :Salut %s", m.channel, m.nick)
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
		line := fmt.Sprintf("%s :%s", m.channel, response[nbr])
		b.Out <- &IrcMsg{
			command: "PRIVMSG",
			args:    []string{line},
		}
	}
}

func ValidConnect(b *IrcBot, m *IrcMsg) {
	s := strings.Join(m.args, " ")
	if strings.Contains(s, b.Nick) {
		fmt.Println("Info : connection terminated")
		b.joined = true
	}
}
