package actions

import "github.com/Zaibon/ircbot"

type Ping struct{}

func (p *Ping) Command() []string {
	return []string{".ping"}
}
func (p *Ping) Usage() string {
	return ".ping : send ping request"
}

func (p *Ping) Do(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
	b.Say(m.Channel(), "pong")
}
