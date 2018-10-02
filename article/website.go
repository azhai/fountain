package article

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Website struct {
	RootDir  string
	Archives []*Link
	Tags     map[string]([]int)
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
		RootDir: root,
		Tags:    make(map[string]([]int)),
		Conf:    NewSetting(),
		Skin:    NewTheme(),
	}
}

func (this *Website) LoadConfig(path string) error {
	data, err := ioutil.ReadFile(this.RootDir + path)
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
	this.Skin.Dir = this.RootDir + fmt.Sprintf("themes/%s/", theme)
	this.Skin.OutDir = this.RootDir + this.Conf.Public
	this.Skin.FuncMap["i18n"] = I18n
	return os.MkdirAll(this.Skin.OutDir, MODE_DIR)
}

func (this *Website) AddArticle(blog *Article, url string) *Link {
	arch := blog.SetUrl(url)
	idx := len(this.Archives)
	this.Archives = append(this.Archives, arch)
	for _, name := range blog.Meta.Tags {
		this.Tags[name] = append(this.Tags[name], idx)
	}
	return arch
}

func (this *Website) GetTargetDir(path, name string) string {
	srcDir := this.RootDir + this.Conf.Source
	pubDir := this.RootDir + this.Conf.Public
	dir := filepath.Dir(path)[len(srcDir):]
	url := dir + "/" + name + ".html"
	this.Debug(path, "->", pubDir+url)
	os.MkdirAll(pubDir+dir, MODE_DIR)
	return url
}

func (this *Website) GetTagArchives(name string) []*Link {
	var arches []*Link
	if indexes, ok := this.Tags[name]; ok {
		for _, idx := range indexes {
			lnk := this.Archives[idx]
			arches = append(arches, lnk)
		}
	}
	return arches
}

func (this *Website) ProcFile(blog *Article, path string) error {
	name, err := blog.ParseFile(path)
	if err != nil {
		this.Debug("ERROR:", err)
		return err
	}
	url := this.GetTargetDir(path, name)
	source := []byte(blog.Source)
	content := this.Convert(source, blog.Format)
	blog.Content = string(content)
	this.AddArticle(blog, url)
	ctx := Table{"Blog": blog, "Tag": "", "Conf": this.Conf}
	err = this.Skin.Render("article", url, ctx)
	if err != nil {
		this.Debug("ERROR:", err)
	}
	return err
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
		ctx := Table{"Cata": cata, "Tag": "", "Conf": this.Conf}
		err = this.Skin.Render("index", url, ctx)
	}
	return err
}

func (this *Website) CreateTags() error {
	var err error
	pubDir := this.RootDir + this.Conf.Public
	err = os.MkdirAll(pubDir+"tag", MODE_DIR)
	if err != nil {
		return err
	}
	for name := range this.Tags {
		url := fmt.Sprintf("tag/%s.html", name)
		this.Debug("√", url)
		ctx := Table{"Site": this, "Tag": name, "Conf": this.Conf}
		err = this.Skin.Render("tag", url, ctx)
	}
	return err

}

func (this *Website) CreateWalkFunc() filepath.WalkFunc {
	blog := NewArticle()
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
		}
		if info.IsDir() {
			this.Debug(path + ":")
			return nil
		}
		return this.ProcFile(blog, path)
	}
}

func (this *Website) BuildFiles() error {
	srcDir := this.RootDir + this.Conf.Source
	walkFunc := this.CreateWalkFunc()
	err := filepath.Walk(srcDir, walkFunc)
	if err == nil {
		this.Debug("Index:")
		this.CreateIndex(this.Conf.Limit)
		this.Debug("Tags:")
		this.CreateTags()
		this.Skin.CopyAssets("static")
	}
	return err
}
