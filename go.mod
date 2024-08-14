module github.com/linode/linodego

require (
	github.com/go-resty/resty/v2 v2.13.1
	github.com/google/go-cmp v0.6.0
	github.com/jarcoal/httpmock v1.3.1
	golang.org/x/net v0.28.0
	golang.org/x/oauth2 v0.22.0
	golang.org/x/text v0.17.0
	gopkg.in/ini.v1 v1.66.6
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.9.0
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

go 1.21

retract v1.0.0 // Accidental branch push
