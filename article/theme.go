package article

import (
	"fmt"
	"os"
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

func (this *Theme) GetOrCreate(name string) *template.Template {
	var err error
	if tpl, ok := this.Tpls[name]; ok {
		return tpl
	}
	tpl := template.New(name).Funcs(this.FuncMap)
	tpl, err = tpl.ParseFiles(this.Dir+name+".html",
		this.Dir+"_head.html", this.Dir+"_header.html")
	if err != nil {
		fmt.Println(err)
	}
	if err == nil {
		this.Tpls[name] = tpl
	}
	return tpl
}

func (this *Theme) Render(name, path string, context interface{}) error {
	file, err := os.Create(this.OutDir + path)
	defer file.Close()
	if err != nil {
		return err
	}
	file.Chmod(MODE_FILE)
	tpl := this.GetOrCreate(name)
	return tpl.ExecuteTemplate(file, name+".html", context)
}

func (this *Theme) RenderArticle(path string, blog *Article) error {
	return this.Render("article", path, blog)
}
