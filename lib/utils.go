package lib

import (
	"regexp"
	"strings"
)

func ToSnakeCase(s string) string {

	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(s, "${1}_${2}")

	re = regexp.MustCompile("([A-Z])([A-Z][a-z])")
	snake = re.ReplaceAllString(snake, "${1}_${2}")

	return strings.ToLower(snake)
}
