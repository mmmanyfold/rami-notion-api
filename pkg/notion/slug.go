package notion

import (
	"regexp"
	"strings"
)

var re = regexp.MustCompile("[^a-z0-9]+")

func Slug(str string) string {
	return strings.Trim(re.ReplaceAllString(strings.ToLower(str), "-"), "-")
}
