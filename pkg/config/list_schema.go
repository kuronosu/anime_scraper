package config

import (
	"github.com/gocolly/colly/v2"
)

func NewParsedLinks() ParsedLinks {
	return ParsedLinks{
		Links: make([]string, 0),
		Data:  make(map[string]interface{}),
	}
}

type ParsedLinks struct {
	Links []string               `json:"links"`
	Data  map[string]interface{} `json:"data"`
}

func (pl *ParsedLinks) _check() {
	if pl.Links == nil {
		pl.Links = make([]string, 0)
	}
	if pl.Data == nil {
		pl.Data = make(map[string]interface{})
	}
}

func (pl *ParsedLinks) Extend(links ParsedLinks) {
	pl._check()
	links._check()
	for _, l := range links.Links {
		pl.Links = append(pl.Links, l)
		if _, ok := pl.Data[l]; !ok {
			pl.Links = append(pl.Links, l)
			pl.Data[l] = links.Data[l]
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

func (ls *LinkField) SafeCompile(e *colly.HTMLElement) string {
	if ls.Selector == "" {
		if ls.Attr == "" {
			return e.Attr("href")
		} else if ls.Attr == "text" {
			return e.Text
		} else {
			return e.Attr(ls.Attr)
		}
	} else {
		if ls.Attr == "" {
			return e.ChildAttr(ls.Selector, "href")
		} else if ls.Attr == "text" {
			return e.ChildText(ls.Selector)
		} else {
			return e.ChildAttr(ls.Selector, ls.Attr)
		}
	}
}

type ListSchema struct {
	ContainerSelector string `yaml:"container_selector"`
	ItemSelector      string `yaml:"item_selector"`
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
	_links := make([]string, 0)
	_data := make(map[string]interface{})
	for _, d := range data {
		_links = append(_links, d.A)
		_data[d.A] = d.Data
	}
	return ParsedLinks{
		Links: _links,
		Data:  _data,
	}
}
