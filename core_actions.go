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
type Pong struct{}

func (p *Pong) Command() []string {
	return []string{"PONG"}
}

func (p *Pong) Usage() string {
	return "pong answer to ping request"
}

func (p *Pong) Do(b *IrcBot, m *IrcMsg) {
	s := fmt.Sprintf("PONG %s", strings.Join(m.Args, " "))
	fmt.Println("irc >> ", s)
	b.writer.PrintfLine(s)
}

// func Pong(b *IrcBot, m *IrcMsg) {

// 	s := fmt.Sprintf("PONG %s", strings.Join(m.Args, " "))
// 	fmt.Println("irc >> ", s)
// 	b.writer.PrintfLine(s)
// }

type ValidConnect struct{}

func (v *ValidConnect) Command() []string {
	return []string{"MODE"}
}

func (v *ValidConnect) Usage() string {
	return "Valid connect set a flag to true is the connection succeed"
}

func (v *ValidConnect) Do(b *IrcBot, m *IrcMsg) {
	//channel is not a good name in this case
	//MODE command put nick at channel place
	if strings.Contains(m.Channel, b.Nick) {
		fmt.Println("Info : connection terminated")
		b.Joined = true
	}
}

// func ValidConnect(b *IrcBot, m *IrcMsg) {
// 	//channel is not a good name in this case
// 	//MODE command put nick at channel place
// 	if strings.Contains(m.Channel, b.Nick) {
// 		fmt.Println("Info : connection terminated")
// 		b.Joined = true
// 	}
// }
