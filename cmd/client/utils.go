package main

import (
	"encoding/json"
	"github.com/alexismrosales/FileDriver/pkg/storage"
	"os"
	"path/filepath"
	"strings"
)

type info struct {
	Ip   string
	Port int
}

// Serialize the info into a json format
func saveInfo(data info) error {
	bData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	strg := storage.NewFileStorage()
	// Saving server info in infoRoute
	err = strg.WriteToFile(infoRoute, bData, false)
	if err != nil {
		return err
	}
	return nil
}

func getInfo() (*info, error) {
	strg := storage.NewFileStorage()
	data, err := strg.ReadFromFile(infoRoute)
	if err != nil {
		return nil, err
	}
	var getInfo info
	err = json.Unmarshal(data, &getInfo)
	if err != nil {
		return nil, err
	}
	return &getInfo, nil
}

func getPathsFull(paths ...string) error {
	execPath, err := os.Executable()
	homePath, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	fullPath := filepath.Join(execPath, "..")
	fullPath = filepath.Clean(fullPath)

	for i, path := range paths {
		if !strings.HasPrefix(path, homePath) {
			tmpPath := filepath.Clean(path)
			paths[i] = filepath.Join(fullPath, tmpPath)
		}
	}
	return nil
}
