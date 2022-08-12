package config

import "github.com/gocolly/colly/v2"

type PaginationLink struct {
	Selector string
	Attr     string
	Prefix   string
	Limit    int
}

func (pl *PaginationLink) GetValue(e *colly.HTMLElement) string {
	if pl.Attr == "" {
		return pl.Prefix + e.Attr("href")
	} else if pl.Attr == "text" {
		return pl.Prefix + e.Text
	} else {
		return pl.Prefix + e.Attr(pl.Attr)
	}
}

type Pagination struct {
	Next     PaginationLink
	Previous PaginationLink
}
