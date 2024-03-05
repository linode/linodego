module github.com/linode/linodego

require (
	github.com/go-resty/resty/v2 v2.11.0
	github.com/google/go-cmp v0.6.0
	github.com/jarcoal/httpmock v1.3.1
	golang.org/x/net v0.22.0
	golang.org/x/oauth2 v0.18.0
	golang.org/x/text v0.14.0
	gopkg.in/ini.v1 v1.66.6
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)

go 1.20

retract v1.0.0 // Accidental branch push
