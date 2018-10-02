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
	Dir     string //结尾有斜杠
	OutDir  string //结尾有斜杠
	FuncMap template.FuncMap
	Tpls    map[string]*template.Template
}

func NewTheme() *Theme {
	return &Theme{
		FuncMap: make(map[string]interface{}),
		Tpls:    make(map[string]*template.Template),
	}
}

func (this *Theme) GetOrCreate(name string) *template.Template {
	var err error
	if tpl, ok := this.Tpls[name]; ok {
		return tpl
	}
	tpl := template.New(name).Funcs(this.FuncMap)
	partials, _ := filepath.Glob(this.Dir + "partials/*.html")
	tpl, err = tpl.ParseFiles(this.Dir + name + ".html")
	tpl, err = tpl.ParseFiles(partials...)
	if err != nil {
		fmt.Println(err)
	}
	if err == nil {
		this.Tpls[name] = tpl
	}
	return tpl
}

func (this *Theme) AddUrlPre(dir string) string {
	prefix := ""
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
		dir = this.Dir + strings.Trim(dir, "/")
		err = utils.CopyDir(dir, this.OutDir)
	}
	return
}

func (this *Theme) Render(name, path string, cxt Table) (err error) {
	var file *os.File
	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, MODE_DIR)
	if err != nil {
		return
	}
	cxt["UrlPre"] = this.AddUrlPre(dir)
	file, err = os.Create(this.OutDir + path)
	defer file.Close()
	if err != nil {
		return
	}
	file.Chmod(MODE_FILE)
	tpl := this.GetOrCreate(name)
	err = tpl.ExecuteTemplate(file, name+".html", cxt)
	return
}
