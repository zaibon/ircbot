package actions

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"code.google.com/p/go.net/html"

	"github.com/Zaibon/ircbot"
)

type TitleExtract struct {
	name string
}

func (u *TitleExtract) Command() []string {
	return []string{
		"PRIVMSG",
	}
}

func (u *TitleExtract) Usage() string {
	return ""
}

func (u *TitleExtract) Do(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
	do(b, m)
}

func do(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
	for _, word := range m.Trailing {

		if !strings.Contains(word, "http") {
			continue
		}

		u, err := url.Parse(word)
		if err != nil {
			fmt.Println("err parse url: ", err)
			continue
		}

		go func() {
			fmt.Println("INFO : start extractTitle,", u.String())
			title, err := extractTitle(u.String())
			if err == nil {
				b.Say(m.Channel(), title)
			}
		}()
	}
}

func extractTitle(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")

	switch {
	case strings.Contains(contentType, "text/html"):
		return parseHTML(resp.Body)
	default:
		return "", fmt.Errorf("mime not supported")
	}
}

func parseHTML(r io.Reader) (string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return "", err
	}
	title := ""
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			title = n.FirstChild.Data
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	if title != "" {
		return title, nil
	}
	return "", fmt.Errorf("no title")
}
