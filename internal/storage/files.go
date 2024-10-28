package storage

import (
	"os"
	"path/filepath"
)

type FileStorage struct{}

func NewFileStorage() *FileStorage {
	return &FileStorage{}
}

// WriteToFile the data on the path given, if the flag "appendData" is up, the const that
// append the data in the file will be set. It is recommended to activate on text files.
func (fs *FileStorage) WriteToFile(filePath string, data any, appendData bool) error {
	filePath, err := GetShortPath(filePath)
	if err != nil {
		return err
	}
	var file *os.File
	// Appending text in case is a text file
	var constFile int
	if appendData {
		constFile = os.O_APPEND
	} else {
		constFile = os.O_TRUNC
	}
	file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|constFile, 0644)

	// Creating a new generic writer
	writer := NewGWriter(file)
	if err != nil {
		return err
	}
	defer file.Close()
	// Writing data
	err = writer.write(data)
	return err
}

// ReadFromFile converts all the file into bytes
func (fs *FileStorage) ReadFromFile(filePath string) ([]byte, error) {
	filePath, err := GetShortPath(filePath)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(filePath)
}

// getShortPath in case the file has ~/ symbol the path will be fixed
func GetShortPath(filePath string) (string, error) {
	if len(filePath) >= 2 && filePath[:2] == "~/" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		filePath = filepath.Join(homeDir, filePath[2:])
	}
	return filePath, nil
}
