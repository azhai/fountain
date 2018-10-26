module fountain

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20181025213731-e84da0312774
	golang.org/x/net => github.com/golang/net v0.0.0-20181023162649-9b4f9f5ad519
	golang.org/x/sync => github.com/golang/sync v0.0.0-20180314180146-1d60e4601c6f
	golang.org/x/sys => github.com/golang/sys v0.0.0-20181025063200-d989b31c8746
	golang.org/x/text => github.com/golang/text v0.3.0
	gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
)

require (
	github.com/kardianos/osext v0.0.0-20170510131534-ae77be60afb1 // indirect
	github.com/kardianos/service v0.0.0-20180910224244-b1866cf76903
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/shurcooL/sanitized_anchor_name v0.0.0-20170918181015-86672fcb3f95 // indirect
	gopkg.in/russross/blackfriday.v2 v2.0.1
	gopkg.in/yaml.v2 v2.2.1
)
