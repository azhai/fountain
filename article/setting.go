package article

import (
	"fmt"
	"time"
)

const (
	DefaultDirMode  = 0755
	DefaultFileMode = 0666
	DefaultTheme    = "default"
)

type Table map[string]any

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
	Show_toc bool
	Top_tags []string
	Github   string
	Authors  map[string]*User
	Footer   *Footer
	Layout   *Table
}

type Footer struct {
	Copyright  string
	From_year  int
	Curr_year  int
	Cn_cert_no string
}

func NewSetting() *Setting {
	return &Setting{
		Title:   "我的博客",
		Lang:    "zh",
		Source:  "source/",
		Public:  "public/",
		Port:    8080,
		Limit:   10,
		Authors: make(map[string]*User),
		Footer:  new(Footer),
		Layout:  new(Table),
	}
}

func (s *Setting) GetTheme() string {
	if s.Theme == "" {
		s.Theme = DefaultTheme
	}
	return s.Theme
}

func (s *Setting) GetFooter() *Footer {
	if s.Footer == nil {
		return nil
	}
	if s.Footer.Copyright == "" && s.Footer.Cn_cert_no == "" {
		s.Footer = nil
	} else {
		s.Footer.Curr_year = time.Now().Year()
	}
	return s.Footer
}

type User struct {
	ID     string
	Name   string
	Intro  string
	Avatar string
}

type Link struct {
	Dir    string
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
	tpl := `<a href="%s" class="art-link" title="%s">%s</a>`
	return fmt.Sprintf(tpl, url, l.Note, l.Title)
}

func I18n(val string) string {
	return val
}
