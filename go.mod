module fountain

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190426145343-a29dc8fdc734
	golang.org/x/net => github.com/golang/net v0.0.0-20190502183928-7f726cade0ab
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190502175342-a43fa875dd82
	golang.org/x/text => github.com/golang/text v0.3.2
	gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
)

require (
	github.com/kardianos/service v1.0.0
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20190429190828-d89cdac9e872 // indirect
	gopkg.in/russross/blackfriday.v2 v2.0.1
	gopkg.in/yaml.v2 v2.2.2
)
