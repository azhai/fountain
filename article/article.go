package article

import (
	"bytes"
	"container/list"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

const (
	SEP_META              = "---"
	SEP_MORE              = "<!--more-->"
	EXT_MARKDOWN          = ".md"
	EXT_RESTRUCTURED_TEXT = ".rst"
)

type MetaData struct {
	Title  string
	Slug   string
	Date   string
	Update string
	Author string
	Tags   []string
	Draft  bool
	Weight int
}

type Catelog struct {
	Site  *Website
	Node  *list.Element
	Start int
	Stop  int
}

func (c Catelog) GetArchives() []*Link {
	return c.Site.Archives[c.Start:c.Stop]
}

func (c Catelog) GetNext() string {
	link := "下一页"
	if node := c.Node.Next(); node != nil {
		url := node.Value.(string)
		link = fmt.Sprintf("<a href=\"./%s\">%s</a>", url, link)
	}
	return link
}

func (c Catelog) GetPrev() string {
	link := "上一页"
	if node := c.Node.Prev(); node != nil {
		url := node.Value.(string)
		link = fmt.Sprintf("<a href=\"./%s\">%s</a>", url, link)
	}
	return link
}

type Article struct {
	Meta    *MetaData
	Author  *User
	Archive *Link
	Format  string
	Source  string
	Content string
}

func NewArticle() *Article {
	return &Article{
		Meta: &MetaData{
			Date: time.Now().Format("2006-01-02"),
		},
	}
}

func (a *Article) SetFormat(ext string) {
	switch ext {
	case EXT_MARKDOWN:
		a.Format = "markdown"
	default:
		a.Format = ""
	}
}

func (a *Article) SetUrl(url string) *Link {
	a.Archive = &Link{
		Url:   url,
		Title: a.Meta.Title,
		Note:  a.Meta.Date,
	}
	return a.Archive
}

func (a *Article) SetDummyAuthor(id string) *User {
	a.Author = &User{
		ID:     id,
		Name:   id,
		Intro:  "",
		Avatar: "",
	}
	return a.Author
}

func (a *Article) SetData(name string, chunk []byte) int {
	text := strings.TrimSpace(string(chunk))
	length := len(text)
	if length == 0 {
		return 0
	}
	if name == "Meta" {
		ParseConfData([]byte(text), &a.Meta)
		if a.Meta.Author != "" {
			a.SetDummyAuthor(a.Meta.Author)
		}
	} else {
		a.Source = text
	}
	return length
}

func (a *Article) SplitSource(data []byte, times int) error {
	var (
		size, length int
		idx          = 0
		sep          = []byte(SEP_META + "\n")
		offset       = len(sep)
	)
	a.Meta.Title = ""
	for {
		times--
		size = bytes.Index(data[idx:], sep)
		// 不是最后一段，元数据还没有找到
		if size >= 0 && times >= 0 && a.Meta.Title == "" {
			length = a.SetData("Meta", data[idx:idx+size])
			// 空段或刚找到元数据
			if length == 0 || a.Meta.Title != "" {
				idx += size + offset
				continue
			}
		}
		a.SetData("Source", data[idx:])
		break
	}
	return nil
}

func (a *Article) ParseFile(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	err = a.SplitSource(data, 2)
	if err != nil {
		return "", err
	}
	a.SetFormat(filepath.Ext(path))
	name := a.Meta.Slug
	if name == "" {
		base := filepath.Base(path)
		ext := filepath.Ext(base)
		name = base[:len(base)-len(ext)]
	}
	return name, nil
}
