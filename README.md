ircbot
======

Simple irc bot package in Go

Example of implementation can be found at ttps://github.com/zaibon/zbibot

##installation
````bash
go get github.com/zaibon/ircbot
````

##usage
````go
import (
	"github.com/zaibon/ircbot"
	"github.com/zaibon/ircbot/actions"
)


func main(){
	//create new bot
	channels := string[]{
		"go-nuts",
	}
	b := ircbot.NewIrcBot("ircbot", "ircbot", "password", "irc.freenode.net", "6667", channels, "irc.db")

	//add custom intern actions
	b.AddInternAction(&actions.Greet{})
	b.AddInternAction(&actions.TitleExtract{})
	b.AddInternAction(actions.NewLogger(b))
	b.AddInternAction(actions.NewURLLog(b))

	//add command fire by users
	b.AddUserAction(&actions.Help{})

	//connectin to server, listen and serve
	b.Connect()

	//block until we send something to b.Exit channel
	<-b.Exit
}

b.Disconnect()
````