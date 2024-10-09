package unit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

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

// loadFixtures loads all JSON files in fixtures directory
func (tf *TestFixtures) loadFixtures() {
	tf.fixtures = make(map[string]interface{})
	fixturesDir := "fixtures"

	err := filepath.Walk(fixturesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			var jsonData interface{}
			if err := json.Unmarshal(data, &jsonData); err != nil {
				return err
			}
			fixtureName := filepath.Base(path)
			tf.fixtures[fixtureName[:len(fixtureName)-5]] = jsonData
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

// GetFixture retrieves the fixture data for the given URL
func (tf *TestFixtures) GetFixture(jsonFile string) (interface{}, error) {
	data, ok := tf.fixtures[jsonFile]
	if !ok {
		return nil, os.ErrNotExist
	}
	fmt.Printf("Loading fixture: %s\n", jsonFile) // TODO:: debugging remove later

	return data, nil
}
