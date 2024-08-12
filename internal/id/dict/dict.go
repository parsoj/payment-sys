package dict

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"runtime"
	"time"
)

// StringArray holds the array of strings
type RandomWordsDictionary struct {
	array []string
}

// LoadFromFile loads the JSON array from the given file
func NewRandomWordsDictionary() (*RandomWordsDictionary, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("unable to get caller information")
	}
	dir := filepath.Dir(filename)
	filepath := filepath.Join(dir, "words_filtered.json")

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	d := RandomWordsDictionary{}

	err = json.Unmarshal(data, &d.array)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &d, nil
}

// GetRandomString returns a random string from the loaded array
func (sa *RandomWordsDictionary) GetRandomString() (string, error) {
	if len(sa.array) == 0 {
		return "", fmt.Errorf("array is empty")
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(sa.array))
	return sa.array[randomIndex], nil
}
