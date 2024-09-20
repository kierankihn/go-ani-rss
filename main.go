package main

import (
	"fmt"
	"time"

	"go-ani-rss/src/rss"
	"go-ani-rss/src/settings"
)

var ()

func main() {
	err := settings.ParserSettings()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		err = rss.ProceedRssItems()
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = settings.SaveSettings()
		if err != nil {
			fmt.Println(err)
			continue
		}

		time.Sleep(60 * time.Second)
	}
}
