module github.com/linode/linodego

require (
	github.com/go-resty/resty/v2 v2.1.1-0.20191201195748-d7b97669fe48
	github.com/google/go-cmp v0.5.7
	github.com/stretchr/testify v1.7.1 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	gopkg.in/ini.v1 v1.66.4
)

require (
	github.com/golang/protobuf v1.2.0 // indirect
	golang.org/x/net v0.0.0-20190628185345-da137c7871d7 // indirect
	google.golang.org/appengine v1.4.0 // indirect
)

go 1.18

retract v1.0.0 // Accidental branch push
