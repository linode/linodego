package k8s

import (
	"encoding/base64"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/transport"

	"github.com/linode/linodego"
)

// NewClientsetFromBytes builds a Clientset from a given Kubeconfig.
//
// Takes an optional transport.WrapperFunc to add request/response middleware to
// api-server requests.
func BuildClientsetFromConfig(
	lkeKubeconfig *linodego.LKEClusterKubeconfig,
	transportWrapper transport.WrapperFunc,
) (kubernetes.Interface, error) {
	kubeConfigBytes, err := base64.StdEncoding.DecodeString(lkeKubeconfig.KubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to decode kubeconfig: %w", err)
	}

	restClientConfig, err := clientcmd.RESTConfigFromKubeConfig(kubeConfigBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse LKE cluster kubeconfig: %w", err)
	}

	if transportWrapper != nil {
		restClientConfig.Wrap(transportWrapper)
	}

	clientset, err := kubernetes.NewForConfig(restClientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build k8s client from LKE cluster kubeconfig: %w", err)
	}
	return clientset, nil
}
