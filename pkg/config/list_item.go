package config

import "github.com/gocolly/colly/v2"

type ListItem struct {
	Plain    bool
	Selector *string
	Field    Field
	Fields   []Field
}

func (li *ListItem) SafeCompile(e *colly.HTMLElement) interface{} {
	if li.Plain {
		_items := make([]interface{}, 0)
		s := li.Field.Selector
		if li.Selector != nil && *li.Selector != "" {
			s = *li.Selector
		}
		e.ForEach(s, func(_ int, ite *colly.HTMLElement) {
			_items = append(_items, li.Field.SafeCompileFlat(ite))
		})
		return _items
	}
	_items := make([]interface{}, 0)
	selector := li.Selector
	if selector != nil && *selector != "" {
		e.ForEach(*selector, func(_ int, ite *colly.HTMLElement) {
			_map := make(map[string]interface{})
			for _, f := range li.Fields {
				if f.Selector != "" {
					ite.ForEachWithBreak(f.Selector, func(_ int, iite *colly.HTMLElement) bool {
						_map[f.Name] = f.SafeCompile(iite)
						return false
					})
				} else {
					_map[f.Name] = f.SafeCompile(ite)
				}
			}
			_items = append(_items, _map)
		})
	}

	return _items
}
