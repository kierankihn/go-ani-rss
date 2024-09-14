package main

import (
	"fmt"
	"go-ani-rss/src/download"
	"go-ani-rss/src/settings"
	"time"

	"github.com/mmcdole/gofeed"
)

var (
	rssFeedUrl   string
	feedProvider *gofeed.Parser
)

func init() {
	rssFeedUrl = "https://api.ani.rip/ani-download.xml"

	feedProvider = gofeed.NewParser()
}

func main() {
	err := settings.ParserSettings()
	if err != nil {
		fmt.Println(err)
	}

	for {
		feed, err := feedProvider.ParseURL(rssFeedUrl)
		if err != nil {
			fmt.Println(err)
		}

		err = download.ProceedRssItems(feed)
		if err != nil {
			fmt.Println(err)
		}

		time.Sleep(60 * time.Second)
	}
}
