package main

import "gopkg.in/russross/blackfriday.v2"

// 将Markdown转为html
func Convert(source []byte) []byte {
	return blackfriday.Run(source)
}
