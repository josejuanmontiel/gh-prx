package models

import (
	"regexp"
	"strings"
)

var InvalidTitleCharsMatcher = regexp.MustCompile(`[^.a-zA-Z0-9]`)

type Issue struct {
	Key                 string
	Title               string
	Type                string
	Parent              string // Optional, populated for Agility issues
	SuggestedBranchName string // Optional, populated for Linear issues
}

func (i *Issue) NormalizedTitle() string {
	title := InvalidTitleCharsMatcher.ReplaceAllString(i.Title, "_")
	title = strings.ToLower(title)
	title = strings.Trim(title, "-")

	return title
}
