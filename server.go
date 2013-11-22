package ircbot

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	// "time"
)

type requestServer struct {
	http.Server
}

func Handler(b *IrcBot, rw http.ResponseWriter, req *http.Request) {

	bob, err := ioutil.ReadAll(req.Body)
	if err != nil {
		b.Error <- err
	}

	ircReq, err := DecodeIrcReq(bob)
	if err != nil {
		b.Error <- err
		time.Sleep(1 * time.Second)
		panic("error decoding json")
	}
	fmt.Printf("WEB << %+v", ircReq)

	b.Out <- &IrcMsg{
		Command: "PRIVMSG",
		Channel: ircReq.Channel,
		Args:    ircReq.Args,
	}

	rw.Write([]byte("send"))
}
