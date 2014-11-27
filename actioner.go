package ircbot

//ActionFunc is the type of function used in Actioner
type ActionFunc func(b *IrcBot, m *IrcMsg)

//Actionner is the interface that object need to implement custom action for the bot
type Actioner interface {
	Command() []string
	Usage() string
	Do(b *IrcBot, m *IrcMsg)
}
