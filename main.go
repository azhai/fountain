package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"runtime"
	"strings"

	"fountain/article"
	"fountain/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/k0kubun/pp"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/toc"
)

const VERSION = "0.6.1"

var (
	serve   bool   // 运行WEB服务
	port    uint   // 服务端口
	root    string // 博客根目录
	theme   string // 皮肤主题
	clean   bool   // 清理旧输出
	verbose bool   // 输出详情

	md = goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)
)

func init() {
	flag.BoolVar(&serve, "s", false, "运行WEB服务")
	flag.UintVar(&port, "p", 0, "服务端口")
	flag.StringVar(&root, "r", "", "博客根目录")
	flag.StringVar(&theme, "t", "", "皮肤主题")
	flag.BoolVar(&clean, "c", false, "清理旧输出")
	flag.BoolVar(&verbose, "v", false, "输出详情")
	flag.Usage = usage
	flag.Parse()
}

func usage() {
	desc := `fountain version: v%s
Usage: fountain [-r root] [-s] [-p port] [-t theme] [-c] [-v]

Options:
`
	fmt.Fprintf(os.Stderr, desc, VERSION)
	flag.PrintDefaults()
}

func main() {
	if runtime.GOOS == "windows" {
		name := "Fountain"
		desc := "Fountain Static Blog Server"
		utils.WinMain(name, desc, run)
	} else {
		run()
	}
}

func run() {
	site := article.NewWebsite(root)
	site.LoadConfig("config.yaml")
	if verbose {
		pp.Println(site.Conf)
	}
	if theme != "" {
		site.Conf.Theme = theme
	}
	site.Convert = func(source []byte, format string) []byte {
		if format == "" {
			return source
		}
		format = strings.ToLower(format)
		if format == "markdown" {
			return MarkdownConvert(source)
		}
		return PluginConvert(source, format)
	}
	site.Debug = func(data ...any) {
		if verbose {
			fmt.Println(data...)
		}
	}

	if clean {
		fmt.Println("Clean ...")
		if len(site.Conf.Public) >= 3 {
			pubDir := site.Root + site.Conf.Public
			utils.CleanDir(pubDir)
		}
	}
	fmt.Println("Build ...")
	site.InitTheme()
	site.BuildFiles()

	if !serve {
		return
	}
	if port == 0 || port > 65535 {
		port = site.Conf.Port
	}
	fmt.Printf("Server at :%d\n", port)

	pudDir := filepath.Join(root, "public")
	app := fiber.New()
	app.Use("/", static.New(pudDir))
	app.Listen(fmt.Sprintf("0.0.0.0:%d", port))
}

func MarkdownConvert(source []byte) []byte {
	source = append([]byte(article.SEP_OUTLINE+"\n"), source...)
	doc := md.Parser().Parse(text.NewReader(source))
	tree, err := toc.Inspect(doc, source, toc.Compact(true))
	if err == nil {
		if list := toc.RenderList(tree); list != nil {
			doc.InsertBefore(doc, doc.FirstChild(), list)
		}
	}
	var buf bytes.Buffer
	if err = md.Renderer().Render(&buf, source, doc); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func PluginConvert(source []byte, format string) []byte {
	plug, err := plugin.Open("./" + format + ".so")
	if err != nil {
		panic(err)
	}
	symb, err := plug.Lookup("Convert")
	if err != nil {
		panic(err)
	}
	conv := symb.(func(source []byte) []byte)
	return conv(source)
}
