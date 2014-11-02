package actions

import (
	"fmt"

	"github.com/zaibon/ircbot"
)

type Greet struct{}

func (g *Greet) Command() []string {
	return []string{"JOIN"}
}

func (g *Greet) Usage() string {
	return ""
}

func (g *Greet) Do(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
	if m.Nick() == b.Nick {
		return
	}

	b.Say(m.Channel(), fmt.Sprintf("Salut %s", m.Nick()))
}
