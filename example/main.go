package main

import (
	"fmt"
	"log"
	"os"

	golinode "github.com/chiefy/go-linode"
)

func main() {
	// Trigger endpoints that accrue a balance
	apiKey, apiOk := os.LookupEnv("LINODE_API_KEY")
	var SpendMoney = true && apiOk

	// Demonstrate endpoints that don't require an account or token
	linodeClient, err := golinode.NewClient(nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	linodeClient.SetDebug(true)

	types, err := linodeClient.ListTypes()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", types)

	kernels, err := linodeClient.ListKernels(nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", kernels)

	if !apiOk || len(apiKey) == 0 {
		log.Fatal("Could not find LINODE_API_KEY, please assert it is set.")
		os.Exit(1)
	}

	// Demonstrate endpoints that require an access token
	linodeClient, err = golinode.NewClient(&apiKey, nil)
	if err != nil {
		log.Fatal(err)
	}
	linodeClient.SetDebug(true)

	var linode *golinode.LinodeInstance

	if SpendMoney {
		linode, err = linodeClient.CreateInstance(&golinode.InstanceCreateOptions{Region: "us-central", Type: "g5-nanode-1"})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%#v", linode)
	}

	linodes, err := linodeClient.ListInstances(nil)

	if len(linodes) == 0 {
		log.Printf("No Linodes to inspect.")
	} else {
		// This is redundantly used for illustrative purposes
		linode, err = linodeClient.GetInstance(linodes[0].ID)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%#v", linode)

		configs, err := linodeClient.ListInstanceConfigs(linode.ID)
		if err != nil {
			log.Fatal(err)
		} else if len(configs) > 0 {
			config, err := linodeClient.GetInstanceConfig(linode.ID, configs[0].ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("First Config: %#v", config)
		}

		disks, err := linodeClient.ListInstanceDisks(linode.ID)
		if err != nil {
			log.Fatal(err)
		} else if len(disks) > 0 {
			disk, err := linodeClient.GetInstanceDisk(linode.ID, disks[0].ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("First Disk: %#v", disk)
		}
	}
}
