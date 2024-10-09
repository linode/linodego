package unit

import (
	"embed"
	"encoding/json"
	"fmt"
	"path/filepath"
)

//go:embed fixtures/*.json
var fixtureFiles embed.FS

// TestFixtures manages loading and retrieving test fixtures
type TestFixtures struct {
	fixtures map[string]interface{}
}

// NewTestFixtures creates a new TestFixtures instance
func NewTestFixtures() *TestFixtures {
	tf := &TestFixtures{}
	tf.loadFixtures()
	return tf
}

// loadFixtures loads all embedded JSON files
func (tf *TestFixtures) loadFixtures() {
	tf.fixtures = make(map[string]interface{})

	entries, err := fixtureFiles.ReadDir("fixtures")
	if err != nil {
		panic(fmt.Sprintf("failed to read embedded fixtures: %v", err))
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			// Read the embedded JSON file
			data, err := fixtureFiles.ReadFile("fixtures/" + entry.Name())
			if err != nil {
				panic(fmt.Sprintf("failed to read fixture %s: %v", entry.Name(), err))
			}

			var jsonData interface{}
			if err := json.Unmarshal(data, &jsonData); err != nil {
				panic(fmt.Sprintf("failed to unmarshal fixture %s: %v", entry.Name(), err))
			}

			// Remove ".json" from the file name
			fixtureName := entry.Name()[:len(entry.Name())-5]
			tf.fixtures[fixtureName] = jsonData
		}
	}
}

// GetFixture retrieves the fixture data for the given name
func (tf *TestFixtures) GetFixture(name string) (interface{}, error) {
	data, ok := tf.fixtures[name]
	if !ok {
		return nil, fmt.Errorf("fixture not found: %s", name)
	}
	return data, nil
}
