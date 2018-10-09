package main

import bf2 "gopkg.in/russross/blackfriday.v2"

func WithOptions(flags bf2.HTMLFlags) bf2.Option {
	params := bf2.HTMLRendererParameters{Flags: flags}
	renderer := bf2.NewHTMLRenderer(params)
	return bf2.WithRenderer(renderer)
}

// 将Markdown转为html
func Convert(source []byte) []byte {
	flags := bf2.CommonHTMLFlags | bf2.TOC
	return bf2.Run(source, WithOptions(flags))
}
