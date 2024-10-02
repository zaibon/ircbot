package ircbot

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/textproto"
	"os"
	"os/signal"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"

	db "github.com/zaibon/ircbot/database"
)

// IrcBot represents the bot in general
type IrcBot struct {
	// identity
	User string
	Nick string

	// server info
	server   string
	port     uint
	channels []string

	// tcp communication
	conn   net.Conn
	reader *textproto.Reader
	writer *textproto.Writer

	// crypto
	encrypted bool
	config    tls.Config

	// data flow

	//channel to send *IrcMsg to the goroutine that handle input message
	//You usualy don't need to send anything there. It's useful when creating custom actionner
	ChIn chan *IrcMsg
	//channel to send *IrcMsg to the goroutine that handle output message
	//every *IrcMsg send in this channel will be send to the server
	//You usualy don't need to send anything there. It's useful when creating custom actionner
	ChOut chan *IrcMsg
	//channel to send *IrcMsg to the goroutine that handle errors
	ChError chan error

	// exit flag
	Exit chan bool

	//action handlers
	handlersIntern map[string][]Actioner //handler of interanl commands
	handlersUser   map[string]Actioner   // handler of commands fired by user

	//database
	db *db.DB
}

func NewIrcBot(user, nick, server string, port uint, channels []string, DBPath string) *IrcBot {
	bot := IrcBot{
		User:     user,
		Nick:     nick,
		server:   server,
		port:     port,
		channels: channels,

		handlersIntern: make(map[string][]Actioner),
		handlersUser:   make(map[string]Actioner),
		ChIn:           make(chan *IrcMsg),
		ChOut:          make(chan *IrcMsg),
		ChError:        make(chan error),
		Exit:           make(chan bool),
	}

	//init database
	var err error
	bot.db, err = db.Open(DBPath)
	if err != nil {
		log.WithFields(log.Fields{
			"db Path": DBPath,
			"error":   err,
		}).Panicln("unable to open database")
	}

	return &bot
}

/////////////////
//Public function
/////////////////

// Connect connects the bot to the server and joins the channels
func (b *IrcBot) Connect(password string) error {
	//launch a go routine that handle errors
	// b.handleError()

	log.WithField("address", b.url()).Debugln("Connection")

	var conn net.Conn
	var err error

	if b.encrypted {
		cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
		if err != nil {
			log.Errorln("error during certificate loading")
			return err
		}

		config := tls.Config{
			MinVersion:   tls.VersionTLS12,
			Certificates: []tls.Certificate{cert},
			ServerName:   "hello",
		}
		config.Rand = rand.Reader
		conn, err = tls.Dial("tcp", b.url(), &config)
		if err != nil {
			log.WithField("address", b.url()).Errorln("Dial server")
			return err
		}
	} else {
		conn, err = net.Dial("tcp", b.url())
		if err != nil {
			log.WithField("address", b.url()).Errorln("Dial server")
			return err
		}
	}
	b.conn = conn

	r := bufio.NewReader(b.conn)
	w := bufio.NewWriter(b.conn)
	b.reader = textproto.NewReader(r)
	b.writer = textproto.NewWriter(w)

	//connect to server
	b.writer.PrintfLine("USER %s 8 * :%s", b.Nick, b.Nick)
	b.writer.PrintfLine("NICK %s", b.Nick)

	//launch go routines that handle requests
	b.listen()
	b.handleActionIn()
	b.handleActionOut()
	b.handlerError()

	b.identify(password)

	return nil
}

// Disconnect sends QUIT command to server and closes connections
func (b *IrcBot) Disconnect() {
	log.Debugln("disconnection...")
	for _, ch := range b.channels {
		log.WithField("channel", ch).Debugln("leaving")
		b.Say(ch, "See you soon...")
	}
	b.writer.PrintfLine("QUIT")
	b.conn.Close()
	log.Debugln("connection close")
	// b.Exit <- true
}

// Say makes the bot say text to channel
func (b *IrcBot) Say(channel string, text string) {
	msg := NewIrcMsg()
	msg.Command = "PRIVMSG"
	msg.CmdParams = []string{channel}
	msg.Trailing = []string{":", text}

	b.ChOut <- msg
}

// AddInternAction add an action to excecute on internal command (join,connect,...)
// command is the internal command to handle, action is an ActionFunc callback
func (b *IrcBot) AddInternAction(a Actioner) {
	addAction(a, b.handlersIntern)
}

// AddUserAction add an action fired by the user to handle
// command is the commands send by user, action is an ActionFunc callback
func (b *IrcBot) AddUserAction(a Actioner) {
	for _, cmd := range a.Command() {
		b.handlersUser[cmd] = a
	}
}

// GetActionnersCmds returns all registred user actioners commands
func (b *IrcBot) GetActionnersCmds() []string {
	var cmds []string
	for cmd, _ := range b.handlersUser {
		fmt.Println(cmd)
		cmds = append(cmds, cmd)
	}
	return cmds
}

// GetActionUsage returns the Actioner from the user actions map or return an error if
// no action if found with this name
// Usefull if you want to access actioner information within other actioner
// see Help actionner for example
func (b *IrcBot) GetActioner(actionName string) (Actioner, error) {
	actioner, ok := b.handlersUser[actionName]
	if !ok {
		log.WithField("action name", actionName).Warningln("action not found")
		return nil, errors.New("no action found with that name")
	}
	return actioner, nil
}

// DBConnection return a new connection do the database. Use it if your custom action need to access the database
func (b *IrcBot) DBConnection() (*db.DB, error) {
	return db.Open(b.db.Path())
}

// String implements the Stringer interface
func (b *IrcBot) String() string {
	s := fmt.Sprintf("server: %s\n", b.server)
	s += fmt.Sprintf("port: %d\n", b.port)
	s += fmt.Sprintf("ssl: %t\n", b.encrypted)

	if len(b.channels) > 0 {
		s += "channels: "
		for _, v := range b.channels {
			s += fmt.Sprintf("%s ", v)
		}
	}

	return s
}

func (b *IrcBot) url() string {
	return fmt.Sprintf("%s:%d", b.server, b.port)
}

func (b *IrcBot) join() {

	for _, v := range b.channels {
		s := fmt.Sprintf("JOIN %s", v)
		log.WithField("channel", v).Debugln("join")
		b.writer.PrintfLine(s)
	}
}

func (b *IrcBot) identify(password string) {
	//idenify with nickserv
	if password != "" {
		s := fmt.Sprintf("PRIVMSG NickServ :identify %s %s", b.Nick, password)
		log.WithFields(log.Fields{
			"nick":   b.Nick,
			"passwd": "*****",
		}).Debugln("identify")
		b.writer.PrintfLine(s)
	}
}
func (b *IrcBot) listen() {

	go func() {

		for {
			//block read line from socket
			line, err := b.reader.ReadLine()
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Errorln("error reading socket")
				b.ChError <- err
			}

			//convert line into IrcMsg
			msg := ParseLine(line)

			// end of MODT
			if msg.Command == "376" {
				b.join()
			}

			if msg.Command == "PING" {
				out := strings.Replace(line, "PING", "PONG", -1)
				b.writer.PrintfLine(out)
				log.Debugln(out)
			}

			if msg.Command == "PRIVMSG" || msg.Command == "JOIN" {
				b.ChIn <- msg
			}
		}

	}()
}

func addAction(a Actioner, acts map[string][]Actioner) {
	if len(a.Command()) > 1 {
		for _, cmd := range a.Command() {
			acts[cmd] = append(acts[cmd], a)
		}
		return
	}
	acts[a.Command()[0]] = append(acts[a.Command()[0]], a)
}

func (b *IrcBot) handleActionIn() {
	go func() {
		for {
			//receive new message
			msg := <-b.ChIn
			log.WithFields(log.Fields{
				"channel": msg.Channel(),
				"nick":    msg.Nick(),
				"command": msg.Command,
				"params":  strings.Join(msg.CmdParams, " "),
				"message": strings.Join(msg.Trailing, " "),
			}).Debugln("receive")

			//action fired by user
			if len(msg.Trailing) > 0 {
				actionUser, ok := b.handlersUser[msg.Trailing[0]]
				if ok {
					if msg.Channel() == b.Nick {
						//query message, respond to user, not channel
						msg.CmdParams[0] = msg.Nick()
					}
					actionUser.Do(b, msg)
				}
			}

			//action fired by event
			actionsIntern, ok := b.handlersIntern[msg.Command]
			if ok && len(actionsIntern) > 0 {
				for _, action := range actionsIntern {
					action.Do(b, msg)
				}
			}
		}
	}()
}

func (b *IrcBot) handleActionOut() {
	go func() {
		for {
			msg := <-b.ChOut

			params := strings.Join(msg.CmdParams, " ")
			content := strings.Join(msg.Trailing, " ")
			s := fmt.Sprintf("%s %s %s", msg.Command, params, content)
			b.writer.PrintfLine(s)
			log.WithFields(log.Fields{
				"channel": msg.Channel(),
				"nick":    msg.Nick(),
				"command": msg.Command,
				"params":  strings.Join(msg.CmdParams, " "),
				"message": strings.Join(msg.Trailing, " "),
			}).Debugln("send")
		}
	}()
}

func (b *IrcBot) handlerError() {
	go func() {
		for {
			err := <-b.ChError
			log.WithField("error", err).Errorln("error")
			// if err != nil {
			// 	b.Disconnect()
			// 	log.Fatalln("ChError ocurs :", err)
			// }
		}
	}()
}

func (b *IrcBot) signlalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	for {
		select {
		case <-c:
			fmt.Println("disconnection")
			b.Disconnect()
		}
	}
}

func (b *IrcBot) errChk(err error) {
	if err != nil {
		b.ChError <- err
	}
}
