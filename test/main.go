package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	golinode "github.com/chiefy/go-linode"
	"github.com/dnaeon/go-vcr/recorder"
)

func main() {

	envErr := fmt.Errorf("env vars LINODE_INSTANCE_ID and LINODE_VOLUME_ID must be set")

	var linodeInstanceID int
	var linodeVolumeID int
	var err error

	if linodeInstanceID, err = strconv.Atoi(os.Getenv("LINODE_INSTANCE_ID")); err != nil {
		log.Fatal(envErr)
	}
	if linodeVolumeID, err = strconv.Atoi(os.Getenv("LINODE_VOLUME_ID")); err != nil {
		log.Fatal(envErr)
	}

	// Start our recorder
	r, err := recorder.New("test/fixtures")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Stop() // Make sure recorder is stopped once done with it

	c := golinode.NewClient(nil, r)
	if err != nil {
		log.Fatal(err)
	}

	_, err = c.GetType("g6-standard-1")
	if err != nil {
		log.Fatal(err)
	}

	_, _ = c.GetType("missing-type")

	_, err = c.ListTypes(nil)
	if err != nil {
		log.Fatal(err)
	}

	_, err = c.ListKernels(nil)
	if err != nil {
		log.Fatal(err)
	}

	filterOpt := golinode.ListOptions{Filter: "{\"label\":\"Recovery - Finnix (kernel)\"}"}
	_, err = c.ListKernels(&filterOpt)
	if err != nil {
		log.Fatal(err)
	}

	_, err = c.ListKernels(golinode.NewListOptions(1, ""))
	if err != nil {
		log.Fatal(err)
	}

	_, err = c.ListImages(nil)
	if err != nil {
		log.Fatal(err)
	}

	_, _ = c.GetImage("does-not-exist")

	_, err = c.GetImage("linode/ubuntu16.04lts")
	if err != nil {
		log.Fatal(err)
	}

	pageOpt := golinode.ListOptions{PageOptions: &golinode.PageOptions{Page: 1}}
	_, err = c.ListLongviewSubscriptions(&pageOpt)
	if err != nil {
		log.Fatal(err)
	}

	_, err = c.ListRegions(nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Succesfully got linode regions")

	_, err = c.ListInstances(nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Succesfully got linode instances")

	_, err = c.GetInstance(linodeInstanceID)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("Succesfully got linode instance ID %d", linodeInstanceID))

	_, err = c.GetInstanceBackups(linodeInstanceID)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("Succesfully got linode backups for instance ID %d", linodeInstanceID))

	_, err = c.ListInstanceDisks(linodeInstanceID, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Succesfully got linode instance disks")

	_, err = c.ListInstanceConfigs(linodeInstanceID, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Succesfully got linode instance configs")

	_, err = c.ListInstanceVolumes(linodeInstanceID, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Succesfully got linode instance volumes")

	_, err = c.ListStackscripts(golinode.NewListOptions(1, ""))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Succesfully got linode public stackscripts (1 page)")

	_, err = c.GetStackscript("7")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Succesfully got linode stackscript ID 7")

	_, err = c.ListVolumes(nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Succesfully got linode volumes (1 page)")

	_, err = c.GetVolume(linodeVolumeID)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("Succesfully got linode volume ID %d", linodeVolumeID))

	log.Printf("Successfully retrieved linode requests!")
}
