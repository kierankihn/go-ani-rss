package format

import (
	"errors"
	"regexp"
	"strings"
)

func getSeasonId(name string) string {
	chineseToDigit := map[string]string{"零": "00", "一": "01", "二": "02", "三": "03", "四": "04", "五": "05", "六": "06", "七": "07", "八": "08", "九": "09", "十": "10"}

	const pattern = `第.季`
	reg, _ := regexp.Compile(pattern)

	matches := reg.FindStringSubmatch(name)

	if len(matches) == 2 && matches[0] == name {
		return chineseToDigit[matches[1]]
	}

	return "01"
}

func FormatPath(originName string, title string, format string) (string, error) {
	// compile regexp
	const pattern = `\[ANi\] (.*) - (\d{2}) [^\.]*\.(.*)`
	reg, _ := regexp.Compile(pattern)

	matches := reg.FindStringSubmatch(originName)

	if len(matches) == 4 && matches[0] == originName {
		name, seasonId, episodeId, extName := title, getSeasonId(matches[1]), matches[2], matches[3]

		format = strings.ReplaceAll(format, `{name}`, name)
		format = strings.ReplaceAll(format, `{season}`, seasonId)
		format = strings.ReplaceAll(format, `{episode}`, episodeId)
		format = strings.ReplaceAll(format, `{ext}`, extName)

		return format, nil
	}
	return "", errors.New("origin name is not valid")
}
