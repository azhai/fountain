package article

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"fountain/utils"
)

type Theme struct {
	tplDir   string //结尾有斜杠
	PubDir   string //结尾有斜杠
	FunDict  template.FuncMap
	TplDict  map[string]*template.Template
	WithSide bool
}

func NewTheme(tplDir string) *Theme {
	return &Theme{
		tplDir:  tplDir,
		FunDict: make(map[string]interface{}),
		TplDict: make(map[string]*template.Template),
	}
}

func (this *Theme) GetOrCreate(name string, incl bool) *template.Template {
	var err error
	if tpl, ok := this.TplDict[name]; ok {
		return tpl
	}
	tpl := template.New(name).Funcs(this.FunDict)
	tpl, err = tpl.ParseFiles(this.tplDir + name + ".html")
	if incl {
		partials, _ := filepath.Glob(this.tplDir + "partials/*.html")
		tpl, err = tpl.ParseFiles(partials...)
	}
	if err != nil {
		fmt.Println(err)
	}
	if err == nil {
		this.TplDict[name] = tpl
	}
	return tpl
}

func (this *Theme) AddUrlPre(dir string) string {
	prefix := "."
	if len(dir) > 0 && dir != "." {
		times := strings.Count(dir, "/") + 1
		prefix = strings.Repeat("../", times)
		prefix = prefix[:len(prefix)-1]
	}
	return prefix
}

func (this *Theme) CopyAssets(dirs ...string) (err error) {
	for _, dir := range dirs {
		if dir == "" {
			continue
		}
		dir = this.tplDir + strings.Trim(dir, "/")
		err = utils.CopyDir(dir, this.PubDir)
	}
	return
}

func (this *Theme) CreateSidebar(path string, archDirs Categoris) (err error) {
	file := this.tplDir + "sidebar.html"
	if utils.PathExist(file) {
		fmt.Println("√", path)
		ctx := Table{"ArchDirs": archDirs, "UrlPre": this.AddUrlPre("")}
		err = this.Render("sidebar", path, ctx)
	}
	return
}

func (this *Theme) Render(name, path string, cxt Table) (err error) {
	var file *os.File
	file, err = os.Create(this.PubDir + path)
	defer file.Close()
	if err != nil {
		return
	}
	file.Chmod(MODE_FILE)
	cxt["Path"] = path
	tpl := this.GetOrCreate(name, true)
	err = tpl.ExecuteTemplate(file, name+".html", cxt)
	return
}
