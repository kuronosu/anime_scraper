package config

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"gopkg.in/yaml.v2"
)

type Regex struct {
	Pattern string
	Group   int
}

type ListItem struct {
	Plain    bool
	Selector *string
	Field    Field
	Fields   []Field
}

type Field struct {
	Name     string
	Selector string
	Type     *string
	Attr     *string
	Replace  map[string]string
	Remove   []string
	Regex    []Regex
	Item     *ListItem
}

func (field Field) SafeCompile(e *colly.HTMLElement) interface{} {
	if field.IsFlat() {
		return field.SafeCompileFlat(e)
	}
	return field.SafeCompileDeep(e)
}

func (field Field) SafeCompileDeep(e *colly.HTMLElement) interface{} {
	if field.IsList() {
		return field.Item.SafeCompile(e)
	}
	return nil
}

func (field Field) SafeCompileFlat(e *colly.HTMLElement) interface{} {
	value, raw, err := field.CompileFlat(e)
	if err != nil {
		value = raw
	}
	return value
}

func (field Field) CompileFlat(e *colly.HTMLElement) (interface{}, string, error) {
	var rawValue string
	if field.Attr != nil {
		rawValue = strings.TrimSpace(e.Attr(*field.Attr))
	} else {
		rawValue = strings.TrimSpace(e.Text)
	}

	if !field.IsFlat() {
		return nil, rawValue, fmt.Errorf("field %s is not flat", field.Name)
	}

	value := rawValue

	for _, pattern := range field.Regex {
		re := regexp.MustCompile(pattern.Pattern)
		match := re.FindStringSubmatch(value)
		if len(match) > 0 && pattern.Group >= 0 && pattern.Group < len(match) {
			value = match[pattern.Group]
		}
	}
	for _, remove := range field.Remove {
		value = strings.ReplaceAll(value, remove, "")
	}

	for k, v := range field.Replace {
		value = strings.ReplaceAll(value, k, v)
	}

	if field.IsInt() {
		v, err := strconv.Atoi(value)
		return v, rawValue, err
	}
	if field.IsFloat() {
		v, err := strconv.ParseFloat(value, 64)
		return v, rawValue, err
	}
	if field.IsString() {
		return value, rawValue, nil
	}

	return value, rawValue, fmt.Errorf("unsupported type %s", field._type())
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

func (field Field) _type() string {
	if field.Type == nil {
		return "string"
	}
	return *field.Type
}

func (field Field) IsFlat() bool {
	return field.IsInt() || field.IsFloat() || field.IsString()
}

func (field Field) IsInt() bool {
	return strings.ToLower(field._type()) == "int"
}

func (field Field) IsFloat() bool {
	_type := strings.ToLower(field._type())
	return _type == "float" || _type == "float32" || _type == "float64"
}

func (field Field) IsString() bool {
	_type := strings.ToLower(field._type())
	return _type == "string" || _type == "str" || _type == ""
}

func (field Field) IsList() bool {
	_type := strings.ToLower(field._type())
	return _type == "list" || _type == "array"
}

type EpisodeScheme struct {
	Fields []Field
}

type AnimeScheme struct {
	Fields []Field
	// Episode EpisodeScheme
}

type PageSchema struct {
	ID         string
	Version    string
	Cloudflare bool
	Anime      AnimeScheme
}

func ReadSchema(filename string) (*PageSchema, error) {
	var out PageSchema

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(content, &out)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s", err)
	}

	return &out, nil
}
