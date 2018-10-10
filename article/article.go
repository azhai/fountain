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

func (this *Catelog) GetArchives() []*Link {
	return this.Site.Archives[this.Start:this.Stop]
}

func (this *Catelog) GetNext() string {
	link := "下一页"
	if node := this.Node.Next(); node != nil {
		url := node.Value.(string)
		link = fmt.Sprintf("<a href=\"./%s\">%s</a>", url, link)
	}
	return link
}

func (this *Catelog) GetPrev() string {
	link := "上一页"
	if node := this.Node.Prev(); node != nil {
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

func (this *Article) SetFormat(ext string) {
	switch ext {
	case EXT_MARKDOWN:
		this.Format = "markdown"
	default:
		this.Format = ""
	}
}

func (this *Article) SetUrl(url string) *Link {
	this.Archive = &Link{
		Url:   url,
		Title: this.Meta.Title,
		Note:  this.Meta.Date,
	}
	return this.Archive
}

func (this *Article) SetDummyAuthor(id string) *User {
	this.Author = &User{
		ID:     id,
		Name:   id,
		Intro:  "",
		Avatar: "",
	}
	return this.Author
}

func (this *Article) SetData(name string, chunk []byte) int {
	text := strings.TrimSpace(string(chunk))
	length := len(text)
	if length == 0 {
		return 0
	}
	if name == "Meta" {
		YamlParse([]byte(text), &this.Meta)
		if this.Meta.Author != "" {
			this.SetDummyAuthor(this.Meta.Author)
		}
	} else {
		this.Source = text
	}
	return length
}

func (this *Article) SplitSource(data []byte, times int) error {
	var (
		size, length int
		idx          = 0
		sep          = []byte(SEP_META + "\n")
		offset       = len(sep)
	)
	this.Meta.Title = ""
	for {
		times--
		size = bytes.Index(data[idx:], sep)
		//不是最后一段，元数据还没有找到
		if size >= 0 && times >= 0 && this.Meta.Title == "" {
			length = this.SetData("Meta", data[idx:idx+size])
			//空段或刚找到元数据
			if length == 0 || this.Meta.Title != "" {
				idx += size + offset
				continue
			}
		}
		this.SetData("Source", data[idx:])
		break
	}
	return nil
}

func (this *Article) ParseFile(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	err = this.SplitSource(data, 2)
	if err != nil {
		return "", err
	}
	this.SetFormat(filepath.Ext(path))
	name := this.Meta.Slug
	if name == "" {
		base := filepath.Base(path)
		ext := filepath.Ext(base)
		name = base[:len(base)-len(ext)]
	}
	return name, nil
}
