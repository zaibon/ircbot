package actions

import "github.com/Zaibon/ircbot"

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

func (h *Help) Do(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
	var output string

	for cmd, _ := range b.HandlersUser {
		output += cmd + ", "
	}
	b.Say(m.Channel(), output[:len(output)-2])
}
