package integration

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/linode/linodego"
)

var (
	client       *linodego.Client
	cleanup      func()
	firewallID   int
	fixturesYaml = "fixtures/TestLinodeCloudFirewall"
)

func TestMain(m *testing.M) {
	if envFixtureMode, ok := os.LookupEnv("LINODE_FIXTURE_MODE"); ok {
		switch envFixtureMode {
		case "record":
			setupCloudFirewall(nil)
		case "play":
			log.Printf("[INFO] Fixture mode play - Test Linode Cloud Firewall not created")
		}
	}

	code := m.Run()

	if envFixtureMode, ok := os.LookupEnv("LINODE_FIXTURE_MODE"); ok && envFixtureMode == "record" {
		deleteCloudFirewall()
	}

	os.Exit(code)
}

func setupCloudFirewall(t *testing.T) {
	client, cleanup = createTestClient(t, fixturesYaml)

	publicIPv4, err := getPublicIPv4()
	if err != nil {
		t.Fatalf("[ERROR] Failed to retrieve public IPv4: %v", err)
	}

	firewallRuleSet := getDefaultFirewallRuleSet(publicIPv4)
	firewallLabel := fmt.Sprintf("cloudfw-%d", time.Now().UnixNano())

	firewall, err := client.CreateFirewall(context.Background(), linodego.FirewallCreateOptions{
		Label: &firewallLabel,
		Rules: firewallRuleSet,
	})
	if err != nil {
		log.Printf("[ERROR] Error creating firewall: %v\n", err)
		os.Exit(1)
	}

	firewallID = firewall.ID
	log.Printf("[INFO] Created Test Linode Cloud Firewall with ID: %d\n", firewallID)
}

func deleteCloudFirewall() {
	if firewallID != 0 {
		err := client.DeleteFirewall(context.Background(), firewallID)
		if err != nil {
			log.Printf("[ERROR] Error deleting Cloud Firewall: %v\n", err)
			os.Exit(1)
		}
		log.Printf("[INFO] Deleted Test Linode Cloud Firewall with ID: %d\n", firewallID)
	}
}

func getDefaultFirewallRuleSet(publicIPv4 string) linodego.FirewallRuleSet {
	cloudFirewallRule := linodego.FirewallRule{
		Label:     "ssh-inbound-accept-local",
		Action:    "ACCEPT",
		Ports:     linodego.Pointer("22"),
		Protocol:  "TCP",
		Addresses: linodego.NetworkAddresses{IPv4: &[]string{publicIPv4}},
	}

	return linodego.FirewallRuleSet{
		Inbound:        []linodego.FirewallRule{cloudFirewallRule},
		InboundPolicy:  "DROP",
		Outbound:       []linodego.FirewallRule{},
		OutboundPolicy: "ACCEPT",
	}
}

func getPublicIPv4() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body) + "/32", nil
}

func GetFirewallID() int {
	return firewallID
}
