package actions

import (
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"

	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/andybalholm/cascadia"

	"github.com/zaibon/ircbot"
)

type TitleExtract struct {
	selector cascadia.Selector
}

func NewTitleExtract() *TitleExtract {
	return &TitleExtract{
		selector: cascadia.MustCompile("title"),
	}
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
	u.do(b, m)
}

func (u *TitleExtract) do(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
	for _, word := range m.Trailing {

		if !strings.Contains(word, "http") {
			continue
		}

		URL, err := url.Parse(word)
		if err != nil {
			log.WithField("url", word).Errorln("error parsing url")
			continue
		}

		go func() {
			log.WithField("url", URL.String()).Debugln("Start extract title")

			title, err := extractTitle(URL.String(), u.selector)
			if err == nil {
				b.Say(m.Channel(), title)
			}

			log.WithFields(log.Fields{
				"url":   URL.String(),
				"title": title,
			}).Debugln("title found")
		}()
	}
}

func extractTitle(url string, selector cascadia.Selector) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")

	switch {
	case strings.Contains(contentType, "text/html"):
		return cssSelectHTML(resp.Body, selector)
	default:
		return "", fmt.Errorf("mime not supported")
	}

}

func cssSelectHTML(r io.Reader, selector cascadia.Selector) (string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", err
	}
	title := ""
	doc.Find("title").Each(func(i int, s *goquery.Selection) {
		title = title + s.Text()
	})
	return title, nil
}
