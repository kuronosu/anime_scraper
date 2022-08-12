package config

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

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

func (field Field) GetRegex() []Regex {
	return field.Regex
}

func (field Field) GetRemove() []string {
	return field.Remove
}

func (field Field) GetReplace() map[string]string {
	return field.Replace
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
		log.Printf("error compiling field %s: %s", field.Name, err)
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

	value := ApplyFilters(field, rawValue)

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
