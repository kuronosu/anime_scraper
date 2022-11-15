package config

import (
	"github.com/gocolly/colly/v2"
)

type ParsedLinks map[string]map[string]interface{}

func (pl ParsedLinks) Extend(links ParsedLinks) {
	if links == nil {
		return
	}
	for key, element := range links {
		if _, ok := pl[key]; !ok {
			pl[key] = element
		}
	}
}

type LinkData struct {
	A    string                 `json:"a"`
	Data map[string]interface{} `json:"data"`
}

type LinkField struct {
	Selector string `yaml:"selector"`
	Attr     string
	Regex    []Regex
	Remove   []string
	Replace  map[string]string
}

func (field LinkField) GetRegex() []Regex {
	return field.Regex
}

func (field LinkField) GetRemove() []string {
	return field.Remove
}

func (field LinkField) GetReplace() map[string]string {
	return field.Replace
}

func (ls *LinkField) SafeCompile(e *colly.HTMLElement) string {
	if ls.Selector == "" {
		if ls.Attr == "" {
			return ApplyFilters(ls, e.Attr("href"))
		} else if ls.Attr == "text" {
			return ApplyFilters(ls, e.Text)
		} else {
			return ApplyFilters(ls, e.Attr(ls.Attr))
		}
	} else {
		if ls.Attr == "" {
			return ApplyFilters(ls, e.ChildAttr(ls.Selector, "href"))
		} else if ls.Attr == "text" {
			return ApplyFilters(ls, e.ChildText(ls.Selector))
		} else {
			return ApplyFilters(ls, e.ChildAttr(ls.Selector, ls.Attr))
		}
	}
}

type ListSchema struct {
	ContainerSelector string `yaml:"container_selector"`
	ItemSelector      string `yaml:"item_selector"`
	Prefix            string `yaml:"prefix"`
	IncludePrefix     bool   `yaml:"include_prefix"`
	Pagination        Pagination
	Link              struct {
		A    LinkField
		Data []Field
	}
}

func (ls *ListSchema) SafeCompile(e *colly.HTMLElement) []LinkData {
	_items := make([]LinkData, 0)
	e.ForEach(ls.ItemSelector, func(_ int, ite *colly.HTMLElement) {
		link := LinkData{}
		link.A = ls.Link.A.SafeCompile(ite)
		if ls.IncludePrefix {
			link.A = ls.Prefix + link.A
		}

		link.Data = make(map[string]interface{})
		for _, f := range ls.Link.Data {
			if f.Selector != "" {
				ite.ForEachWithBreak(f.Selector, func(_ int, iite *colly.HTMLElement) bool {
					link.Data[f.Name] = f.SafeCompile(iite)
					return false
				})
			} else {
				link.Data[f.Name] = f.SafeCompile(ite)
			}
		}
		_items = append(_items, link)
	})
	return _items
}

func ParseLinkData(data []LinkData) ParsedLinks {
	pl := make(ParsedLinks)
	for _, d := range data {
		pl[d.A] = d.Data
	}
	return pl
}

func ExtractLinks(data []LinkData) []string {
	links := make([]string, 0)
	for _, d := range data {
		links = append(links, d.A)
	}
	return links
}
