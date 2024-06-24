module fountain

go 1.22

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.1.0

require (
	github.com/BurntSushi/toml v1.4.0
	github.com/kardianos/service v1.2.2
	gopkg.in/russross/blackfriday.v2 v2.1.0
)

require golang.org/x/sys v0.21.0 // indirect
