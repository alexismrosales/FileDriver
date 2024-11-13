package main

import (
	"github.com/alexismrosales/FileDriver/pkg/storage"
	"os"
	"path/filepath"
)

// pathExists verify if the path is valid
func pathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// createDir Create a new directory allowing short paths
func createDir(path string) error {
	sPath, err := storage.GetShortPath(path)
	if err != nil {
		return err
	}
	err = os.MkdirAll(sPath, 0755)
	if err != nil {
		return err
	}
	return nil
}

// safePath Using the base directory path, joins a new path not allowing the user navigate out the dir
func safePath(path string) (string, error) {
	serverPath, err := storage.GetShortPath(baseDir)
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(serverPath, path)
	return fullPath, nil
}

// getBaseShortPath Get a short path with the server direction
func getBaseShortPath(path string) string {
	if len(path) >= 2 && path[:2] == "~/" {
		path = filepath.Join(baseDir, path[2:])
		return path[len(baseDir)+1:]
	}
	return filepath.Join(baseDir, path)
}
