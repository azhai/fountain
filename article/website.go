package article

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

func (this *Website) LoadConfig(path string) error {
	data, err := ioutil.ReadFile(this.Root + path)
	if err != nil {
		this.Debug("ERROR:", err)
		return err
	}
	YamlParse(data, &this.Conf)
	return nil
}

func (this *Website) InitTheme() error {
	theme := "default"
	if this.Conf.Theme != "" {
		theme = this.Conf.Theme
	}
	themeDir := fmt.Sprintf("themes/%s/", theme)
	this.Skin = NewTheme(this.Root + themeDir)
	this.Skin.PubDir = this.Root + this.Conf.Public
	this.Skin.FunDict["i18n"] = I18n
	this.Skin.FunDict["getArchiveString"] = this.GetArchiveString
	return os.MkdirAll(this.Skin.PubDir, MODE_DIR)
}

func (this *Website) AddArticle(blog *Article, dir, name string) string {
	url := dir + "/" + name + ".html"
	arch := blog.SetUrl(url)
	idx := len(this.Archives)
	this.Archives = append(this.Archives, arch)
	if authorId := blog.Meta.Author; authorId != "" {
		if author, ok := this.Conf.Authors[authorId]; ok {
			blog.Author = author
		}
	}
	for _, name := range blog.Meta.Tags {
		this.ArchTags[name] = append(this.ArchTags[name], idx)
	}
	if _, ok := this.ArchDirs[dir]; !ok {
		fullDir := this.Skin.PubDir + dir
		this.Debug(fullDir + ":")
		os.MkdirAll(fullDir, MODE_DIR)
	}
	this.ArchDirs[dir] = append(this.ArchDirs[dir], idx)
	return url
}

func (this *Website) GetArchive(idx int) *Link {
	count := len(this.Archives)
	if idx >= 0 && idx < count {
		return this.Archives[idx]
	}
	return &Link{}
}

func (this *Website) GetArchiveString(idx int) string {
	lnk := this.GetArchive(idx)
	if lnk == nil || lnk.Title == "" {
		return ""
	}
	return lnk.ToString("@URLPRE@")
}

func (this *Website) GetTagArchives(name string) []*Link {
	var arches []*Link
	if indexes, ok := this.ArchTags[name]; ok {
		for _, idx := range indexes {
			lnk := this.Archives[idx]
			arches = append(arches, lnk)
		}
	}
	return arches
}

func (this *Website) CreateIndex(pageSize int) error {
	var (
		err      error
		catalogs []*Catelog
		count    = len(this.Archives)
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
		cata.Site = this
		url = cata.Node.Value.(string)
		this.Debug("√", url)
		ctx := Table{"Cata": cata, "Tag": ""}
		err = this.Prepare("", ctx, false)
		if err != nil {
			return err
		}
		err = this.Skin.Render("index", url, ctx)
	}
	return err
}

func (this *Website) CreateTags() error {
	var err error
	err = os.MkdirAll(this.Skin.PubDir+"tag", MODE_DIR)
	if err != nil {
		return err
	}
	for name := range this.ArchTags {
		url := fmt.Sprintf("tag/%s.html", name)
		this.Debug("√", url)
		ctx := Table{"Site": this, "Tag": name}
		err = this.Prepare("tag", ctx, false)
		if err != nil {
			return err
		}
		err = this.Skin.Render("tag", url, ctx)
	}
	return err
}

func (this *Website) Prepare(dir string, cxt Table, createDir bool) (err error) {
	if createDir {
		err = os.MkdirAll(this.Skin.PubDir+dir, MODE_DIR)
		if err != nil {
			return
		}
	}
	cxt["Dir"] = dir
	cxt["UrlPre"] = this.Skin.AddUrlPre(dir)
	cxt["Conf"] = this.Conf
	return
}

func (this *Website) ProcFile(blog *Article, path string, prelen int) error {
	name, err := blog.ParseFile(path)
	if err != nil {
		this.Debug("ERROR:", err)
		return err
	}
	source := []byte(blog.Source)
	content := this.Convert(source, blog.Format)
	blog.Content = string(content)
	path = path[prelen:]
	dir := filepath.Dir(path)
	url := this.AddArticle(blog, dir, name)
	this.Debug(path, "->", url)
	ctx := Table{"Blog": blog, "Tag": ""}
	err = this.Prepare(dir, ctx, false)
	if err != nil {
		this.Debug("ERROR:", err)
	}
	err = this.Skin.Render("article", url, ctx)
	if err != nil {
		this.Debug("ERROR:", err)
	}
	return err
}

func (this *Website) CreateWalkFunc() filepath.WalkFunc {
	blog := NewArticle()
	prelen := len(this.Root + this.Conf.Source)
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			this.Debug("ERROR:", err)
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
				return this.ProcFile(blog, path, prelen)
			}
		}
	}
}

func (this *Website) BuildFiles() error {
	srcDir := this.Root + this.Conf.Source
	walkFunc := this.CreateWalkFunc()
	err := filepath.Walk(srcDir, walkFunc)
	if err == nil {
		this.Debug("Index:")
		this.CreateIndex(this.Conf.Limit)
		this.Debug("Tags:")
		this.CreateTags()
		this.Skin.CopyAssets("static")
		//这个放在复制静态文件之后，可以不用创建目录
		path := "static/js/app.js"
		this.Skin.CreateSidebar(path, this.ArchDirs)
	}
	return err
}
