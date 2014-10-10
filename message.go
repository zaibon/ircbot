package ircbot

import "strings"

//IrcMsg represent a message receive or send to the irc server
type IrcMsg struct {
	Raw    string
	Prefix string

	Command   string
	CmdParams []string

	Trailing []string
}

func NewIrcMsg() *IrcMsg {
	return &IrcMsg{}
}

func (m *IrcMsg) parseline(line string) {
	prefixEnd, trailingStart := -1, len(line)
	m.Prefix, m.Command = "", ""

	if strings.HasPrefix(line, ":") {
		prefixEnd = strings.Index(line, " ")
		m.Prefix = line[1:prefixEnd]
	}

	trailingStart = strings.Index(line, " :")
	if trailingStart >= 0 {
		m.Trailing = strings.Fields(line[trailingStart+2:])
	} else {
		trailingStart = len(line) - 1
	}

	cmdAndParams := strings.Fields(line[(prefixEnd + 1) : trailingStart+1])
	if len(cmdAndParams) > 0 {
		m.Command = cmdAndParams[0]
	}
	if len(cmdAndParams) > 1 {
		m.CmdParams = cmdAndParams[1:]
	}
}

//Channel return the channel that send the IrcMsg
func (m *IrcMsg) Channel() string {
	return m.CmdParams[0]
}

//Channel return the nickname of the user who send IrcMsg
func (m *IrcMsg) Nick() string {
	if strings.Contains(m.Prefix, "!") {
		tmp := strings.SplitAfterN(m.Prefix, "!", 2)[0]
		return tmp[:len(tmp)-1]
	}
	return ""
}

//ParseLine parse a line receive from server and return a new IrcMsg object
func ParseLine(line string) *IrcMsg {
	msg := NewIrcMsg()
	msg.parseline(line)
	return msg
}
