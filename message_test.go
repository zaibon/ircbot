package ircbot

import (
	"testing"
)

var parseLineTestTable = []struct {
	input  string
	expect *IrcMsg
}{
	{
		":weber.freenode.net NOTICE * :*** Looking up your hostname... ",
		&IrcMsg{
			Prefix:    "weber.freenode.net",
			Command:   "NOTICE",
			CmdParams: []string{"*"},
			Trailing:  []string{"***", "Looking", "up", "your", "hostname..."},
		},
	},

	{
		":zbitest MODE zbitest :+i ",
		&IrcMsg{
			Prefix:    "zbitest",
			Command:   "MODE",
			CmdParams: []string{"zbitest"},
			Trailing:  []string{"+i"},
		},
	},

	{
		":Zaibon!~zaibon@142.ip-37-187-37.eu PRIVMSG #zbitest :.h",
		&IrcMsg{
			Prefix:    "Zaibon!~zaibon@142.ip-37-187-37.eu",
			Command:   "PRIVMSG",
			CmdParams: []string{"#zbitest"},
			Trailing:  []string{".h"},
		},
	},
	{
		":zbiii!~zbiii@17.123-247-81.adsl-dyn.isp.belgacom.be JOIN #zbitest",
		&IrcMsg{
			Prefix:    "zbiii!~zbiii@17.123-247-81.adsl-dyn.isp.belgacom.be",
			Command:   "JOIN",
			CmdParams: []string{"#zbitest"},
			Trailing:  []string{""},
		},
	},
}

func TestParseLine(t *testing.T) {
	for _, tt := range parseLineTestTable {
		actual := &IrcMsg{}
		actual.parseline(tt.input)
		if !equal(actual, tt.expect) {
			t.Errorf("\ninput %s\nexpected %+v\nactual %+v\n", tt.input, tt.expect, actual)
		}
	}
}

var nickTestTable = []struct {
	input  string
	expect string
}{
	{
		":weber.freenode.net NOTICE * :*** Looking up your hostname... ",
		"",
	},
	{
		":zbitest MODE zbitest :+i ",
		"",
	},
	{
		":Zaibon!~zaibon@142.ip-37-187-37.eu PRIVMSG #zbitest :.h",
		"Zaibon",
	},
}

func TestNick(t *testing.T) {
	for _, tt := range nickTestTable {
		actual := &IrcMsg{}
		actual.parseline(tt.input)
		if actual.Nick() != tt.expect {
			t.Errorf("\ninput %s\nexpected %+v\nactual %+v\n", tt.input, tt.expect, actual.Nick())
		}
	}
}

var channelTestTable = []struct {
	input  string
	expect string
}{
	{
		":weber.freenode.net NOTICE * :*** Looking up your hostname... ",
		"*",
	},
	{
		":zbitest MODE zbitest :+i ",
		"zbitest",
	},
	{
		":Zaibon!~zaibon@142.ip-37-187-37.eu PRIVMSG #zbitest :.h",
		"#zbitest",
	},
}

func TestChannel(t *testing.T) {
	for _, tt := range channelTestTable {
		actual := &IrcMsg{}
		actual.parseline(tt.input)
		if actual.Channel() != tt.expect {
			t.Errorf("\ninput %s\nexpected %+v\nactual %+v\n", tt.input, tt.expect, actual.Channel())
		}
	}
}

func equal(m1, m2 *IrcMsg) bool {
	if m1.Prefix != m2.Prefix {
		return false
	}

	if m1.Command != m2.Command {
		return false
	}

	for i := range m1.CmdParams {
		if m1.CmdParams[i] != m2.CmdParams[i] {
			return false
		}
	}

	for i := range m1.Trailing {
		if m1.Trailing[i] != m2.Trailing[i] {
			return false
		}
	}
	return true
}
