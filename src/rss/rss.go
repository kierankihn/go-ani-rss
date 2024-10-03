package rss

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"go-ani-rss/src/download"
	"go-ani-rss/src/settings"

	"github.com/docker/go-units"
	"github.com/mmcdole/gofeed"
)

type ItemInfo struct {
	Name       string
	Path       string
	Season     int
	Episode    int
	Ext        string
	Url        string
	Size       int64
	ItemConfig *settings.ItemConfig
}

var (
	ItemInfos []ItemInfo

	rssFeedUrl   string
	feedProvider *gofeed.Parser
)

func init() {
	rssFeedUrl = "https://api.ani.rip/"

	feedProvider = gofeed.NewParser()
}

func parserRssFeed() error {
	feed, err := feedProvider.ParseURL(rssFeedUrl)
	if err != nil {
		return err
	}

	const pattern = `\[ANi\] (.*) - (\d{2}) [^\.]*\.(.*)`
	reg, _ := regexp.Compile(pattern)

	ItemInfos = nil

	for _, item := range feed.Items {
		matches := reg.FindStringSubmatch(item.Title)

		if len(matches) == 4 && matches[0] == item.Title {
			for pos, itemConfig := range settings.Config.ItemConfigs {
				itemReg, err := regexp.Compile(itemConfig.Filter)
				if err != nil {
					return err
				}

				if itemReg.MatchString(item.Title) {
					var itemInfo ItemInfo
					itemInfo.Name, itemInfo.Season, itemInfo.Ext, itemInfo.Url, itemInfo.ItemConfig = itemConfig.Name, itemConfig.Season, matches[3], item.Link, &settings.Config.ItemConfigs[pos]

					itemInfo.Episode, _ = strconv.Atoi(matches[2])
					itemInfo.Size, err = units.FromHumanSize(strings.ReplaceAll(item.Extensions["anime"]["size"][0].Value, " ", ""))
					if err != nil {
						return nil
					}

					itemInfo.Path = itemConfig.Path
					itemInfo.Path = strings.ReplaceAll(itemInfo.Path, `{name}`, itemInfo.Name)
					itemInfo.Path = strings.ReplaceAll(itemInfo.Path, `{season}`, fmt.Sprintf("%02d", itemInfo.Season))
					itemInfo.Path = strings.ReplaceAll(itemInfo.Path, `{episode}`, fmt.Sprintf("%02d", itemInfo.Episode))
					itemInfo.Path = strings.ReplaceAll(itemInfo.Path, `{ext}`, itemInfo.Ext)

					ItemInfos = append(ItemInfos, itemInfo)
				}
			}
		}
	}

	return nil
}

func downloadRssItem(itemInfo ItemInfo) error {
	fmt.Printf("Started downloading %s\n", itemInfo.Path)

	// download file
	err := download.DownloadFileFromUrl(itemInfo.Url, itemInfo.Path, itemInfo.Size)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Printf("\nDownloaded %s\n", itemInfo.Path)

	itemInfo.ItemConfig.Progress = itemInfo.Episode

	return nil
}

func ProceedRssItems() error {
	err := parserRssFeed()
	if err != nil {
		return err
	}

	for _, itemInfo := range ItemInfos {
		if itemInfo.ItemConfig.Progress < itemInfo.Episode {
			err := downloadRssItem(itemInfo)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
