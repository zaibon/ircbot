package ircbot

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"
	// "time"
)

type requestServer struct {
	http.Server
}

func handler(b *IrcBot, rw http.ResponseWriter, req *http.Request) {

	bob, err := ioutil.ReadAll(req.Body)
	if err != nil {
		b.ChError <- err
	}

	ircReq, err := DecodeIrcReq(bob)
	if err != nil {
		b.ChError <- err
		time.Sleep(1 * time.Second)
		panic("error decoding json")
	}
	fmt.Printf("WEB << %+v", ircReq)

	b.ChOut <- &IrcMsg{
		Command: ircReq.Command,
		Channel: ircReq.Channel,
		Args:    ircReq.Args,
	}

	rw.Write([]byte("send"))
}

func send(b *IrcBot, rw http.ResponseWriter, req *http.Request) {

	req.ParseForm()
	channel := req.PostFormValue("channel")
	text := req.PostFormValue("text")

	b.ChOut <- &IrcMsg{
		Command: "PRIVMSG",
		Channel: channel,
		Args: []string{
			":",
			text,
		},
	}
	http.Redirect(rw, req, "/qg", http.StatusFound)
}

var tmpl string = `<html>
	<head>
		<title>ZbiBot Centre QG</title>
	</head>
	<body>
		<div id="content">
			<form action="/send" method="post">
				<div><input type="text" name="channel"/></div>
				<div><input type="text" name="text"/></div>
				<div><input type="submit" value="send"></div>
			</form>
		</div>
	</body>
</html>`

func gui(rw http.ResponseWriter, req *http.Request) {
	t := template.New("form")
	t.Parse(tmpl)
	t.ExecuteTemplate(rw, "form", nil)
}
