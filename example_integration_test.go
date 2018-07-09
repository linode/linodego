package linodego_test

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/chiefy/linodego"
)

var spendMoney = false

func init() {
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

	// Wether or not we will walk example endpoints that cost money
	if envSpend, spendSet := os.LookupEnv("LINODE_SPEND"); apiOk && spendSet {
		if apiSpend, err := strconv.Atoi(envSpend); err == nil {
			log.Println("LINODE_SPEND being set to", apiSpend > 0)
			spendMoney = apiSpend > 0
		} else {
			log.Fatalln("LINODE_SPEND should be an integer, 0 or 1")
		}
	}
}

func ExampleGetAccount() {
	// Example readers, Ignore this bit of setup code needed to record test fixtures
	linodeClient, teardown := createTestClient(nil, "fixtures/ExampleGetAccount")
	defer teardown()

	account, err := linodeClient.GetAccount()
	if err != nil {
		log.Fatalln("* While getting account: ", err)
	}
	fmt.Println("Account email has @:", strings.Contains(account.Email, "@"))

	// Output:
	// Account email has @: true
}

func Example() {
	// Example readers, Ignore this bit of setup code needed to record test fixtures
	linodeClient, teardown := createTestClient(nil, "fixtures/Example")
	defer teardown()

	var linode *linodego.Instance
	linode, err := linodeClient.GetInstance(1231)
	fmt.Println("## Instance request with Invalid ID")
	fmt.Println("### Linode:", linode)
	fmt.Println("### Error:", err)

	if spendMoney {
		linode, err = linodeClient.CreateInstance(&linodego.InstanceCreateOptions{Region: "us-central", Type: "g5-nanode-1"})
		if err != nil {
			log.Fatalln("* While creating instance: ", err)
		}
		linode, err = linodeClient.UpdateInstance(linode.ID, &linodego.InstanceUpdateOptions{Label: linode.Label + "-renamed"})
		if err != nil {
			log.Fatalln("* While renaming instance: ", err)
		}
		fmt.Println("## Created Instance")
		event, err := linodeClient.WaitForEventFinished(linode.ID, linodego.EntityLinode, linodego.ActionLinodeCreate, *linode.Created, 240)
		if err != nil {
			log.Fatalf("* Failed to wait for Linode %d to finish creation: %s", linode.ID, err)
		}
		if err := linodeClient.MarkEventRead(event); err != nil {
			log.Fatalln("* Failed to mark Linode create event seen", err)
		}

		diskSwap, err := linodeClient.CreateInstanceDisk(linode.ID, linodego.InstanceDiskCreateOptions{Size: 50, Filesystem: "swap", Label: "linodego_swap"})
		if err != nil {
			log.Fatalln("* While creating swap disk:", err)
		}
		eventSwap, err := linodeClient.WaitForEventFinished(linode.ID, linodego.EntityLinode, linodego.ActionDiskCreate, diskSwap.Created, 240)
		// @TODO it is not sufficient that a disk was created. Which disk was it?
		// Sounds like we'll need a WaitForEntityStatus function.
		if err != nil {
			log.Fatalf("* Failed to wait for swap disk %d to finish creation: %s", diskSwap.ID, err)
		}
		if err := linodeClient.MarkEventRead(eventSwap); err != nil {
			log.Fatalln("* Failed to mark swap disk create event seen", err)
		}

		diskRaw, err := linodeClient.CreateInstanceDisk(linode.ID, linodego.InstanceDiskCreateOptions{Size: 50, Filesystem: "raw", Label: "linodego_raw"})
		if err != nil {
			log.Fatalln("* While creating raw disk:", err)
		}
		eventRaw, err := linodeClient.WaitForEventFinished(linode.ID, linodego.EntityLinode, linodego.ActionDiskCreate, diskRaw.Created, 240)
		// @TODO it is not sufficient that a disk was created. Which disk was it?
		// Sounds like we'll need a WaitForEntityStatus function.
		if err != nil {
			log.Fatalf("* Failed to wait for raw disk %d to finish creation: %s", diskRaw.ID, err)
		}
		if err := linodeClient.MarkEventRead(eventRaw); err != nil {
			log.Fatalln("* Failed to mark raw disk create event seen", err)
		}

		diskDebian, err := linodeClient.CreateInstanceDisk(
			linode.ID,
			linodego.InstanceDiskCreateOptions{
				Size:       1500,
				Filesystem: "ext4",
				Image:      "linode/debian9",
				Label:      "linodego_debian",
				RootPass:   randPassword(),
			},
		)
		if err != nil {
			log.Fatalln("* While creating Debian disk:", err)
		}
		eventDebian, err := linodeClient.WaitForEventFinished(linode.ID, linodego.EntityLinode, linodego.ActionDiskCreate, diskDebian.Created, 240)
		// @TODO it is not sufficient that a disk was created. Which disk was it?
		// Sounds like we'll need a WaitForEntityStatus function.
		if err != nil {
			log.Fatalf("* Failed to wait for Debian disk %d to finish creation: %s", diskDebian.ID, err)
		}
		if err := linodeClient.MarkEventRead(eventDebian); err != nil {
			log.Fatalln("* Failed to mark Debian disk create event seen", err)
		}
		fmt.Println("### Created Disks")

		createOpts := linodego.InstanceConfigCreateOptions{
			Devices: &linodego.InstanceConfigDeviceMap{
				SDA: &linodego.InstanceConfigDevice{DiskID: diskDebian.ID},
				SDB: &linodego.InstanceConfigDevice{DiskID: diskRaw.ID},
				SDC: &linodego.InstanceConfigDevice{DiskID: diskSwap.ID},
			},
			Kernel: "linode/direct-disk",
			Label:  "example config label",
			// RunLevel:   "default",
			// VirtMode:   "paravirt",
			Comments: "example config comment",
			// RootDevice: "/dev/sda",
			Helpers: &linodego.InstanceConfigHelpers{
				Network:    true,
				ModulesDep: false,
			},
		}
		config, err := linodeClient.CreateInstanceConfig(linode.ID, createOpts)
		if err != nil {
			log.Fatalln("* Failed to create Config", err)
		}
		fmt.Println("### Created Config:")
		updateOpts := linodego.InstanceConfigUpdateOptions{
			Comments: "updated example config comment",
		}
		config, err = linodeClient.UpdateInstanceConfig(linode.ID, config.ID, updateOpts)
		if err != nil {
			log.Fatalln("* Failed to update Config", err)
		}
		fmt.Println("### Updated Config:")

		booted, err := linodeClient.BootInstance(linode.ID, config.ID)
		if err != nil || !booted {
			log.Fatalln("* Failed to boot Instance", err)
		}
		fmt.Println("### Booted Instance")

		eventBooted, err := linodeClient.WaitForEventFinished(linode.ID, linodego.EntityLinode, linodego.ActionLinodeBoot, *config.Updated, 240)
		if err != nil {
			fmt.Println("### Boot Instance failed as expected:", err)
		} else {
			log.Fatalln("* Expected boot Instance to fail")
		}

		if err := linodeClient.MarkEventRead(eventBooted); err != nil {
			log.Fatalln("* Failed to mark boot event seen", err)
		}

		err = linodeClient.DeleteInstanceConfig(linode.ID, config.ID)
		if err != nil {
			log.Fatalln("* Failed to delete Config", err)
		}
		fmt.Println("### Deleted Config")

		err = linodeClient.DeleteInstanceDisk(linode.ID, diskSwap.ID)
		if err != nil {
			log.Fatalln("* Failed to delete Disk", err)
		}
		fmt.Println("### Deleted Disk")

		err = linodeClient.DeleteInstance(linode.ID)
		if err != nil {
			log.Fatalln("* Failed to delete Instance", err)
		}
		fmt.Println("### Deleted Instance")

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

		fmt.Println("## First Linode")

		configs, err := linodeClient.ListInstanceConfigs(linode.ID, nil)
		if err != nil {
			log.Fatal(err)
		} else if len(configs) > 0 {
			config, err := linodeClient.GetInstanceConfig(linode.ID, configs[0].ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("### First Config:", config.ID > 0)
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
			fmt.Println("### First Disk:", disk.ID > 0)
		} else {
			fmt.Println("### No Disks")
		}

		backups, err := linodeClient.GetInstanceBackups(linode.ID)
		if err != nil {
			log.Fatal(err)
		}
		if len(backups.Automatic) > 0 {
			fmt.Println("### First Auto Backup")
		} else {
			fmt.Println("### No Auto Backups")
		}
		fmt.Println("### Snapshots")
		if backups.Snapshot.Current != nil {
			// snapshot fetched will be exactly the same as backups.Snapshot.Current
			// just being redundant for illustrative purposes
			if snapshot, err := linodeClient.GetInstanceSnapshot(linode.ID, backups.Snapshot.Current.ID); err == nil {
				fmt.Println("#### Current:", snapshot.ID > 0)
			} else {
				fmt.Println("#### No Current Snapshot:", err)
			}
		} else {
			fmt.Println("### No Current Snapshot")
		}

		volumes, err := linodeClient.ListInstanceVolumes(linode.ID, nil)
		if err != nil {
			log.Fatal(err)
		} else if len(volumes) > 0 {
			volume, err := linodeClient.GetVolume(volumes[0].ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("### First Volume:", volume.ID > 0)
		} else {
			fmt.Println("### No Volumes")
		}

		stackscripts, err := linodeClient.ListStackscripts(&linodego.ListOptions{Filter: "{\"mine\":true}"})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("## Your Stackscripts:", len(stackscripts) > 0)
	}

	// Output:
	// ## Instance request with Invalid ID
	// ### Linode: <nil>
	// ### Error: [404] Not found
	// ## List Instances
	// ## First Linode
	// ### First Config: true
	// ### First Disk: true
	// ### No Auto Backups
	// ### Snapshots
	// #### Current: true
	// ### First Volume: true
	// ## Your Stackscripts: true
}

// randPassword generates a password sufficient to pass the Linode API standards,
// don't use it outside of this example script where the Linode is immediately destroyed.
func randPassword() string {
	const lowerBytes = "abcdefghijklmnopqrstuvwxyz"
	const upperBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const digits = "0123456789"
	const symbols = "/-=+@#$^&*()~!`|[]{}\\?,.<>;:'"

	length := 64 // must be divisible by character class count (4)
	b := make([]byte, length)

	for i := 0; i < length; i += 4 {
		b[i] = lowerBytes[rand.Intn(len(lowerBytes))]
		b[i+1] = upperBytes[rand.Intn(len(upperBytes))]
		b[i+2] = digits[rand.Intn(len(digits))]
		b[i+3] = symbols[rand.Intn(len(symbols))]
	}

	for i := range b {
		j := rand.Intn(i + 1)
		b[i], b[j] = b[j], b[i]
	}

	return string(b)
}
