package actions

import (
	"fmt"
	"strings"

	"github.com/zaibon/ircbot"
)

type Help struct{}

func (h *Help) Command() []string {
	return []string{
		".help",
		".h",
	}
}

func (h *Help) Usage() string {
	return ".h <command>"
}

func (h *Help) Do(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
	if len(m.Trailing) < 2 {
		var output string
		fmt.Printf("%+v", b.GetActionnersCmds())
		for _, cmd := range b.GetActionnersCmds() {
			output += fmt.Sprintf("%s, ", cmd)
		}
		b.Say(m.Channel(), output[:len(output)-2])
		return
	}

	cmd := m.Trailing[1]
	if !strings.HasPrefix(cmd, ".") {
		cmd = fmt.Sprintf(".%s", cmd)
	}
	actioner, err := b.GetActioner(cmd)
	if err != nil {
		b.Say(m.Channel(), err.Error())
		return
	}

	lines := strings.Split(actioner.Usage(), "\n")
	for _, line := range lines {
		b.Say(m.Channel(), line)
	}
}
