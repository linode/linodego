package main

import (
	"fmt"
	"log"
	"os"

	golinode "github.com/chiefy/go-linode"
)

var linodeClient = golinode.NewClient(nil, nil)
var spendMoney = false

func main() {
	// Trigger endpoints that accrue a balance
	apiToken, apiOk := os.LookupEnv("LINODE_TOKEN")
	spendMoney = true && apiOk

	var err error
	if err != nil {
		log.Fatal(err)
	}
	linodeClient.SetDebug(false)

	if !apiOk || len(apiToken) == 0 {
		log.Fatal("Could not find LINODE_TOKEN, please assert it is set.")
	}

	// Demonstrate endpoints that require an access token
	linodeClient = golinode.NewClient(&apiToken, nil)
	if err != nil {
		log.Fatal(err)
	}

	moreExamples_authenticated()
}

func moreExamples_authenticated() {
	var linode *golinode.Instance

	linode, err := linodeClient.GetInstance(1231)
	fmt.Printf("%#v", linode)

	if spendMoney {
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

		configs, err := linodeClient.ListInstanceConfigs(linode.ID, nil)
		if err != nil {
			log.Fatal(err)
		} else if len(configs) > 0 {
			config, err := linodeClient.GetInstanceConfig(linode.ID, configs[0].ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("First Config: %#v", config)
		}

		disks, err := linodeClient.ListInstanceDisks(linode.ID, nil)
		if err != nil {
			log.Fatal(err)
		} else if len(disks) > 0 {
			disk, err := linodeClient.GetInstanceDisk(linode.ID, disks[0].ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("First Disk: %#v", disk)
		}

		volumes, err := linodeClient.ListInstanceVolumes(linode.ID, nil)
		if err != nil {
			log.Fatal(err)
		} else if len(volumes) > 0 {
			volume, err := linodeClient.GetInstanceVolume(linode.ID, volumes[0].ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("First Volume: %#v", volume)
		}

		stackscripts, err := linodeClient.ListStackscripts(&golinode.ListOptions{Filter: "{\"mine\":true}"})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%+v", stackscripts)
	}
}
