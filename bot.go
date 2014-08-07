package ircbot

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"
	"time"
)

type IrcBot struct {
	// identity
	User     string
	Nick     string
	password string

	// server info
	server   string
	port     string
	channels []string

	// tcp communication
	conn   net.Conn
	reader *textproto.Reader
	writer *textproto.Writer

	// web interface
	webEnable bool
	webPort   string

	// crypto
	encrypted bool
	config    tls.Config

	// data flow
	ChIn    chan *IrcMsg
	ChOut   chan *IrcMsg
	ChError chan error

	// exit flag
	Exit chan bool

	//action handlers
	handlersIntern map[string][]Actioner //handler of interanl commands
	handlersUser   map[string]Actioner   // handler of commands fired by user

	//are we joined ChIn channel?
	joined bool

	db *DB
}

func NewIrcBot(user, nick, password, server, port string, channels []string) *IrcBot {
	bot := IrcBot{
		User:     user,
		Nick:     nick,
		password: password,
		server:   server,
		port:     port,
		channels: channels,

		handlersIntern: make(map[string][]Actioner),
		handlersUser:   make(map[string]Actioner),
		ChIn:           make(chan *IrcMsg),
		ChOut:          make(chan *IrcMsg),
		ChError:        make(chan error),
		Exit:           make(chan bool),
		joined:         false,
	}

	//defautl actions, needed to run proprely
	bot.AddInternAction(&pong{})
	bot.AddInternAction(&validConnect{})
	bot.AddInternAction(&Help{})

	//init database
	if err := bot.db.open("irc.db"); err != nil {
		panic(err)
	}
	if err := bot.db.init(); err != nil {
		panic(err)
	}

	return &bot
}

/////////////////
//Public function
/////////////////

//Connect connects the bot to the server and joins the channels
func (b *IrcBot) Connect() {
	//launch a go routine that handle errors
	// b.handleError()

	log.Println("Info> connection to", b.url())

	var tcpCon net.Conn
	var err error
	if b.encrypted {
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

	//connect to server
	b.writer.PrintfLine("USER %s 8 * :%s", b.Nick, b.Nick)
	b.writer.PrintfLine("NICK %s", b.Nick)

	//launch go routines that handle requests
	b.listen()
	b.handleActionIn()
	b.handleActionOut()
	b.handlerError()
	if b.webEnable {
		b.handleWeb()
	}

	//join all channels
	b.join()

	b.identify()
}

//Disconnect sends QUIT command to server and closes connections
func (b *IrcBot) Disconnect() {
	b.writer.PrintfLine("QUIT")
	b.conn.Close()
	b.Exit <- true
}

func (b *IrcBot) IsJoined() bool {
	return b.joined
}

//Say makes the bot say text to channel
func (b *IrcBot) Say(channel string, text string) {
	msg := NewIrcMsg()
	msg.Command = "PRIVMSG"
	msg.Channel = channel
	msg.Args = append(msg.Args, ":"+text)

	b.ChOut <- msg
}

//AddInternAction add an action to excutre on internal command (join,connect,...)
//command is the internal command to handle, action is an ActionFunc callback
func (b *IrcBot) AddInternAction(a Actioner) {
	addAction(a, b.handlersIntern)
}

//AddUserAction add an action fired by the user to handle
//command is the commands send by user, action is an ActionFunc callback
func (b *IrcBot) AddUserAction(a Actioner) {
	for _, cmd := range a.Command() {
		b.handlersUser[cmd] = a
	}
}

//String implements the Stringer interface
func (b *IrcBot) String() string {
	s := fmt.Sprintf("server: %s\n", b.server)
	s += fmt.Sprintf("port: %s\n", b.port)
	s += fmt.Sprintf("ssl: %t\n", b.encrypted)

	if len(b.channels) > 0 {
		s += "channels: "
		for _, v := range b.channels {
			s += fmt.Sprintf("%s ", v)
		}
	}

	return s
}

//EnableWeb enables the web interface of the bot
func (b *IrcBot) EnableWeb() {
	b.webEnable = true
}

//SetWebPort sets the port on wich the web interface will listen
func (b *IrcBot) SetWebPort(port int) {
	b.webPort = strconv.Itoa(port)
}

func (b *IrcBot) url() string {
	return fmt.Sprintf("%s:%s", b.server, b.port)
}

func (b *IrcBot) join() {

	//prevent to send JOIN command before we are conected
	for {
		if !b.joined {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	for _, v := range b.channels {
		s := fmt.Sprintf("JOIN %s", v)
		fmt.Println("irc >> ", s)
		b.writer.PrintfLine(s)
	}
	b.joined = true
}

func (b *IrcBot) identify() {
	//idenify with nickserv
	if b.password != "" {
		s := fmt.Sprintf("PRIVMSG NickServ :identify %s %s", b.Nick, b.password)
		fmt.Println("irc >> ", s)
		b.writer.PrintfLine(s)
	}
}

func (b *IrcBot) listen() {

	go func() {

		for {
			//block read line from socket
			line, err := b.reader.ReadLine()
			if err != nil {
				b.ChError <- err
			}
			//convert line into IrcMsg
			msg := ParseLine(line)
			b.ChIn <- msg
			if err := logMsg(msg, b.db); err != nil {
				b.ChError <- err
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
			fmt.Println("irc << ", msg.Raw)

			if msg.Command == "JOIN" && msg.Nick == b.Nick {
				b.joined = true
			}

			if msg.Command == "PRIVMSG" && strings.HasPrefix(msg.Args[0], ":.") {
				action, ok := b.handlersUser[strings.TrimPrefix(msg.Args[0], ":")]
				if ok {
					action.Do(b, msg)
				}
			} else {
				actions, ok := b.handlersIntern[msg.Command]
				//handle action
				if ok && len(actions) > 0 {
					for _, action := range actions {
						action.Do(b, msg)
					}
				}
			}
		}
	}()
}

func (b *IrcBot) handleActionOut() {
	go func() {
		for {
			msg := <-b.ChOut

			//we send nothing before we sure we join channel
			if b.joined == false {
				continue
			}

			s := fmt.Sprintf("%s %s %s", msg.Command, msg.Channel, strings.Join(msg.Args, " "))
			fmt.Println("irc >> ", s)
			b.writer.PrintfLine(s)
		}
	}()
}

func (b *IrcBot) handlerError() {
	go func() {
		for {
			err := <-b.ChError
			fmt.Printf("error >> %s", err)
			if err != nil {
				b.Disconnect()
				log.Fatalln("ChError ocurs :", err)
			}
		}
	}()
}

//handleWeb handles requests receive on http server
func (b *IrcBot) handleWeb() {
	go func() {
		http.HandleFunc("/qg", gui)
		http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
			send(b, w, r)
		})
		http.HandleFunc("/ircbot", func(w http.ResponseWriter, r *http.Request) {
			handler(b, w, r)
		})
		http.ListenAndServe(":"+b.webPort, nil)
	}()
}

func (b *IrcBot) errChk(err error) {
	if err != nil {
		log.Println("ChError> ", err)
		b.ChError <- err
	}
}
