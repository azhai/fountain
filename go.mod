module fountain

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20180910181607-0e37d006457b
	golang.org/x/net => github.com/golang/net v0.0.0-20180925072008-f04abc6bdfa7
	golang.org/x/sync => github.com/golang/sync v0.0.0-20180314180146-1d60e4601c6f
	golang.org/x/sys => github.com/golang/sys v0.0.0-20180925112736-b09afc3d579e
	golang.org/x/text => github.com/golang/text v0.3.0
	gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
)

require (
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/shurcooL/sanitized_anchor_name v0.0.0-20170918181015-86672fcb3f95 // indirect
	gopkg.in/russross/blackfriday.v2 v2.0.1
	gopkg.in/yaml.v2 v2.2.1
)
