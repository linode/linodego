package main

import (
	"context"
	"fmt"
	"os"

	"github.com/linode/linodego"
)

func main() {
	client, err := linodego.NewClientFromEnv(nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create Encrypted VM
	encryptedLinode, err := client.CreateInstance(context.TODO(), linodego.InstanceCreateOptions{
		Region: "us-east",
		Type:   "g6-standard-1",
		Label:  "encrypted",
		// SET_ME_SECURELY
		RootPass:       "",
		Image:          "linode/debian12-kube-v1.28.3",
		DiskEncryption: &linodego.InstanceDiskEncryptionEnabled,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%v\n", encryptedLinode)

	unencryptedLinode, err := client.CreateInstance(context.TODO(), linodego.InstanceCreateOptions{
		Region: "us-east",
		Type:   "g6-standard-1",
		Label:  "not",
		// SET_ME_SECURELY
		RootPass:       "",
		Image:          "linode/debian12-kube-v1.28.3",
		Tags:           []string{"not"},
		DiskEncryption: &linodego.InstanceDiskEncryptionDisabled,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%v\n", unencryptedLinode)

	c, err := client.CreateLKECluster(context.TODO(), linodego.LKEClusterCreateOptions{
		NodePools: []linodego.LKENodePoolCreateOptions{
			{
				Count:          1,
				Type:           "g6-standard-1",
				Tags:           []string{"linodego"},
				DiskEncryption: &linodego.InstanceDiskEncryptionEnabled,
			},
		},
		Label:      "enabled2",
		Region:     "us-east",
		K8sVersion: "1.28",
		Tags:       []string{"bingo"},
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	nodePools, err := client.ListLKENodePools(context.TODO(), c.ID, &linodego.ListOptions{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%v", nodePools)
}
