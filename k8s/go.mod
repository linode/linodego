module github.com/linode/linodego/k8s

require (
	github.com/linode/linodego v0.20.1
	k8s.io/api v0.21.1
	k8s.io/apimachinery v0.21.1
	k8s.io/client-go v0.21.1
)

replace github.com/linode/linodego => ../

go 1.15
