package ircbot

type ActionFunc func(b *IrcBot, m *IrcMsg)

type Actioner interface {
	Command() []string
	Usage() string
	Do(b *IrcBot, m *IrcMsg)
}
