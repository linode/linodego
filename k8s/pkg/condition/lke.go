package condition

import (
	"context"
	"fmt"

	"github.com/linode/linodego"
	"github.com/linode/linodego/k8s"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterHasReadyNode is a ClusterConditionFunc which polls for at least one node to have the
// condition NodeReady=True.
func ClusterHasReadyNode(ctx context.Context, options linodego.ClusterConditionOptions) (bool, error) {
	clientset, err := k8s.BuildClientsetFromConfig(options.LKEClusterKubeconfig, options.TransportWrapper)
	if err != nil {
		return false, err
	}

	nodes, err := clientset.CoreV1().Nodes().List(ctx, v1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to get nodes for cluster: %w", err)
	}

	for _, node := range nodes.Items {
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue {
				return true, nil
			}
		}
	}
	return false, nil
}

// ClusterNodesReady is a ClusterConditionFunc which polls for all nodes to have the
// condition NodeReady=True.
func ClusterNodesReady(ctx context.Context, options linodego.ClusterConditionOptions) (bool, error) {
	clientset, err := k8s.BuildClientsetFromConfig(options.LKEClusterKubeconfig, options.TransportWrapper)
	if err != nil {
		return false, err
	}

	nodes, err := clientset.CoreV1().Nodes().List(ctx, v1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to get nodes for cluster: %w", err)
	}

	for _, node := range nodes.Items {
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady && condition.Status != corev1.ConditionTrue {
				return false, nil
			}
		}
	}
	return true, nil
}

// WaitForLKEClusterReady polls with a given timeout for the LKE Cluster's api-server
// to be healthy and for the cluster to have at least one node with the NodeReady
// condition true.
func WaitForLKEClusterReady(ctx context.Context, client linodego.Client, clusterID int, options linodego.LKEClusterPollOptions) error {
	return client.WaitForLKEClusterConditions(ctx, clusterID, options, ClusterHasReadyNode)
}

// WaitForLKEClusterAndNodesReady polls with a given timeout for the LKE
// Cluster's api-server to be healthy and for all cluster nodes to have the
// NodeReady condition true.
func WaitForLKEClusterAndNodesReady(ctx context.Context, client linodego.Client, clusterID int, options linodego.LKEClusterPollOptions) error {
	return client.WaitForLKEClusterConditions(ctx, clusterID, options, ClusterNodesReady)
}
