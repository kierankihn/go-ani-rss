package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/docker/go-units"
	"github.com/mmcdole/gofeed"
	"github.com/schollz/progressbar/v3"

	"go-ani-rss/src/format"
	"go-ani-rss/src/settings"
)

func downloadFileFromUrl(url string, path string, size int64) error {
	// create output file
	outputDir := filepath.Dir(path)
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}
	output, err := os.Create(path)
	if err != nil {
		return err
	}
	defer output.Close()

	// make http request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// progress bar
	bar := progressbar.DefaultBytes(size, "downloading")

	// write into target file
	io.Copy(io.MultiWriter(output, bar), response.Body)

	return nil
}

func proceedRssItem(item *gofeed.Item, itemConfig settings.ItemConfig) error {
	fmt.Printf("Started downloading %s\n", item.Title)

	// convert file size
	size, err := units.FromHumanSize(strings.Replace(item.Extensions["anime"]["size"][0].Value, " ", "", -1))
	if err != nil {
		return nil
	}

	// convert download path
	path, err := format.FormatPath(item.Title, itemConfig.Name, itemConfig.Path)
	if err != nil {
		return err
	}
	// download file
	err = downloadFileFromUrl(item.Link, path, size)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Printf("\nDownloaded %s\n", item.Title)
	return nil
}

func ProceedRssItems(feed *gofeed.Feed, ItemConfigs []settings.ItemConfig) error {
	for _, item := range feed.Items {
		for _, itemConfig := range ItemConfigs {
			reg, err := regexp.Compile(itemConfig.Filter)
			if err != nil {
				return err
			}

			if reg.MatchString(item.Title) {
				if item.PublishedParsed.After(settings.Config.LastDownloadTime) {
					err = proceedRssItem(item, itemConfig)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	err := settings.SaveSettings()
	if err != nil {
		return err
	}

	return nil
}
