package article

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

const (
	MODE_DIR  = 0755
	MODE_FILE = 0666
)

func YamlParse(data []byte, storage interface{}) error {
	return yaml.Unmarshal(data, storage)
}

type Table map[string]interface{}

type Setting struct {
	Title    string
	Logo     string
	Lang     string
	Source   string
	Public   string
	Theme    string
	Port     uint
	Limit    int
	Sort     string
	Top_tags []string
	Authors  map[string]*User
	Layout   *Table
	Github   *Table
}

func NewSetting() *Setting {
	return &Setting{
		Title:   "我的博客",
		Lang:    "zh",
		Source:  "source/",
		Public:  "public/",
		Theme:   "default",
		Port:    8080,
		Limit:   10,
		Authors: make(map[string]*User),
		Layout:  new(Table),
		Github:  new(Table),
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

func (l Link) ToString(urlpre string) string {
	url := urlpre + "/" + l.Url
	if l.Anchor != "" {
		url += "#" + l.Anchor
	}
	tpl := `<a href="%s" title="%s">%s</a>`
	return fmt.Sprintf(tpl, url, l.Note, l.Title)
}

func I18n(val string) string {
	return val
}
