module test

go 1.13

require (
	github.com/dnaeon/go-vcr v1.0.1
	github.com/google/go-cmp v0.5.2
	github.com/linode/linodego v0.20.1
	github.com/linode/linodego/k8s v0.0.0-00010101000000-000000000000
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	gopkg.in/yaml.v2 v2.3.0 // indirect
	k8s.io/client-go v0.18.3
)

replace github.com/linode/linodego => ../

replace github.com/linode/linodego/k8s => ../k8s
