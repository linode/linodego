module github.com/linode/linodego

require (
	github.com/go-resty/resty/v2 v2.7.0
	github.com/google/go-cmp v0.5.7
	github.com/stretchr/testify v1.7.1 // indirect
	gopkg.in/ini.v1 v1.66.6
)

require (
	golang.org/x/net v0.14.0
	golang.org/x/text v0.12.0 // indirect
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

go 1.20

retract v1.0.0 // Accidental branch push
