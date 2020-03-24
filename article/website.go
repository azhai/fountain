package article

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Categoris map[string]([]int)

type Website struct {
	Root     string
	Archives []*Link
	ArchDirs Categoris
	ArchTags Categoris
	Conf     *Setting
	Skin     *Theme
	Convert  func(source []byte, format string) []byte
	Debug    func(data ...interface{})
}

func NewWebsite(root string) *Website {
	if root == "" || root == "." || root == "./" {
		root = ""
	} else if root[len(root)-1] != '/' {
		root += "/"
	}
	return &Website{
		Root:     root,
		ArchDirs: make(Categoris),
		ArchTags: make(Categoris),
		Conf:     NewSetting(),
	}
}

func (w *Website) LoadConfig(path string) error {
	data, err := ioutil.ReadFile(w.Root + path)
	if err != nil {
		w.Debug("ERROR:", err)
		return err
	}
	YamlParse(data, &w.Conf)
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
	w.Skin.FunDict["getArchiveString"] = w.GetArchiveString
	return os.MkdirAll(w.Skin.PubDir, MODE_DIR)
}

func (w *Website) AddArticle(blog *Article, dir, name string) string {
	url := dir + "/" + name + ".html"
	arch := blog.SetUrl(url)
	idx := len(w.Archives)
	w.Archives = append(w.Archives, arch)
	if authorId := blog.Meta.Author; authorId != "" {
		if author, ok := w.Conf.Authors[authorId]; ok {
			blog.Author = author
		}
	}
	for _, name := range blog.Meta.Tags {
		w.ArchTags[name] = append(w.ArchTags[name], idx)
	}
	if _, ok := w.ArchDirs[dir]; !ok {
		fullDir := w.Skin.PubDir + dir
		w.Debug(fullDir + ":")
		os.MkdirAll(fullDir, MODE_DIR)
	}
	w.ArchDirs[dir] = append(w.ArchDirs[dir], idx)
	return url
}

func (w Website) GetArchive(idx int) *Link {
	count := len(w.Archives)
	if idx >= 0 && idx < count {
		return w.Archives[idx]
	}
	return &Link{}
}

func (w Website) GetArchiveString(idx int) string {
	lnk := w.GetArchive(idx)
	if lnk == nil || lnk.Title == "" {
		return ""
	}
	return lnk.ToString("@URLPRE@")
}

func (w Website) GetTagArchives(name string) []*Link {
	var arches []*Link
	if indexes, ok := w.ArchTags[name]; ok {
		for _, idx := range indexes {
			lnk := w.Archives[idx]
			arches = append(arches, lnk)
		}
	}
	return arches
}

func (w *Website) CreateIndex(pageSize int) error {
	var (
		err      error
		catalogs []*Catelog
		count    = len(w.Archives)
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

func (w Website) CreateTags() error {
	var err error
	createDir := true
	for name := range w.ArchTags {
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

func (w Website) CreatePages() error {
	thDir, createDir := w.Skin.GetDir(), true
	pages, err := filepath.Glob(thDir + "pages/*.html")
	if err != nil {
		return err
	}
	thPrelen := len(thDir)
	for _, p := range pages {
		url := p[thPrelen:]
		w.Debug("√", url)
		ctx := Table{"Site": w, "Tag": ""}
		err = w.Prepare("pages", ctx, createDir)
		if err != nil {
			return err
		}
		createDir = false
		if name := filepath.Base(url); strings.HasSuffix(name, ".html") {
			err = w.Skin.Render(name, url, ctx)
		}
	}
	return err
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
	return
}

func (w *Website) RenderFile(tpl, url, path string) error {
	dir := filepath.Dir(path)
	ctx := Table{"Blog": nil, "Tag": ""}
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

func (w *Website) RenderArticle(blog *Article, path string, prelen int) error {
	name, err := blog.ParseFile(path)
	if err != nil {
		w.Debug("ERROR:", err)
		return err
	}
	source := []byte(blog.Source)
	content := w.Convert(source, blog.Format)
	blog.Content = string(content)
	path = path[prelen:]
	dir := filepath.Dir(path)
	url := w.AddArticle(blog, dir, name)
	w.Debug(path, "->", url)
	return w.RenderFile("article", url, path)
}

func (w *Website) ProcFile(blog *Article, path string, prelen int) error {
	name, err := blog.ParseFile(path)
	if err != nil {
		w.Debug("ERROR:", err)
		return err
	}
	source := []byte(blog.Source)
	content := w.Convert(source, blog.Format)
	blog.Content = string(content)
	path = path[prelen:]
	dir := filepath.Dir(path)
	url := w.AddArticle(blog, dir, name)
	w.Debug(path, "->", url)
	ctx := Table{"Blog": blog, "Tag": ""}
	err = w.Prepare(dir, ctx, false)
	if err != nil {
		w.Debug("ERROR:", err)
	}
	err = w.Skin.Render("article", url, ctx)
	if err != nil {
		w.Debug("ERROR:", err)
	}
	return err
}

func (w *Website) CreateWalkFunc() filepath.WalkFunc {
	blog := NewArticle()
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
				return w.ProcFile(blog, path, prelen)
			}
		}
	}
}

func (w *Website) BuildFiles() error {
	srcDir := w.Root + w.Conf.Source
	walkFunc := w.CreateWalkFunc()
	err := filepath.Walk(srcDir, walkFunc)
	if err == nil {
		w.Debug("Index:")
		w.CreateIndex(w.Conf.Limit)
		w.Debug("Tags:")
		w.CreateTags()
		w.Debug("Pages:")
		w.CreatePages()
		w.Skin.CopyAssets("static")
		//这个放在复制静态文件之后，可以不用创建目录
		w.Debug("Static:")
		path := "static/js/app.js"
		w.Skin.CreateSidebar(path, w.ArchDirs)
		w.Debug("√", path)
	}
	return err
}
