module github.com/linode/linodego

require (
	github.com/go-resty/resty/v2 v2.12.0
	github.com/google/go-cmp v0.6.0
	github.com/jarcoal/httpmock v1.3.1
	golang.org/x/net v0.24.0
	golang.org/x/oauth2 v0.20.0
	golang.org/x/text v0.14.0
	gopkg.in/ini.v1 v1.66.6
)

require github.com/stretchr/testify v1.9.0 // indirect

go 1.21

retract v1.0.0 // Accidental branch push
