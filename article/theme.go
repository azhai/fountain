package article

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"fountain/utils"
)

func AddUrlPre(dir string) string {
	prefix := "."
	if len(dir) > 0 && dir != "." {
		times := strings.Count(dir, "/") + 1
		prefix = strings.Repeat("../", times)
		prefix = prefix[:len(prefix)-1]
	}
	return prefix
}

type Theme struct {
	dir      string // 结尾有斜杠
	PubDir   string // 结尾有斜杠
	FunDict  template.FuncMap
	TplDict  map[string]*template.Template
	WithSide bool
}

func NewTheme(dir string) *Theme {
	return &Theme{
		dir:     dir,
		FunDict: make(map[string]interface{}),
		TplDict: make(map[string]*template.Template),
	}
}

func (t *Theme) HasTemplate(name string) bool {
	if _, ok := t.TplDict[name]; ok {
		return true
	}
	return utils.PathExist(filepath.Join(t.dir, name+".html"))
}

func (t *Theme) GetOrCreate(name, path string, incl bool) *template.Template {
	var err error
	if tpl, ok := t.TplDict[name]; ok {
		return tpl
	}
	tpl := template.New(name).Funcs(t.FunDict)
	tpl, err = tpl.ParseFiles(t.dir + path)
	if incl {
		tpl, err = tpl.ParseGlob(t.dir + "partials/*.html")
	}
	if err != nil {
		fmt.Println(err)
	}
	if err == nil {
		t.TplDict[name] = tpl
	}
	return tpl
}

func (t Theme) GetDir() string {
	return t.dir
}

func (t Theme) CopyAssets(dir string) (err error) {
	dir = strings.TrimSpace(dir)
	if dir == "" || dir == "." || dir == ".." {
		return
	}
	dir = t.dir + strings.Trim(dir, "/")
	err = utils.CopyDir(dir, t.PubDir)
	return
}

func (t Theme) CreateSidebar(path string, archDirs Categoris) (err error) {
	file := t.dir + "sidebar.html"
	if utils.PathExist(file) {
		ctx := Table{"ArchDirs": archDirs, "UrlPre": AddUrlPre("")}
		err = t.Render("sidebar", path, ctx)
	}
	return
}

func (t *Theme) Render(name, path string, cxt Table) (err error) {
	var file *os.File
	file, err = os.Create(t.PubDir + path)
	defer file.Close()
	if err != nil {
		return
	}
	file.Chmod(MODE_FILE)
	cxt["Path"] = path
	if !strings.HasSuffix(name, ".html") {
		path = name + ".html"
	}
	name = filepath.Base(path)
	tpl := t.GetOrCreate(name, path, true)
	err = tpl.ExecuteTemplate(file, tpl.Name(), cxt)
	return
}
