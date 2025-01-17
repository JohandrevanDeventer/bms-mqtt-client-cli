package persist

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type FilePersister struct {
	filePath string
	data     map[string]interface{}
}

func Test() {
	// fmt.Println("Test")
}

func NewFilePersister(filePath string) (*FilePersister, error) {
	// filePathDir := "./config/states"

	persister := &FilePersister{
		filePath: filePath,
		data:     make(map[string]interface{}),
	}

	filePathDir := filepath.Dir(filePath)

	_, err := os.Stat(filePathDir)
	if os.IsNotExist(err) {
		os.MkdirAll(filePathDir, os.ModePerm)
	}

	if err := persister.load(); err != nil {
		return nil, err
	}

	return persister, nil
}

func (p *FilePersister) load() error {
	file, err := os.Open(p.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist, start with empty data
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&p.data)
}

func (p *FilePersister) save() error {
	file, err := os.Create(p.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(p.data)
}

// Set allows setting a value using a nested key like "key1.key2.key3".
func (p *FilePersister) Set(key string, value interface{}) {
	keys := strings.Split(key, ".")
	lastKey := keys[len(keys)-1]
	current := p.data

	// Traverse or create intermediate maps
	for _, k := range keys[:len(keys)-1] {
		if _, ok := current[k]; !ok {
			current[k] = make(map[string]interface{})
		}
		if nextMap, ok := current[k].(map[string]interface{}); ok {
			current = nextMap
		} else {
			// Handle non-map conflicts
			current[k] = make(map[string]interface{})
			current = current[k].(map[string]interface{})
		}
	}

	// Set the final value (supports arrays, maps, etc.)
	current[lastKey] = value
	p.save()
}

// Get retrieves a value using a nested key like "key1.key2.key3".
func (p *FilePersister) Get(key string) interface{} {
	keys := strings.Split(key, ".")
	current := p.data

	// Traverse nested maps
	for _, k := range keys {
		if nextMap, ok := current[k].(map[string]interface{}); ok {
			current = nextMap
		} else {
			return current[k] // Return the value if it's not a map
		}
	}
	return nil // Key not found
}
