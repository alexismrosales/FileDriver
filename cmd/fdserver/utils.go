package main

import (
	"encoding/json"
	"fmt"
	"github.com/alexismrosales/FileDriver/pkg/storage"
	"os"
	"time"
)

type status struct {
	Path string
}

// Serialize the info into a json format
func saveStatus(path string) error {
	status := &status{Path: path}
	bData, err := json.Marshal(status)
	if err != nil {
		return err
	}
	strg := storage.NewFileStorage()
	// Saving server info in infoRoute
	err = strg.WriteToFile(statusPath, bData, false)
	if err != nil {
		return err
	}
	return nil
}

func getStatus() (*status, error) {
	strg := storage.NewFileStorage()
	data, err := strg.ReadFromFile(statusPath)
	if err != nil {
		return nil, err
	}
	var getStatus status
	err = json.Unmarshal(data, &getStatus)
	if err != nil {
		return nil, err
	}
	return &getStatus, nil
}

func getStatusPath() (string, error) {
	sts, err := getStatus()
	if err != nil {
		return "", err
	}
	return sts.Path, nil
}

// createBaseDir Create a directory using the base directory
func createBaseDir() error {
	err := createDir(baseDir)
	if err != nil {
		// In case the directory is already created
		if os.IsNotExist(err) {
			return nil
		}
	}
	return err
}

func createZipFileName() string {
	timestamp := time.Now().Format("2006-01-02-15:04")
	return fmt.Sprintf("filedrive-download-%s.zip", timestamp)
}
