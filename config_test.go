package linodego

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestConfig_LoadWithDefaults(t *testing.T) {
	client := NewClient(nil)

	file := createTestConfig(t, configLoadWithDefault)

	err := client.LoadConfig(&LoadConfigOptions{
		Path: file.Name(),
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(client.configProfiles) != 1 {
		fmt.Println(client.configProfiles)
		t.Fatalf("mismatched profile count: %d != %d", len(client.configProfiles), 1)
	}

	p, ok := client.configProfiles["default"]
	if !ok {
		t.Fatalf("default profile does not exist")
	}

	if p.APIToken != "blah" {
		t.Fatalf("mismatched api token: %s != %s", p.APIToken, "blah")
	}

	if p.APIURL != "api.cool.linode.com" {
		t.Fatalf("mismatched api url: %s != %s", p.APIURL, "api.cool.linode.com")
	}

	if p.APIVersion != "v4beta" {
		t.Fatalf("mismatched api version: %s != %s", p.APIVersion, "v4beta")
	}

	expectedURL := "https://api.cool.linode.com/v4beta"

	if client.resty.HostURL != expectedURL {
		t.Fatalf("mismatched host url: %s != %s", client.resty.HostURL, expectedURL)
	}

	if client.resty.Header.Get("Authorization") != "Bearer "+p.APIToken {
		t.Fatalf("token not found in auth header: %s", p.APIToken)
	}
}

func TestConfig_OverrideDefaults(t *testing.T) {
	client := NewClient(nil)

	file := createTestConfig(t, configOverrideDefaults)

	err := client.LoadConfig(&LoadConfigOptions{
		Path:    file.Name(),
		Profile: "cool",
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(client.configProfiles) != 2 {
		fmt.Println(client.configProfiles)
		t.Fatalf("mismatched profile count: %d != %d", len(client.configProfiles), 2)
	}

	p, ok := client.configProfiles["cool"]
	if !ok {
		t.Fatalf("cool profile does not exist")
	}

	if p.APIToken != "blah" {
		t.Fatalf("mismatched api token: %s != %s", p.APIToken, "blah")
	}

	if p.APIURL != "api.cool.linode.com" {
		t.Fatalf("mismatched api url: %s != %s", p.APIURL, "api.cool.linode.com")
	}

	if p.APIVersion != "v4" {
		t.Fatalf("mismatched api version: %s != %s", p.APIVersion, "v4")
	}

	expectedURL := "https://api.cool.linode.com/v4"

	if client.resty.HostURL != expectedURL {
		t.Fatalf("mismatched host url: %s != %s", client.resty.HostURL, expectedURL)
	}

	if client.resty.Header.Get("Authorization") != "Bearer "+p.APIToken {
		t.Fatalf("token not found in auth header: %s", p.APIToken)
	}
}

func TestConfig_NoDefaults(t *testing.T) {
	client := NewClient(nil)

	file := createTestConfig(t, configNoDefaults)

	err := client.LoadConfig(&LoadConfigOptions{
		Path:    file.Name(),
		Profile: "cool",
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(client.configProfiles) != 2 {
		fmt.Println(client.configProfiles)
		t.Fatalf("mismatched profile count: %d != %d", len(client.configProfiles), 2)
	}

	p, ok := client.configProfiles["cool"]
	if !ok {
		t.Fatalf("cool profile does not exist")
	}

	if p.APIToken != "mytoken" {
		t.Fatalf("mismatched api token: %s != %s", p.APIToken, "mytoken")
	}

	if client.resty.Header.Get("Authorization") != "Bearer "+p.APIToken {
		t.Fatalf("token not found in auth header: %s", p.APIToken)
	}
}

func createTestConfig(t *testing.T, conf string) *os.File {
	file, err := ioutil.TempFile("", "linode")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Fprint(file, conf)

	t.Cleanup(func() {
		file.Close()
		os.Remove(file.Name())
	})

	return file
}

const configLoadWithDefault = `
[default]
token = blah
api_url = api.cool.linode.com
api_version = v4beta
`

const configOverrideDefaults = `
[default]
token = blah
api_url = api.cool.linode.com
api_version = v4beta

[cool]
api_version = v4
# Values from default are inherited here
`

const configNoDefaults = `
[cool]
token = mytoken
# Linodego default values are inherited here
`
