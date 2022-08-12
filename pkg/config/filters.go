package config

import (
	"regexp"
	"strings"
)

type Filter interface {
	GetRegex() []Regex
	GetRemove() []string
	GetReplace() map[string]string
}

func ApplyFilters(filter Filter, value string) string {
	for _, pattern := range filter.GetRegex() {
		re := regexp.MustCompile(pattern.Pattern)
		match := re.FindStringSubmatch(value)
		if len(match) > 0 && pattern.Group >= 0 && pattern.Group < len(match) {
			value = match[pattern.Group]
		}
	}
	for _, remove := range filter.GetRemove() {
		value = strings.ReplaceAll(value, remove, "")
	}

	for k, v := range filter.GetReplace() {
		value = strings.ReplaceAll(value, k, v)
	}
	return value
}
