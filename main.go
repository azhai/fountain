package main

import (
	"flag"
	"fmt"
	"plugin"

	"fountain/article"
)

const VERSION = "0.39.4"

var (
	build bool   //构建html
	debug bool   //输出详情
	port  uint   //服务端口
	root  string //博客根目录
)

func init() {
	flag.BoolVar(&build, "build", false, "构建html")
	flag.BoolVar(&debug, "debug", false, "输出详情")
	flag.UintVar(&port, "port", 8080, "服务端口")
	flag.StringVar(&root, "root", "", "博客根目录")
	flag.Parse()
}

func main() {
	site := article.NewWebsite(root)
	site.LoadConfig("config.yml")
	site.Convert = Convert
	site.Debug = func(data ...interface{}) {
		if debug {
			fmt.Println(data...)
		}
	}
	if build {
		fmt.Println("Build ...")
		site.InitTheme()
		site.BuildFiles()
	} else {
		fmt.Printf("Server at :%d\n", port)
	}
}

func Convert(source []byte, format string) []byte {
	if format == "" {
		return source
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
