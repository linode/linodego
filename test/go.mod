module test

go 1.15

require (
	github.com/dnaeon/go-vcr v1.1.0
	github.com/google/go-cmp v0.5.7
	github.com/linode/linodego v0.20.1
	github.com/linode/linodego/k8s v0.0.0-00010101000000-000000000000
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f
	k8s.io/client-go v0.23.4
)

replace github.com/linode/linodego => ../

replace github.com/linode/linodego/k8s => ../k8s
