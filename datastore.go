package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

var once sync.Once
var instance *DataStore

// Max cache size
const Max = 1000000

type DataStore struct {
	store        map[string]string
	filename     string
	numEntries   int64
	cacheEnabled bool
}

func (ds *DataStore) load() error {
	file, err := os.Open(ds.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid line: %s", line)
		}
		ds.store[parts[0]] = parts[1]
		ds.numEntries++

		if ds.numEntries > Max {
			fmt.Println("Max capacity exceeded, turning cache Mode on")
			ds.cacheEnabled = true
			return nil
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (ds *DataStore) Get(key string) (string, bool) {
	value, ok := ds.store[key]

	if ds.cacheEnabled && !ok {
		v, err := ds.fetchValue(key)
		if err != nil {
			fmt.Printf("Encountered cache error: %s\n", err.Error())
			return "", false
		}
		return v, true
	}
	return value, ok
}

// only used when caching is enabled
func (ds *DataStore) fetchValue(key string) (string, error) {
	file, err := os.Open(ds.filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid line: %s", line)
		}

		if parts[0] == key {
			return parts[1], nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("key not found: %s", key)
}

func initDataStore(dsfile string) *DataStore {
	return &DataStore{
		filename:     dsfile,
		store:        make(map[string]string),
		numEntries:   0,
		cacheEnabled: false, //disabled by default
	}
}

func getInstance(dsfile string) *DataStore {
	once.Do(func() {
		instance = initDataStore(dsfile)
		if err := instance.load(); err != nil {
			fmt.Printf("Error in reading datastore file %s, %s\n", instance.filename, err.Error())
		}
	})
	return instance
}
