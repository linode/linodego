package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/chiefy/linodego"
)

var linodeClient = linodego.NewClient(nil, nil)
var spendMoney = false

func main() {
	// Trigger endpoints that accrue a balance
	apiToken, apiOk := os.LookupEnv("LINODE_TOKEN")
	spendMoney = spendMoney && apiOk

	var err error
	if err != nil {
		log.Fatal(err)
	}

	if !apiOk || len(apiToken) == 0 {
		log.Fatal("Could not find LINODE_TOKEN, please verify that it is set.")
	}

	// Demonstrate endpoints that require an access token
	linodeClient = linodego.NewClient(&apiToken, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Wether or not we will walk example endpoints that cost money
	if envSpend, spendSet := os.LookupEnv("LINODE_SPEND"); apiOk && spendSet {
		if apiSpend, err := strconv.Atoi(envSpend); err == nil {
			log.Println("LINODE_SPEND being set to", apiSpend > 0)
			spendMoney = apiSpend > 0
		} else {
			log.Fatalln("LINODE_SPEND should be an integer, 0 or 1")
		}
	}

	// Wether or not we will enable Resty debugging output
	if envDebug, apiOk := os.LookupEnv("LINODE_DEBUG"); apiOk {
		if apiDebug, err := strconv.Atoi(envDebug); err == nil {
			log.Println("LINODE_DEBUG being set to", apiDebug > 0)
			linodeClient.SetDebug(apiDebug > 0)
		} else {
			log.Fatalln("LINODE_DEBUG should be an integer, 0 or 1")
		}
	}

	moreExamples_authenticated()
}

func moreExamples_authenticated() {
	var linode *linodego.Instance

	linode, err := linodeClient.GetInstance(1231)
	fmt.Println("## Instance request with Invalid ID")
	fmt.Println("### Linode\n", linode, "\n### Error\n", err)

	if spendMoney {
		linode, err = linodeClient.CreateInstance(&linodego.InstanceCreateOptions{Region: "us-central", Type: "g5-nanode-1"})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("## Created Instance\n", linode)
		fmt.Println("### Creating Disk")
		for i := 0; i < 30; i++ {
			disk, err := linodeClient.CreateInstanceDisk(linode.ID, linodego.InstanceDiskCreateOptions{Size: 50, Filesystem: "raw", Label: "linodego_disk"})
			if err != nil {
				linodeErr, ok := err.(*linodego.Error)
				if ok && linodeErr.Code >= 400 {
					log.Printf("- Waiting 1s for disk (got %s) Pass %d", linodeErr.Error(), i)
					time.Sleep(time.Second)
					continue
				}
				log.Fatalln("* While creating disk:", err)
			}
			log.Println(disk)
			err = linodeClient.WaitForEventFinished(disk.ID, linodego.EntityDisk, linodego.ActionDiskCreate, disk.Created, 5)
			if err != nil {
				log.Fatalf("* Failed to wait for Linode disk %d to finish creation: %s", disk.ID, err)
			}
			break
		}
	}

	linodes, err := linodeClient.ListInstances(nil)
	fmt.Println("## List Instances")

	if len(linodes) == 0 {
		log.Println("No Linodes to inspect.")
	} else {
		// This is redundantly used for illustrative purposes
		linode, err = linodeClient.GetInstance(linodes[0].ID)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("## First Linode\n", linode)

		configs, err := linodeClient.ListInstanceConfigs(linode.ID, nil)
		if err != nil {
			log.Fatal(err)
		} else if len(configs) > 0 {
			config, err := linodeClient.GetInstanceConfig(linode.ID, configs[0].ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("### First Config:\n", config)
		} else {
			fmt.Println("### No Configs")
		}

		disks, err := linodeClient.ListInstanceDisks(linode.ID, nil)
		if err != nil {
			log.Fatal(err)
		} else if len(disks) > 0 {
			disk, err := linodeClient.GetInstanceDisk(linode.ID, disks[0].ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("### First Disk\n", disk)
		} else {
			fmt.Println("### No Disks")
		}

		backups, err := linodeClient.GetInstanceBackups(linode.ID)
		if err != nil {
			log.Fatal(err)
		}
		if len(backups.Automatic) > 0 {
			fmt.Println("### First Auto Backup\n", backups.Automatic[0])
		} else {
			fmt.Println("### No Auto Backups")
		}
		fmt.Println("### Snapshots\n", backups.Snapshot)
		if backups.Snapshot.Current != nil {
			// snapshot fetched will be exactly the same as backups.Snapshot.Current
			// just being redundant for illustrative purposes
			if snapshot, err := linodeClient.GetInstanceSnapshot(linode.ID, backups.Snapshot.Current.ID); err == nil {
				fmt.Println("#### Current\n", snapshot)
			} else {
				fmt.Println("#### No Current Snapshot\n", err)
			}
		} else {
			fmt.Println("### No Current Snapshot")
		}

		volumes, err := linodeClient.ListInstanceVolumes(linode.ID, nil)
		if err != nil {
			log.Fatal(err)
		} else if len(volumes) > 0 {
			volume, err := linodeClient.GetInstanceVolume(linode.ID, volumes[0].ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("### First Volume\n", volume)
		} else {
			fmt.Println("### No Volumes")
		}

		stackscripts, err := linodeClient.ListStackscripts(&linodego.ListOptions{Filter: "{\"mine\":true}"})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("## Your Stackscripts\n", stackscripts)
	}
}
