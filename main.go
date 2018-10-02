package main

import (
	"flag"
	"fmt"
	"plugin"

	"fountain/article"
	"fountain/utils"
	"gopkg.in/russross/blackfriday.v2"
)

const VERSION = "0.39.7"

var (
	serve   bool   //运行WEB服务
	port    uint   //服务端口
	root    string //博客根目录
	clean   bool   //清理旧输出
	verbose bool   //输出详情
)

func init() {
	flag.BoolVar(&serve, "s", false, "运行WEB服务")
	flag.UintVar(&port, "p", 0, "服务端口")
	flag.StringVar(&root, "r", "", "博客根目录")
	flag.BoolVar(&clean, "c", false, "清理旧输出")
	flag.BoolVar(&verbose, "v", false, "输出详情")
	flag.Parse()
}

func main() {
	site := article.NewWebsite(root)
	site.LoadConfig("config.yml")
	site.Convert = Convert
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
			pubDir := site.RootDir + site.Conf.Public
			if len(site.Conf.Public) >= 3 {
				utils.CleanDir(pubDir)
			}
		}
		fmt.Println("Build ...")
		site.InitTheme()
		site.BuildFiles()
	}
}

func Convert(source []byte, format string) []byte {
	if format == "" {
		return source
	} else if format == "markdown" {
		return blackfriday.Run(source)
	}
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
