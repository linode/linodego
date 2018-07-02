package linodego_test

import (
	"fmt"
	"log"

	"github.com/chiefy/linodego"
)

func ExampleCreateNodeBalancer() {
	// Example readers, Ignore this bit of setup code needed to record test fixtures
	linodeClient, teardown := createTestClient(nil, "fixtures/ExampleCreateNodeBalancer")
	defer teardown()

	fmt.Println("## NodeBalancer create")
	var nbID int
	var nb = &linodego.NodeBalancer{
		ClientConnThrottle: 20,
		Region:             "us-east",
	}

	createOpts := nb.GetCreateOptions()
	nb, err := linodeClient.CreateNodeBalancer(&createOpts)
	if err != nil {
		log.Fatal(err)
	}
	nbID = nb.ID

	fmt.Println("### Get")
	nb, err = linodeClient.GetNodeBalancer(nbID)
	if err != nil {
		log.Fatal(err)
	}

	updateOpts := nb.GetUpdateOptions()
	*updateOpts.Label += "_renamed"
	nb, err = linodeClient.UpdateNodeBalancer(nbID, updateOpts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("### Delete")
	if err := linodeClient.DeleteNodeBalancer(nbID); err != nil {
		log.Fatal(err)
	}

	// Output:
	// ## NodeBalancer create
	// ### Get
	// ### Delete
}
