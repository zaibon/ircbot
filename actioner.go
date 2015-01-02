package ircbot

//ActionFunc is the type of function used in Actioner
type ActionFunc func(b *IrcBot, m *IrcMsg)

//Actioner is the interface that objects need to implement custom action for the bot
type Actioner interface {
	// array of IRC or user command on which the action will be triggered (MSG)
	// example for IRC command: JOIN, PRIVMSG
	// example for user command: .myCmd
	Command() []string

	// used to display help about the action
	Usage() string

	// actual action done by the bot
	Do(b *IrcBot, m *IrcMsg)
}
