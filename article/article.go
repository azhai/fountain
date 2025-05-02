package article

import (
	"bytes"
	"container/list"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	// "github.com/k0kubun/pp"
)

const (
	SEP_META              = "---"
	SEP_MORE              = "<!--more-->"
	SEP_OUTLINE           = "<!--outline-->"
	SEP_SAFE_MODE         = "<!-- raw HTML omitted -->"
	EXT_MARKDOWN          = ".md"
	EXT_RESTRUCTURED_TEXT = ".rst"
)

type MetaData struct {
	Title  string   `yaml:"title"`
	Slug   string   `yaml:"slug,omitempty"`
	Date   string   `yaml:"date,omitempty"`
	Update string   `yaml:"update,omitempty"`
	Author string   `yaml:"author,omitempty"`
	Tags   []string `yaml:"tags,omitempty"`
	Draft  bool     `yaml:"draft,omitempty"`
	Weight int      `yaml:"weight,omitempty"`
}

type Catelog struct {
	Site  *Website
	Node  *list.Element
	Start int
	Stop  int
}

func CreateCatelogs(count, pageSize int) []*Catelog {
	var catalogs []*Catelog
	lst, url := list.New(), "index.html"
	for i := 0; i < count; i += pageSize {
		if pageNo := i / pageSize; pageNo > 0 {
			url = fmt.Sprintf("index-%d.html", pageNo)
		}
		stop := min(i+pageSize, count)
		cata := &Catelog{
			Node:  lst.PushBack(url),
			Start: i,
			Stop:  stop,
		}
		catalogs = append(catalogs, cata)
	}
	return catalogs
}

func (c Catelog) GetArchives() []*Link {
	links := make([]*Link, 0)
	for i := c.Start; i < c.Stop; i++ {
		art := c.Site.Articles[i]
		links = append(links, art.Archive)
	}
	return links
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
	Outline string
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

func (a *Article) SetDirUrl(dir, url string) *Link {
	a.Archive = &Link{
		Dir:   dir,
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
		// fmt.Println(text)
		if err := yaml.Unmarshal([]byte(text), a.Meta); err != nil {
			panic(err)
		}
		if a.Meta.Author != "" {
			a.SetDummyAuthor(a.Meta.Author)
		}
		// pp.Println(a.Meta)
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

func (a *Article) SplitContent(data []byte) error {
	content := strings.TrimSpace(string(data))
	pieces := strings.SplitN(content, SEP_SAFE_MODE, 2)
	if len(pieces) != 2 {
		pieces = strings.SplitN(content, SEP_OUTLINE, 2)
	}
	if len(pieces) == 2 {
		a.Outline, a.Content = pieces[0], pieces[1]
	} else {
		a.Content = content
	}
	return nil
}

func (a *Article) ParseFile(path string) (string, error) {
	data, err := os.ReadFile(path)
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
