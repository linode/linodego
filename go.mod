module github.com/linode/linodego

require (
	github.com/go-resty/resty/v2 v2.17.1
	github.com/google/go-cmp v0.7.0
	github.com/google/go-querystring v1.2.0
	github.com/jarcoal/httpmock v1.4.1
	golang.org/x/net v0.49.0
	golang.org/x/oauth2 v0.34.0
	golang.org/x/text v0.33.0
	gopkg.in/ini.v1 v1.67.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.11.1
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

go 1.24.0

toolchain go1.25.1

retract v1.0.0 // Accidental branch push
