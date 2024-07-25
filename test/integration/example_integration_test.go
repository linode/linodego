package integration

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
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

func ExampleClient_GetAccount() {
	// Example readers, Ignore this bit of setup code needed to record test fixtures
	linodeClient, teardown := createTestClient(nil, "fixtures/ExampleGetAccount")
	defer teardown()

	account, err := linodeClient.GetAccount(context.Background())
	if err != nil {
		log.Fatalln("* While getting account: ", err)
	}
	fmt.Println("Account has email:", len(account.Email) > 0)

	// Output:
	// Account has email: true
}

func ExampleClient_ListUsers() {
	// Example readers, Ignore this bit of setup code needed to record test fixtures
	linodeClient, teardown := createTestClient(nil, "fixtures/ExampleListUsers")
	defer teardown()

	users, err := linodeClient.ListUsers(context.Background(), nil)
	if err != nil {
		log.Fatalln("* While getting users: ", err)
	}
	fmt.Println("User exists:", len(users) > 0)

	// Output:
	// User exists: true
}

const (
	lowerBytes = "abcdefghijklmnopqrstuvwxyz"
	upperBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits     = "0123456789"
	symbols    = "/-=+@#$^&*()~!`|[]{}\\?,.<>;:'"
)

func randString(length int, characterClasses ...string) string {
	quotient := (0.0 + float64(length)) / float64(len(characterClasses))
	if quotient != math.Trunc(quotient) {
		panic("length must be divisible by characterClasses count")
	}

	b := make([]byte, length)

	for i := 0; i < length; i += len(characterClasses) {
		for j, characterClass := range characterClasses {
			randPos := rand.Intn(len(characterClass))
			b[i+j] = characterClass[randPos]
		}
	}

	for i := range b {
		j := rand.Intn(i + 1)
		b[i], b[j] = b[j], b[i]
	}

	return string(b)
}

func randPassword() string {
	return randString(64, lowerBytes, upperBytes, digits, symbols)
}

func randLabel() string {
	return randString(12, lowerBytes, digits)
}
