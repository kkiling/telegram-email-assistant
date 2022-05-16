package common

import (
	"regexp"
)

func GetContentId(value string) string {
	if value == "" {
		return ""
	}
	var re = regexp.MustCompile(`(?m)\<(.*?)\>`)
	result := re.FindAllStringSubmatch(value, -1)
	return result[0][1]
}
