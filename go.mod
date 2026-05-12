module github.com/linode/linodego

require (
	github.com/go-resty/resty/v2 v2.17.2
	github.com/google/go-cmp v0.7.0
	github.com/google/go-querystring v1.2.0
	github.com/jarcoal/httpmock v1.4.1
	golang.org/x/net v0.54.0
	golang.org/x/oauth2 v0.36.0
	golang.org/x/text v0.37.0
	gopkg.in/ini.v1 v1.67.2
)

require (
	github.com/kr/pretty v0.3.1 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	golang.org/x/time v0.14.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/stretchr/testify v1.11.1
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

go 1.26.0

retract v1.0.0 // Accidental branch push
