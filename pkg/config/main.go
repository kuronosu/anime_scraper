package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type DetailSchema struct {
	Fields []Field
}

type PageSchema struct {
	ID         string
	Version    string
	Cloudflare bool
	Detail     DetailSchema
	List       ListSchema
}

func ReadSchema(filename string) (*PageSchema, error) {
	var out PageSchema

	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(content, &out)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s", err)
	}

	return &out, nil
}
