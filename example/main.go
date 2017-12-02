package main

import (
	"fmt"
	"log"
	"os"

	golinode "github.com/chiefy/go-linode"
)

func main() {
	apiKey, ok := os.LookupEnv("LINODE_API_KEY")
	if !ok {
		log.Fatal("Could not find LINODE_API_KEY, please assert it is set.")
	}
	linodeClient, err := golinode.NewClient(apiKey)
	if err != nil {
		log.Fatal(err)
	}
	linodeClient.SetDebug(true)
	res, err := linodeClient.GetInstance(4090913)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", res)

}
