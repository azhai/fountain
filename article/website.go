package article

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

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
	Debug      func(data ...any)
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
	theme := w.Conf.GetTheme()
	themeDir := fmt.Sprintf("themes/%s/", theme)
	w.Skin = NewTheme(w.Root + themeDir)
	w.Skin.PubDir = w.Root + w.Conf.Public
	w.Skin.FunDict["i18n"] = I18n
	return os.MkdirAll(w.Skin.PubDir, DefaultDirMode)
}

func (w *Website) AddArticle(blog *Article, dir string) error {
	if authorId := blog.Meta.Author; authorId != "" {
		if author, ok := w.Conf.Authors[authorId]; ok {
			blog.Author = author
		}
	}

	if _, ok := w.DirLinks[dir]; !ok {
		w.DirList = append(w.DirList, dir)
		fullDir := w.Skin.PubDir + dir
		w.Debug(fullDir + ":")
		os.MkdirAll(fullDir, DefaultDirMode)
	}
	lnk := blog.Archive.ToString(AddUrlPre(dir))
	w.DirLinks[dir] = append(w.DirLinks[dir], lnk)

	idx := len(w.Articles)
	for _, tag := range blog.Meta.Tags {
		w.TagIndexes[tag] = append(w.TagIndexes[tag], idx)
	}

	w.Articles = append(w.Articles, blog)
	return nil
}

func (w *Website) SortByDate() {
	sort.Slice(w.Articles, func(i, j int) bool {
		return w.Articles[i].Meta.Date > w.Articles[j].Meta.Date
	})
}

func (w *Website) CreateIndex(pageSize int) (err error) {
	ctx := Table{"Cata": "", "Tag": ""}
	if err = w.Prepare("", ctx, true); err != nil {
		return
	}
	catalogs := CreateCatelogs(len(w.Articles), pageSize)
	for _, cata := range catalogs {
		url := cata.Node.Value.(string)
		w.Debug("√", url)
		cata.Site = w
		ctx["Cata"] = cata
		err = w.Skin.Render("index", url, ctx)
	}
	return
}

func (w Website) CreateDirs() error {
	var err error
	for _, name := range w.DirList {
		url := fmt.Sprintf("%s/index.html", name)
		w.Debug("√", url)
		ctx := Table{"Site": w, "Dir": name}
		err = w.Prepare(name, ctx, true)
		if err != nil {
			return err
		}
		err = w.Skin.Render("dir", url, ctx)
	}
	return err
}

func (w Website) CreateTags() (err error) {
	ctx := Table{"Site": w, "Tag": ""}
	if err = w.Prepare("tags", ctx, true); err != nil {
		return
	}
	for name := range w.TagIndexes {
		url := fmt.Sprintf("tags/%s.html", name)
		w.Debug("√", url)
		ctx["Tag"] = name
		err = w.Skin.Render("tag", url, ctx)
	}
	return err
}

func (w Website) GlobPages(thDir string) ([]string, error) {
	htmlPages, err1 := filepath.Glob(thDir + "pages/*.html")
	mdPages, err2 := filepath.Glob(thDir + "pages/*.md")
	if err1 != nil {
		return mdPages, err2
	}
	if err2 != nil {
		return htmlPages, err1
	}
	return append(htmlPages, mdPages...), nil
}

func (w Website) CreatePages(thDir string, thPrelen int) (err error) {
	ctx := Table{"Site": w, "Tag": ""}
	if err = w.Prepare("pages", ctx, true); err != nil {
		return
	}
	htmlPages, _ := filepath.Glob(thDir + "pages/*.html")
	for _, path := range htmlPages {
		url := path[thPrelen:]
		w.Debug("√", url)
		err = w.Skin.Render(url, url, ctx)
	}

	var blog *Article
	mdPages, _ := filepath.Glob(thDir + "pages/*.md")
	for _, path := range mdPages {
		url := path[thPrelen:len(path)-3] + ".html"
		w.Debug("√", url)
		if blog, err = w.ProcFile(path, url); err != nil {
			return
		}
		err = w.RenderBlog("channel", blog)
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
		err = os.MkdirAll(w.Skin.PubDir+dir, DefaultDirMode)
		if err != nil {
			return
		}
	}
	cxt["Conf"], cxt["Dir"] = w.Conf, dir
	cxt["Footer"] = w.Conf.GetFooter()
	cxt["Github"] = w.Conf.Github
	cxt["UrlPre"] = AddUrlPre(dir)
	cxt["ArchDirs"] = w.DirLinks
	return
}

func (w *Website) RenderBlog(tpl string, blog *Article) error {
	if blog == nil || blog.Archive == nil {
		return fmt.Errorf("invalid blog: %v", blog)
	}
	dir, url := blog.Archive.Dir, blog.Archive.Url
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

func (w *Website) ProcFile(fullpath, path string) (*Article, error) {
	blog := NewArticle()
	name, err := blog.ParseFile(fullpath)
	if err != nil {
		w.Debug("ERROR:", err)
		return nil, err
	}

	source := []byte(blog.Source)
	content := w.Convert(source, blog.Format)
	_ = blog.SplitContent(content)

	dir := filepath.Dir(path)
	url := dir + "/" + name + ".html"
	blog.SetDirUrl(dir, url)
	w.Debug(path, "->", url)
	return blog, err
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
				var blog *Article
				blog, err = w.ProcFile(path, path[prelen:])
				dir := blog.Archive.Dir
				err = w.AddArticle(blog, dir)
				return err
			}
		}
	}
}

func (w *Website) BuildFiles() error {
	thDir := w.Skin.GetDir()
	w.Debug("Pages:")
	err := w.CreatePages(thDir, len(thDir))
	if err != nil {
		return err
	}

	walkFunc := w.CreateWalkFunc()
	srcDir := w.Root + w.Conf.Source
	err = filepath.Walk(srcDir, walkFunc)
	if err != nil {
		return err
	}
	// pp.Println("Articles", w.Articles)

	for _, blog := range w.Articles {
		err = w.RenderBlog("article", blog)
		if err != nil {
			return err
		}
	}

	w.Debug("Index:")
	w.SortByDate()
	if err = w.CreateIndex(w.Conf.Limit); err != nil {
		return err
	}

	if w.Skin.HasTemplate("dir") {
		w.Debug("Dirs:")
		if err = w.CreateDirs(); err != nil {
			return err
		}
	}
	if w.Skin.HasTemplate("tag") {
		w.Debug("Tags:")
		if err = w.CreateTags(); err != nil {
			return err
		}
	}

	err = w.Skin.CopyAssets("static")
	if err != nil {
		return err
	}
	if w.Skin.WithSide {
		// 这个放在复制静态文件之后，可以不用创建目录
		w.Debug("Static:")
		path := "static/js/app.js"
		err = w.Skin.CreateSidebar(path, w.DirLinks)
		w.Debug("√", path)
	}
	return err
}
