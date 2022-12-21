module github.com/linode/linodego

require (
	github.com/go-resty/resty/v2 v2.1.1-0.20191201195748-d7b97669fe48
	github.com/google/go-cmp v0.5.7
	github.com/stretchr/testify v1.7.1 // indirect
	gopkg.in/ini.v1 v1.66.6
)

require (
	golang.org/x/net v0.0.0-20190628185345-da137c7871d7
	golang.org/x/text v0.3.0 // indirect
)

go 1.18

retract v1.0.0 // Accidental branch push
