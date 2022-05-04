package linodego

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"strings"
)

const DefaultConfigPath = "%s/.config/linode" // home dir
const DefaultConfigProfile = "default"

type ConfigProfile struct {
	APIToken   string `ini:"linode_api_token"`
	APIVersion string `ini:"linode_api_version"`
	APIURL     string `ini:"linode_api_url"`
}

type LoadConfigOptions struct {
	Path    string
	Profile string
}

// LoadConfig loads a Linode config according to the options argument.
// If no options are specified, the following defaults will be used:
// Path: ~/.config/linode
// Profile: default
func (c *Client) LoadConfig(options *LoadConfigOptions) error {
	path := options.Path
	if path == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		path = fmt.Sprintf(DefaultConfigPath, homeDir)
	}

	profileOption := options.Profile
	if profileOption == "" {
		profileOption = DefaultConfigProfile
	}

	cfg, err := ini.Load(path)
	if err != nil {
		return err
	}

	defaultConfig := ConfigProfile{
		APIToken:   "",
		APIURL:     APIHost,
		APIVersion: APIVersion,
	}

	if cfg.HasSection("default") {
		err := cfg.Section("default").MapTo(&defaultConfig)
		if err != nil {
			return fmt.Errorf("failed to map default profile: %s", err)
		}
	}

	result := make(map[string]ConfigProfile)

	for _, profile := range cfg.Sections() {
		name := strings.ToLower(profile.Name())

		f := defaultConfig
		if err := profile.MapTo(&f); err != nil {
			return fmt.Errorf("failed to map values: %s", err)
		}

		result[name] = f
	}

	c.configProfiles = result

	if err := c.UseProfile(profileOption); err != nil {
		return fmt.Errorf("unable to use profile %s: %s", profileOption, err)
	}

	return nil
}

// UseProfile switches client to use the specified profile.
// The specified profile must be already be loaded using client.LoadConfig(...)
func (c *Client) UseProfile(name string) error {
	name = strings.ToLower(name)

	profile, ok := c.configProfiles[name]
	if !ok {
		return fmt.Errorf("profile %s does not exist", name)
	}

	if profile.APIToken == "" {
		return fmt.Errorf("unable to resolve linode_api_token for profile %s", name)
	}

	if profile.APIURL == "" {
		return fmt.Errorf("unable to resolve linode_api_url for profile %s", name)
	}

	if profile.APIVersion == "" {
		return fmt.Errorf("unable to resolve linode_api_version for profile %s", name)
	}

	c.SetToken(profile.APIToken)
	c.SetBaseURL(profile.APIURL)
	c.SetAPIVersion(profile.APIVersion)

	return nil
}
