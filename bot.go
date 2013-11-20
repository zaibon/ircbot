package ircbot

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"strings"
	"time"
)

type IrcBot struct {
	// identity
	User string
	Nick string

	// server info
	Server  string
	Port    string
	Channel []string

	// tcp communication
	conn   net.Conn
	reader *textproto.Reader
	writer *textproto.Writer

	// data flow
	In    chan *IrcMsg
	Out   chan *IrcMsg
	Error chan error

	// exit flag
	Exit chan bool

	//action handlers
	Handlers map[string]ActionFunc

	//are we joined in channel?
	joined bool
}

func NewIrcBot() *IrcBot {
	return &IrcBot{
		Handlers: make(map[string]ActionFunc),
		In:       make(chan *IrcMsg),
		Out:      make(chan *IrcMsg),
		Error:    make(chan error),
		Exit:     make(chan bool),
		joined:   false,
	}
}

func (b *IrcBot) url() string {
	return fmt.Sprintf("%s:%s", b.Server, b.Port)
}

func (b *IrcBot) Connect() {
	//launch a go routine that handle errors
	// b.handleError()

	log.Println("Info> connection to", b.url())
	tcpCon, err := net.Dial("tcp", b.url())
	if err != nil {
		log.Println("Error> ", err)
		b.Error <- err
	}

	b.conn = tcpCon
	r := bufio.NewReader(b.conn)
	w := bufio.NewWriter(b.conn)
	b.reader = textproto.NewReader(r)
	b.writer = textproto.NewWriter(w)

	b.writer.PrintfLine("USER %s 8 * :%s", b.Nick, b.Nick)
	b.writer.PrintfLine("NICK %s", b.Nick)
}

func (b *IrcBot) Join() {
	for _, v := range b.Channel {
		b.writer.PrintfLine("JOIN %s", v)
	}
	time.Sleep(2 * time.Second)
	b.joined = true
}

func (b *IrcBot) Listen() {

	go func() {

		for {
			//block read line from socket
			line, err := b.reader.ReadLine()
			if err != nil {
				b.Error <- err
			}
			//convert line into IrcMsg
			msg := Parseline(line)
			b.In <- msg
		}

	}()
}

func (b *IrcBot) Say(s string) {
	msg := NewIrcMsg()
	msg.command = "PRIVMSG"
	msg.args = append(msg.args, s)

	b.Out <- msg
}

func (b *IrcBot) HandleActionIn() {
	go func() {
		for {
			//receive new message
			msg := <-b.In
			log.Println(msg.raw)
			//handle action
			action := b.Handlers[msg.command]
			action(b, &msg)
		}
	}()
}

func (b *IrcBot) HandleActionOut() {
	go func() {
		for {
			msg := <-b.Out

			//we send nothing before we sure we join channel
			if b.joined == false {
				continue
			}

			s := fmt.Sprintf("%s %s", msg.command, strings.Join(msg.args, " "))
			fmt.Println("irc >> ", s)
			b.writer.PrintfLine(s)
		}
	}()
}

func (b *IrcBot) HandleError() {
	go func() {
		for {
			err := b.Error
			log.Printf("error > %s", err)
			if err != nil {
				b.Disconnect()
				log.Fatalln("Error ocurs :", err)
			}
		}
	}()
}

func (b *IrcBot) Disconnect() {
	b.writer.PrintfLine("QUIT")
	b.conn.Close()
}
