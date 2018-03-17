package main

import (
	"log"

	golinode "github.com/chiefy/go-linode"
	"github.com/dnaeon/go-vcr/recorder"
)

func main() {
	// Start our recorder
	r, err := recorder.New("test/fixtures")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Stop() // Make sure recorder is stopped once done with it

	c, err := golinode.NewClient(nil, r)
	if err != nil {
		log.Fatalf("Failed to create linode client: %s", err)
	}
	c.SetDebug(false)

	_, err = c.ListRegions()
	if err != nil {
		log.Fatalf("Failed to get linode regions: %s", err)
	}
	log.Println("Succesfully got linode regions")

	_, err = c.ListInstances()
	if err != nil {
		log.Fatalf("Failed to get linode instances: %s", err)
	}
	log.Println("Succesfully got linode instances")

	_, err = c.ListImages()
	if err != nil {
		log.Fatalf("Failed to get linode images: %s", err)
	}
	log.Println("Succesfully got linode images")

	_, err = c.GetInstance(6809519)
	if err != nil {
		log.Fatalf("Failed to get linode instance ID 6809519: %s", err)
	}
	log.Println("Succesfully got linode instance ID 6809519")

	_, err = c.ListStackscripts()
	if err != nil {
		log.Fatalf("Failed to get linode stackscripts: %s", err)
	}
	log.Println("Succesfully got linode public stackscripts (1 page)")

	_, err = c.GetStackscript(7)
	if err != nil {
		log.Fatalf("Failed to get linode stackscript ID 7: %s", err)
	}
	log.Println("Succesfully got linode stackscript ID 7")

	_, err = c.ListVolumes()
	if err != nil {
		log.Fatalf("Failed to get linode volumes: %s", err)
	}
	log.Println("Succesfully got linode volumes (1 page)")

	_, err = c.GetVolume(4880)
	if err != nil {
		log.Fatalf("Failed to get linode volume ID 4880: %s", err)
	}
	log.Println("Succesfully got linode volume ID 4880")

	log.Printf("Successfully retrieved linode requests!")
}
