package linodego_test

import (
	"context"
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
	nb, err := linodeClient.CreateNodeBalancer(context.TODO(), &createOpts)
	if err != nil {
		log.Fatal(err)
	}
	nbID = nb.ID

	fmt.Println("### Get")
	nb, err = linodeClient.GetNodeBalancer(context.TODO(), nbID)
	if err != nil {
		log.Fatal(err)
	}

	updateOpts := nb.GetUpdateOptions()
	*updateOpts.Label += "_renamed"
	nb, err = linodeClient.UpdateNodeBalancer(context.TODO(), nbID, updateOpts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("### Delete")
	if err := linodeClient.DeleteNodeBalancer(context.TODO(), nbID); err != nil {
		log.Fatal(err)
	}

	// Output:
	// ## NodeBalancer create
	// ### Get
	// ### Delete
}

func ExampleCreateNodeBalancerConfig() {
	// Example readers, Ignore this bit of setup code needed to record test fixtures
	linodeClient, teardown := createTestClient(nil, "fixtures/ExampleCreateNodeBalancerConfig")
	defer teardown()

	fmt.Println("## NodeBalancer create")
	clientConnThrottle := 20
	nb, err := linodeClient.CreateNodeBalancer(context.TODO(), &linodego.NodeBalancerCreateOptions{
		ClientConnThrottle: &clientConnThrottle,
		Region:             "us-east",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("## NodeBalancer Config create")
	createOpts := linodego.NodeBalancerConfigCreateOptions{
		Port: 80,
		/*
			Protocol:      linodego.ProtocolHTTP,
			Algorithm:     linodego.AlgorithmLeastConn,
			Stickiness:    linodego.StickinessHTTPCookie,
			Check:         linodego.CheckHTTP,
			CheckInterval: 30,
			CheckAttempts: 5,
			CipherSuite:   linodego.CipherRecommended,
		*/
	}
	nbc, err := linodeClient.CreateNodeBalancerConfig(context.TODO(), nb.ID, &createOpts)
	if err != nil {
		log.Fatal(err)
	}
	nbcID := nbc.ID

	fmt.Println("## NodeBalancer Config update")
	updateOpts := nbc.GetUpdateOptions()
	updateOpts.Port += 8000
	nbc, err = linodeClient.UpdateNodeBalancerConfig(context.TODO(), nb.ID, nbc.ID, updateOpts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("### List")
	configs, err := linodeClient.ListNodeBalancerConfigs(context.TODO(), nb.ID, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("### Get")
	nbc, err = linodeClient.GetNodeBalancerConfig(context.TODO(), nb.ID, configs[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("### Delete")
	if nbc.ID != nbcID {
		log.Fatalf("Unexpected Nodebalancer Config ID %d != %d", nbc.ID, nbcID)
	}
	if err := linodeClient.DeleteNodeBalancerConfig(context.TODO(), nb.ID, nbc.ID); err != nil {
		log.Fatal(err)
	}

	if err := linodeClient.DeleteNodeBalancer(context.TODO(), nb.ID); err != nil {
		log.Fatal(err)
	}

	// Output:
	// ## NodeBalancer create
	// ## NodeBalancer Config create
	// ## NodeBalancer Config update
	// ### List
	// ### Get
	// ### Delete
}

func ExampleCreateNodeBalancerNode() {
	// Example readers, Ignore this bit of setup code needed to record test fixtures
	linodeClient, teardown := createTestClient(nil, "fixtures/ExampleCreateNodeBalancerNode")
	defer teardown()

	fmt.Println("## NodeBalancer create")
	clientConnThrottle := 20
	nb, err := linodeClient.CreateNodeBalancer(context.TODO(), &linodego.NodeBalancerCreateOptions{
		ClientConnThrottle: &clientConnThrottle,
		Region:             "us-east",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("## NodeBalancer Config create")

	nbc, err := linodeClient.CreateNodeBalancerConfig(context.TODO(), nb.ID, &linodego.NodeBalancerConfigCreateOptions{
		Port: 80,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("## NodeBalancer Node create")
	createOpts := linodego.NodeBalancerNodeCreateOptions{
		Address: "192.168.129.255:80",
		Label:   "192.168.129.255-80",
	}
	nbn, err := linodeClient.CreateNodeBalancerNode(context.TODO(), nb.ID, nbc.ID, &createOpts)
	if err != nil {
		log.Fatal(err)
	}
	nbnID := nbn.ID

	fmt.Println("## NodeBalancer Node update")
	updateOpts := nbn.GetUpdateOptions()
	updateOpts.Address = "192.168.129.0:8080"
	nbn, err = linodeClient.UpdateNodeBalancerNode(context.TODO(), nb.ID, nbc.ID, nbn.ID, updateOpts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("### List")
	nodes, err := linodeClient.ListNodeBalancerNodes(context.TODO(), nb.ID, nbc.ID, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("### Get")
	nbn, err = linodeClient.GetNodeBalancerNode(context.TODO(), nb.ID, nbc.ID, nodes[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("### Delete")
	if nbn.ID != nbnID {
		log.Fatalf("Unexpected Nodebalancer Node ID %d != %d", nbn.ID, nbnID)
	}
	if err := linodeClient.DeleteNodeBalancerNode(context.TODO(), nb.ID, nbc.ID, nbn.ID); err != nil {
		log.Fatal(err)
	}

	if err := linodeClient.DeleteNodeBalancerConfig(context.TODO(), nb.ID, nbc.ID); err != nil {
		log.Fatal(err)
	}

	if err := linodeClient.DeleteNodeBalancer(context.TODO(), nb.ID); err != nil {
		log.Fatal(err)
	}

	// Output:
	// ## NodeBalancer create
	// ## NodeBalancer Config create
	// ## NodeBalancer Node create
	// ## NodeBalancer Node update
	// ### List
	// ### Get
	// ### Delete
}
