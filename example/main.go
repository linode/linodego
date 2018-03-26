package main

import (
	"fmt"
	"log"
	"os"

	golinode "github.com/chiefy/go-linode"
)

func main() {
	apiKey, ok := os.LookupEnv("LINODE_API_KEY")
	if !ok {
		log.Fatal("Could not find LINODE_API_KEY, please assert it is set.")
	}
	linodeClient, err := golinode.NewClient(&apiKey, nil)
	if err != nil {
		log.Fatal(err)
	}
	linodeClient.SetDebug(true)

	types, err := linodeClient.ListTypes()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", types)

	kernels, err := linodeClient.ListKernels()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", kernels)

	linodes, err := linodeClient.ListInstances()

	if len(linodes) == 0 {
		log.Printf("No Linodes to inspect.")
	} else {
		// This is redundantly used for illustrative purposes
		linode, err := linodeClient.GetInstance(linodes[0].ID)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%v", linode)

		configs, err := linodeClient.ListInstanceConfigs(linode.ID)
		if err != nil {
			log.Fatal(err)
		} else {
			config, err := linodeClient.GetInstanceConfig(linode.ID, configs[0].ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("First Config: %v", config)
		}

		disks, err := linodeClient.ListInstanceDisks(linode.ID)
		if err != nil {
			log.Fatal(err)
		} else {
			disk, err := linodeClient.GetInstanceDisk(linode.ID, disks[0].ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("First Disk: %v", disk)
		}
	}
}
