package config

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type Regex struct {
	Pattern string
	Group   int
}

type Field struct {
	Name     string
	Selector string
	Type     *string
	Attr     *string
	Replace  map[string]string
	Remove   []string
	Regex    []Regex
}

func (field Field) Compile(value string) (interface{}, error) {
	if field.IsInt() {
		v, err := strconv.Atoi(value)
		return v, err
	}
	if field.IsFloat() {
		v, err := strconv.ParseFloat(value, 64)
		return v, err
	}
	if field.IsString() {
		return value, nil
	}
	return nil, fmt.Errorf("unsupported type %s", field._type())
}

func (field Field) _type() string {
	if field.Type == nil {
		return "string"
	}
	return *field.Type
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
