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

	fmt.Println("## Stackscript create")

	var ss *linodego.Stackscript
	for rev := 1; rev < 4; rev++ {
		fmt.Println("### Revision ", rev)
		if rev == 1 {
			stackscript := linodego.Stackscript{}.GetCreateOptions()
			stackscript.Description = "description for example stackscript " + time.Now().String()
			// stackscript.Images = make([]string, 2, 2)
			stackscript.Images = []string{"linode/debian9", "linode/ubuntu18.04"}
			stackscript.IsPublic = false
			stackscript.Label = "example stackscript " + time.Now().String()
			stackscript.RevNote = "revision " + strconv.Itoa(rev)
			stackscript.Script = "#!/bin/bash\n"
			ss, err = linodeClient.CreateStackscript(&stackscript)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			update := ss.GetUpdateOptions()
			update.RevNote = "revision " + strconv.Itoa(rev)
			update.Label = strconv.Itoa(rev) + " " + ss.Label
			update.Script += "echo " + strconv.Itoa(rev) + "\n"
			ss, err = linodeClient.UpdateStackscript(ss.ID, update)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	fmt.Println("### Delete ")
	err = linodeClient.DeleteStackscript(ss.ID)
	if err != nil {
		log.Fatal(err)
	}

	if spendMoney {
		linode, err = linodeClient.CreateInstance(&linodego.InstanceCreateOptions{Region: "us-central", Type: "g5-nanode-1"})
		if err != nil {
			log.Fatalln("* While creating instance: ", err)
		}

		fmt.Println("## Created Instance\n", linode)
		event, err := linodeClient.WaitForEventFinished(linode.ID, linodego.EntityLinode, linodego.ActionLinodeCreate, *linode.Created, 240)
		if err != nil {
			log.Fatalf("* Failed to wait for Linode %d to finish creation: %s", linode.ID, err)
		}
		if err := linodeClient.MarkEventRead(event); err != nil {
			log.Fatalln("* Failed to mark Linode create event seen", err)
		}
		disk, err := linodeClient.CreateInstanceDisk(linode.ID, linodego.InstanceDiskCreateOptions{Size: 50, Filesystem: "raw", Label: "linodego_disk"})
		if err != nil {
			log.Fatalln("* While creating disk:", err)
		}

		fmt.Println("### Created Disk\n", disk)
		event, err = linodeClient.WaitForEventFinished(linode.ID, linodego.EntityLinode, linodego.ActionDiskCreate, disk.Created, 240)
		if err := linodeClient.MarkEventRead(event); err != nil {
			log.Fatalln("* Failed to mark Disk create event seen", err)
		}

		createOpts := linodego.InstanceConfigCreateOptions{
			Devices: &linodego.InstanceConfigDeviceMap{
				SDA: &linodego.InstanceConfigDevice{DiskID: disk.ID},
			},
			Kernel:     "linode/direct-disk",
			Label:      "example config label",
			RunLevel:   "default",
			VirtMode:   "paravirt",
			Comments:   "example config comment",
			RootDevice: "/dev/sda",
			Helpers: &linodego.InstanceConfigHelpers{
				Network:    true,
				ModulesDep: false,
			},
		}
		config, err := linodeClient.CreateInstanceConfig(linode.ID, createOpts)
		if err != nil {
			log.Fatalln("* Failed to create Config", err)
		}
		fmt.Println("### Created Config:\n", config)
		updateOpts := linodego.InstanceConfigUpdateOptions{
			Comments: "updated example config comment",
		}
		config, err = linodeClient.UpdateInstanceConfig(linode.ID, config.ID, updateOpts)
		if err != nil {
			log.Fatalln("* Failed to update Config", err)
		}
		fmt.Println("### Updated Config:\n", config)

		err = linodeClient.DeleteInstanceConfig(linode.ID, config.ID)
		if err != nil {
			log.Fatalln("* Failed to delete Config", err)
		}
		fmt.Println("### Deleted Config")

		// @TODO it is not sufficient that a disk was created. Which disk was it?
		// Sounds like we'll need a WaitForEntityStatus function.
		if err != nil {
			log.Fatalf("* Failed to wait for Linode disk %d to finish creation: %s", disk.ID, err)
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
