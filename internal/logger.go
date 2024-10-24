package internal

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger struct {
	file    *os.File
	gWriter *GenericWriter
}

// New logger created, setting a file to write logs on a specific path
// setting the new file type and call the function NewWriter giving it the file
// as a parameter to be manipulated easier on the next functions
func NewLogger(logFilePath string) (*Logger, error) {
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	// Writing log on file
	log.SetOutput(file)
	return &Logger{file: file, gWriter: NewGWriter(file)}, nil
}

// Print show the message formatted and then proceeds to write the log on the file
func (logger *Logger) Print(level, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formatedMessage := fmt.Sprintf("[%s][%s]: %s", timestamp, level, message)
	log.Println(formatedMessage)
	logger.gWriter.Writer(formatedMessage)
}
