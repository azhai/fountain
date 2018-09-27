package article

import (
	"gopkg.in/yaml.v2"
)

const (
	MODE_DIR  = 0755
	MODE_FILE = 0666
)

func YamlParse(data []byte, storage interface{}) error {
	return yaml.Unmarshal(data, storage)
}

type Section map[string]interface{}

type Setting struct {
	Title   string
	Logo    string
	Lang    string
	Source  string
	Public  string
	Theme   string
	Port    uint
	Limit   int
	Authors map[string]*User
	Layout  *Section
	Repo    *Section
}

func NewSetting() *Setting {
	return &Setting{
		Title:   "",
		Logo:    "",
		Lang:    "",
		Source:  "",
		Public:  "",
		Theme:   "",
		Port:    8080,
		Limit:   20,
		Authors: make(map[string]*User),
		Layout:  new(Section),
		Repo:    new(Section),
	}
}

type User struct {
	ID     string
	Name   string
	Intro  string
	Avatar string
}

type Link struct {
	Url    string
	Anchor string
	Title  string
	Note   string
}

func I18n(val string) string {
	return val
}
