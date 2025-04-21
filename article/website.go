package article

import (
	"container/list"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"
	// "github.com/k0kubun/pp"
)

type Website struct {
	Root       string
	DirList    []string
	DirLinks   map[string][]string
	TagIndexes map[string][]int
	Articles   []*Article
	Conf       *Setting
	Skin       *Theme
	Convert    func(source []byte, format string) []byte
	Debug      func(data ...interface{})
}

func NewWebsite(root string) *Website {
	if root == "" || root == "." || root == "./" {
		root = ""
	} else if root[len(root)-1] != '/' {
		root += "/"
	}
	return &Website{
		Root:       root,
		DirLinks:   make(map[string][]string),
		TagIndexes: make(map[string][]int),
		Conf:       NewSetting(),
	}
}

func (w *Website) LoadConfig(path string) error {
	data, err := os.ReadFile(w.Root + path)
	if err != nil {
		w.Debug("ERROR:", err)
		return err
	}
	if err := yaml.Unmarshal(data, w.Conf); err != nil {
		panic(err)
	}
	return nil
}

func (w *Website) InitTheme() error {
	theme := "default"
	if w.Conf.Theme != "" {
		theme = w.Conf.Theme
	}
	themeDir := fmt.Sprintf("themes/%s/", theme)
	w.Skin = NewTheme(w.Root + themeDir)
	w.Skin.PubDir = w.Root + w.Conf.Public
	w.Skin.FunDict["i18n"] = I18n
	return os.MkdirAll(w.Skin.PubDir, MODE_DIR)
}

func (w *Website) AddArticle(blog *Article, dir, name string) string {
	url := dir + "/" + name + ".html"
	blog.SetDirUrl(dir, url)
	if authorId := blog.Meta.Author; authorId != "" {
		if author, ok := w.Conf.Authors[authorId]; ok {
			blog.Author = author
		}
	}

	if _, ok := w.DirLinks[dir]; !ok {
		w.DirList = append(w.DirList, dir)
		fullDir := w.Skin.PubDir + dir
		w.Debug(fullDir + ":")
		os.MkdirAll(fullDir, MODE_DIR)
	}
	lnk := blog.Archive.ToString(AddUrlPre(dir))
	w.DirLinks[dir] = append(w.DirLinks[dir], lnk)

	idx := len(w.Articles)
	for _, tag := range blog.Meta.Tags {
		w.TagIndexes[tag] = append(w.TagIndexes[tag], idx)
	}

	w.Articles = append(w.Articles, blog)
	return url
}

func (w *Website) SortByDate() {
	sort.Slice(w.Articles, func(i, j int) bool {
		return w.Articles[i].Meta.Date > w.Articles[j].Meta.Date
	})
}

func (w *Website) CreateIndex(pageSize int) error {
	var (
		err      error
		catalogs []*Catelog
		count    = len(w.Articles)
		url      = "index.html"
	)
	lst := list.New()
	for i := 0; i < count; i += pageSize {
		if pageNo := i / pageSize; pageNo > 0 {
			url = fmt.Sprintf("index-%d.html", pageNo)
		}
		stop := i + pageSize
		if stop > count {
			stop = count
		}
		cata := &Catelog{
			Node:  lst.PushBack(url),
			Start: i,
			Stop:  stop,
		}
		catalogs = append(catalogs, cata)
	}
	for _, cata := range catalogs {
		cata.Site = w
		url = cata.Node.Value.(string)
		w.Debug("√", url)
		ctx := Table{"Cata": cata, "Tag": ""}
		err = w.Prepare("", ctx, false)
		if err != nil {
			return err
		}
		err = w.Skin.Render("index", url, ctx)
	}
	return err
}

func (w Website) CreateDirs() error {
	var err error
	for _, name := range w.DirList {
		url := fmt.Sprintf("%s/index.html", name)
		w.Debug("√", url)
		ctx := Table{"Site": w, "Dir": name}
		err = w.Prepare(name, ctx, false)
		if err != nil {
			return err
		}
		err = w.Skin.Render("dir", url, ctx)
	}
	return err
}

func (w Website) CreateTags() error {
	var err error
	createDir := true
	for name := range w.TagIndexes {
		url := fmt.Sprintf("tags/%s.html", name)
		w.Debug("√", url)
		ctx := Table{"Site": w, "Tag": name}
		err = w.Prepare("tags", ctx, createDir)
		if err != nil {
			return err
		}
		createDir = false
		err = w.Skin.Render("tag", url, ctx)
	}
	return err
}

func (w Website) GlobPages(thDir string) ([]string, error) {
	return filepath.Glob(thDir + "pages/*.html")
}

func (w Website) CreatePages(pages []string, thPrelen int) (err error) {
	createDir := true
	for _, p := range pages {
		url := p[thPrelen:]
		w.Debug("√", url)
		ctx := Table{"Site": w, "Tag": ""}
		err = w.Prepare("pages", ctx, createDir)
		if err != nil {
			return
		}
		createDir = false
		if name := filepath.Base(url); strings.HasSuffix(name, ".html") {
			err = w.Skin.Render(name, url, ctx)
		}
	}
	return
}

// func (w Website) GetDirArchives(name string) []*Link {
// 	var arches []*Link
// 	if indexes, ok := w.DirLinks[name]; ok {
// 		for _, idx := range indexes {
// 			lnk := w.Articles[idx].Archive
// 			arches = append(arches, lnk)
// 		}
// 	}
// 	return arches
// }

func (w Website) GetTagArchives(name string) []*Link {
	var arches []*Link
	if indexes, ok := w.TagIndexes[name]; ok {
		for _, idx := range indexes {
			lnk := w.Articles[idx].Archive
			arches = append(arches, lnk)
		}
	}
	return arches
}

func (w Website) Prepare(dir string, cxt Table, createDir bool) (err error) {
	if createDir {
		err = os.MkdirAll(w.Skin.PubDir+dir, MODE_DIR)
		if err != nil {
			return
		}
	}
	cxt["Conf"], cxt["Dir"] = w.Conf, dir
	cxt["UrlPre"] = AddUrlPre(dir)
	cxt["ArchDirs"] = w.DirLinks
	return
}

func (w *Website) RenderFile(tpl, dir, url string, blog *Article) error {
	ctx := Table{"Blog": blog, "Tag": ""}
	err := w.Prepare(dir, ctx, false)
	if err != nil {
		w.Debug("ERROR:", err)
	}
	err = w.Skin.Render(tpl, url, ctx)
	if err != nil {
		w.Debug("ERROR:", err)
	}
	return err
}

func (w *Website) ProcFile(fullpath, path string) error {
	blog := NewArticle()
	name, err := blog.ParseFile(fullpath)
	if err != nil {
		w.Debug("ERROR:", err)
		return err
	}
	source := []byte(blog.Source)
	content := w.Convert(source, blog.Format)
	blog.Content = string(content)
	dir := filepath.Dir(path)
	url := w.AddArticle(blog, dir, name)
	w.Debug(path, "->", url)
	return nil
}

func (w *Website) CreateWalkFunc() filepath.WalkFunc {
	prelen := len(w.Root + w.Conf.Source)
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			w.Debug("ERROR:", err)
			return err
		}
		// 跳过.开头的目录或文件
		base := filepath.Base(path)
		if base[0] == '.' {
			if info.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			}
		} else {
			if info.IsDir() {
				return nil
			} else {
				return w.ProcFile(path, path[prelen:])
			}
		}
	}
}

func (w *Website) BuildFiles() error {
	srcDir := w.Root + w.Conf.Source
	walkFunc := w.CreateWalkFunc()
	err := filepath.Walk(srcDir, walkFunc)
	// pp.Println("Articles", w.Articles)
	for _, blog := range w.Articles {
		dir, url := "", ""
		if blog != nil && blog.Archive != nil {
			dir, url = blog.Archive.Dir, blog.Archive.Url
		}
		w.RenderFile("article", dir, url, blog)
	}
	if err != nil {
		return err
	}
	w.Debug("Index:")
	w.SortByDate()
	w.CreateIndex(w.Conf.Limit)

	if w.Skin.HasTemplate("dir") {
		w.Debug("Dirs:")
		w.CreateDirs()
	}
	if w.Skin.HasTemplate("tag") {
		w.Debug("Tags:")
		w.CreateTags()
	}

	var pages []string
	thDir := w.Skin.GetDir()
	pages, err = w.GlobPages(thDir)
	if err == nil && len(pages) > 0 {
		w.Debug("Pages:")
		w.CreatePages(pages, len(thDir))
	}

	w.Skin.CopyAssets("static")
	if w.Skin.WithSide {
		// 这个放在复制静态文件之后，可以不用创建目录
		w.Debug("Static:")
		path := "static/js/app.js"
		w.Skin.CreateSidebar(path, w.DirLinks)
		w.Debug("√", path)
	}
	return err
}
