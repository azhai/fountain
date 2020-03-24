package main

import (
	"flag"
	"fmt"
	"os"
	"plugin"
	"runtime"

	"fountain/article"
	"fountain/utils"
	bf2 "gopkg.in/russross/blackfriday.v2"
)

const VERSION = "0.41.2"

var (
	serve   bool   //运行WEB服务
	port    uint   //服务端口
	root    string //博客根目录
	theme   string //皮肤主题
	clean   bool   //清理旧输出
	verbose bool   //输出详情
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
Usage: fountain [-t theme] [-r root] [-p port] [-scv]

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
	site.LoadConfig("config.yml")
	if theme != "" {
		site.Conf.Theme = theme
	}
	site.Convert = func(source []byte, format string) []byte {
		if format == "" {
			return source
		} else if format == "markdown" {
			flags := bf2.CommonHTMLFlags
			if site.Conf.Theme != "night" {
				flags = flags | bf2.TOC
			}
			return bf2.Run(source, WithOptions(flags))
		}
		return PluginConvert(source, format)
	}
	site.Debug = func(data ...interface{}) {
		if verbose {
			fmt.Println(data...)
		}
	}
	if serve || port > 0 {
		if port == 0 || port > 65535 {
			port = site.Conf.Port
		}
		fmt.Printf("Server at :%d\n", port)
	} else {
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
	}
}

func WithOptions(flags bf2.HTMLFlags) bf2.Option {
	params := bf2.HTMLRendererParameters{Flags: flags}
	renderer := bf2.NewHTMLRenderer(params)
	return bf2.WithRenderer(renderer)
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
