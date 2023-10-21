module github.com/linode/linodego

require (
	github.com/go-resty/resty/v2 v2.10.0
	github.com/google/go-cmp v0.6.0
	golang.org/x/net v0.17.0
	golang.org/x/text v0.13.0
	gopkg.in/ini.v1 v1.66.6
)

require github.com/stretchr/testify v1.8.4 // indirect

go 1.20

retract v1.0.0 // Accidental branch push
