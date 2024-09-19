package format

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"go-ani-rss/src/settings"
)

func FormatPath(originName string, itemConfig settings.ItemConfig) (string, error) {
	// compile regexp
	const pattern = `\[ANi\] (.*) - (\d{2}) [^\.]*\.(.*)`
	reg, _ := regexp.Compile(pattern)

	matches := reg.FindStringSubmatch(originName)

	newName := itemConfig.Path

	if len(matches) == 4 && matches[0] == originName {
		name, seasonId, episodeId, extName := itemConfig.Name, strconv.Itoa(itemConfig.Season), matches[2], matches[3]

		newName = strings.ReplaceAll(newName, `{name}`, name)
		newName = strings.ReplaceAll(newName, `{season}`, seasonId)
		newName = strings.ReplaceAll(newName, `{episode}`, episodeId)
		newName = strings.ReplaceAll(newName, `{ext}`, extName)

		return newName, nil
	}
	return "", errors.New("origin name is not valid")
}
