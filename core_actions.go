package ircbot

import (
	"fmt"
	"strings"
)

type ActionFunc func(b *IrcBot, m *IrcMsg)

type Actioner interface {
	Command() []string
	Usage() string
	Do(b *IrcBot, m *IrcMsg)
}

//Pong sends a pong response to ping
type pong struct{}

func (p *pong) Command() []string {
	return []string{"PING"}
}

func (p *pong) Usage() string {
	return "pong answer to ping request"
}

func (p *pong) Do(b *IrcBot, m *IrcMsg) {
	s := fmt.Sprintf("PONG %s", strings.Join(m.Args, " "))
	fmt.Println("irc >> ", s)
	b.writer.PrintfLine(s)
}

//ValidConnect sets a flag to true is the connection succeed
type validConnect struct{}

func (v *validConnect) Command() []string {
	return []string{"MODE"}
}

func (v *validConnect) Usage() string {
	return "ValidConnect set a flag to true is the connection succeed"
}

func (v *validConnect) Do(b *IrcBot, m *IrcMsg) {
	//channel is not a good name in this case
	//MODE command put nick at channel place
	if strings.Contains(m.Channel, b.Nick) {
		fmt.Println("Info : connection terminated")
		b.joined = true
	}
}

type Help struct{}

func (h *Help) Command() []string {
	return []string{
		".help",
		".h",
	}
}

func (h *Help) Usage() string {
	return ".help .h : display this message"
}

func (h *Help) Do(b *IrcBot, m *IrcMsg) {
	var output string

	for cmd, _ := range b.handlersUser {
		output += cmd + ", "
	}
	b.Say(m.Channel, output[:len(output)-2])
}
