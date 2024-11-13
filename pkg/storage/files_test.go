package storage

import (
	"bufio"
	"bytes"
	"testing"
)

// TestFileW create a new file called .test writing the Text
// that said "filetest"
func TestFileW(test *testing.T) {
	file := NewFileStorage()
	file.WriteToFile("~/.test", "filetest", true)
	bytes, err := file.ReadFromFile("~/.test")
	if err != nil {
		test.Log(err)
	}
	test.Log(string(bytes))
}

// TestFileR read a file called .test
func TestFileR(test *testing.T) {
	file := NewFileStorage()
	data, err := file.ReadFromFile("~/.test")
	if err != nil {
		test.Log(err)
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	if scanner.Scan() {
		firstLine := scanner.Text()
		test.Log("Test line readed: ", firstLine)
	}
}
