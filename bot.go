package ircbot

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
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

	// crypto
	Encrypted bool
	config    tls.Config

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

func (b *IrcBot) loadCert() {

}

func (b *IrcBot) Connect() {
	//launch a go routine that handle errors
	// b.handleError()

	log.Println("Info> connection to", b.url())

	var tcpCon net.Conn
	var err error
	if b.Encrypted {
		cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
		b.errChk(err)

		config := tls.Config{Certificates: []tls.Certificate{cert}}
		config.Rand = rand.Reader
		tcpCon, err = tls.Dial("tcp", b.url(), &config)
		b.errChk(err)

	} else {
		tcpCon, err = net.Dial("tcp", b.url())
		b.errChk(err)
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

	//prevent to send JOIN command before we are conected
	for {
		if !b.joined {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	for _, v := range b.Channel {
		s := fmt.Sprintf("JOIN %s", v)
		fmt.Println("irc >> ", s)
		b.writer.PrintfLine(s)
	}
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
			fmt.Println("irc << ", msg.raw)
			//handle action
			if action := b.Handlers[msg.command]; action != nil {
				action(b, msg)
			}
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
			err := <-b.Error
			fmt.Printf("error > %s", err)
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
