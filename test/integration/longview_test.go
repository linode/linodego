package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
)

func TestLongviewClient_List(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestLongviewClient_List")
	defer teardown()

	longviewClients, err := client.ListLongviewClients(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing longview clients, expected struct - error %v", err)
	}
	if len(longviewClients) == 0 {
		t.Errorf("Expected a list longview clients - %v", longviewClients)
	}
}

func TestLongviewClient_Get(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestLongviewClient_Get")
	defer teardown()

	createOptions := linodego.LongviewClientCreateOptions{Label: "testing"}

	// Create a longview client to later fetch using get
	testingLongviewClient, createErr := client.CreateLongviewClient(context.TODO(), createOptions)
	if createErr != nil {
		t.Errorf("Error creating longview client:%s", createErr)
	}

	t.Cleanup(func() {
		if err := client.DeleteLongviewClient(context.Background(), testingLongviewClient.ID); err != nil {
			t.Fatal(err)
		}
	})

	// Fetch the ID of the newly created longview client
	testingID := testingLongviewClient.ID

	// If there is no error, then GetLongviewClient works properly
	longviewClient, getErr := client.GetLongviewClient(context.Background(), testingID)
	if getErr != nil {
		t.Errorf("Error getting longview client:%s", getErr)
	}

	longviewClientList, listErr := client.ListLongviewClients(context.Background(), nil)
	if listErr != nil {
		t.Errorf("Error listing longview clients:%s", listErr)
	}

	found := false
	for _, element := range longviewClientList {
		if element.ID == longviewClient.ID {
			found = true
		}
	}

	if !found {
		t.Errorf("Longview client not found in list.")
	}
}

func TestLongviewClient_Create(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestLongviewClient_Create")
	defer teardown()

	createOptions := linodego.LongviewClientCreateOptions{Label: "testing"}

	// Create a longview client to later fetch using get
	testingLongviewClient, createErr := client.CreateLongviewClient(context.TODO(), createOptions)
	if createErr != nil {
		t.Errorf("Error creating longview client:%s", createErr)
	}

	t.Cleanup(func() {
		if err := client.DeleteLongviewClient(context.Background(), testingLongviewClient.ID); err != nil {
			t.Fatal(err)
		}
	})

	testingID := testingLongviewClient.ID

	// If there is no error, then TestLongviewClient_Create works properly
	_, getErr := client.GetLongviewClient(context.Background(), testingID)
	if getErr != nil {
		t.Errorf("Error getting longview client:%s", getErr)
	}
}

func TestLongviewClient_Delete(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestLongviewClient_Delete")
	defer teardown()

	createOptions := linodego.LongviewClientCreateOptions{Label: "testing"}

	// Create a longview client to later delete
	testingLongviewClient, createErr := client.CreateLongviewClient(context.TODO(), createOptions)
	if createErr != nil {
		t.Errorf("Error creating longview client:%s", createErr)
	}

	// Fetch the ID of the newly created longview client
	testingID := testingLongviewClient.ID

	// Make sure the longview client exists
	longviewClient, getErr := client.GetLongviewClient(context.Background(), testingID)
	if getErr != nil {
		t.Errorf("Error getting longview client:%s", getErr)
	}

	// If there is no error, the longview client was delete properly
	if err := client.DeleteLongviewClient(context.Background(), testingLongviewClient.ID); err != nil {
		t.Fatal(err)
	}

	longviewClientList, listErr := client.ListLongviewClients(context.Background(), nil)
	if listErr != nil {
		t.Errorf("Error listing longview clients:%s", listErr)
	}

	found := false
	for _, element := range longviewClientList {
		if element.ID == longviewClient.ID {
			found = true
		}
	}

	// If the longview client still appears in the list, then it was not deleted properly
	if found {
		t.Errorf("Longview client was found in list.")
	}
}

func TestLongviewClient_Update(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestLongviewClient_Update")
	defer teardown()

	createOptions := linodego.LongviewClientCreateOptions{Label: "testing"}

	// Create a longview client to later update
	testingLongviewClient, createErr := client.CreateLongviewClient(context.TODO(), createOptions)
	if createErr != nil {
		t.Errorf("Error creating longview client:%s", createErr)
	}

	// Fetch the ID of the newly created longview client
	testingID := testingLongviewClient.ID

	updateOptions := linodego.LongviewClientUpdateOptions{Label: "testing_updated"}

	// Update the longview client
	updatedTestingLongviewClient, updateErr := client.UpdateLongviewClient(context.TODO(), testingID, updateOptions)
	if createErr != nil {
		t.Errorf("Error updating longview client:%s", updateErr)
	}

	t.Cleanup(func() {
		if err := client.DeleteLongviewClient(context.Background(), updatedTestingLongviewClient.ID); err != nil {
			t.Fatal(err)
		}
	})

	// If the label does not match what it was updated to, the update was not successful
	if updatedTestingLongviewClient.Label != "testing_updated" {
		t.Errorf("Longview client was not updated.")
	}
}

func TestLongviewPlan_Get(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestLongviewPlan_Get")
	defer teardown()

	// Get the current longview plan
	testingLongviewPlan, err := client.GetLongviewPlan(context.Background())
	if err != nil {
		t.Errorf("Error getting longview plan:%s", err)
	}

	// Ensure that a plan was returned
	if testingLongviewPlan == nil {
		t.Errorf("Expected a longview plan, recieved nil")
	}

	validIDs := []string{"longview-3", "longview-10", "longview-40", "longview-100", ""}

	// Ensure that the id of the longview plan is a valid one
	if !contains(validIDs, testingLongviewPlan.ID) {
		t.Errorf("Invalid longview plan ID:%s", testingLongviewPlan.ID)
	}
}

func TestLongviewPlan_Update(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestLongviewPlan_Update")
	defer teardown()

	// Get the current longview plan
	testingLongviewPlan, getErr := client.GetLongviewPlan(context.Background())
	if getErr != nil {
		t.Errorf("Error getting longview plan:%s", getErr)
	}

	testingID := testingLongviewPlan.ID
	updateOptions := linodego.LongviewPlanUpdateOptions{LongviewSubscription: "longview-10"}

	// Update the longview plan
	updatedLongviewPlan, updateErr := client.UpdateLongviewPlan(context.Background(), updateOptions)
	if updateErr != nil {
		t.Errorf("Error updating longview plan:%s", updateErr)
	}

	// Set the longview plan to what it was before (after the test)
	t.Cleanup(func() {
		resetOptions := linodego.LongviewPlanUpdateOptions{LongviewSubscription: testingID}

		if _, err := client.UpdateLongviewPlan(context.Background(), resetOptions); err != nil {
			t.Fatal(err)
		}
	})

	// Ensure the longview plan was updated correctly
	if updatedLongviewPlan.ID != "longview-10" {
		t.Errorf("Longview plan not updated")
	}

}

func contains(arr []string, elem string) bool {
	for _, curr := range arr {
		if curr == elem {
			return true
		}
	}
	return false
}
