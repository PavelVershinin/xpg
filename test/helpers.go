package test

import (
	"regexp"
	"strings"
)

func ClearQuery(s string) string {
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = regexp.MustCompile(`[\s]{2,}`).ReplaceAllString(s, " ")
	s = regexp.MustCompile(`(["|'])\s,`).ReplaceAllString(s, "$1,")
	s = regexp.MustCompile(`,\s(["|'])`).ReplaceAllString(s, ",$1")
	return strings.TrimSpace(s)
}
