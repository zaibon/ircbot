package ircbot

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestA(t *testing.T) {
}

func TestConnect(t *testing.T) {
	b := NewIrcBot()
	b.Server = "irc.freenode.net"
	b.Port = "6667"
	b.Nick = "ZbiBot"
	b.User = b.Nick

	b.Channel = append(b.Channel, "#testgigx")

	b.Handlers["PONG"] = Pong

	b.Connect()
	fmt.Println("connected")
	b.Listen()
	// fmt.Println("listening")
	// b.HandleActionIn()
	// fmt.Println("handling action")
	// b.HandleActionOut()

	// b.Say("hello world")
	// time.Sleep(5 * time.Second)
	// b.Say("good bye")
	fmt.Println("wait")
	time.Sleep(30 * time.Second)
	b.Disconnect()

	log.SetPrefix("irc> ")
}
