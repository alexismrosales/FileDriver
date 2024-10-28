package internal

import (
	"cmd/client/internal/storage"
	"fmt"
	"os"
	"time"
)

const (
	INFO = iota
	WARN
	ERROR
	DEBUG
)

type Logger struct {
	logFilePath string
	file        *os.File
	fStorage    *storage.FileStorage
}

// New logger created, setting a file to write logs on a specific path
// setting the new file type and call the function NewWriter giving it the file
// as a parameter to be manipulated easier on the next functions
func NewLogger(logFilePath string) (*Logger, error) {
	// In case the file has ~/ symbol, the path is fixed
	logFilePath, err := storage.GetShortPath(logFilePath)
	if err != nil {
		return nil, err
	}
	// Creation or read of the file to print the log
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return &Logger{logFilePath: logFilePath, file: file, fStorage: storage.NewFileStorage()}, nil
}

// Print show the message formatted and then proceeds to write the log on the file
func (logger *Logger) Print(message string, level int) {
	var levelStr string
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelStr = levelToString(level)
	formatedMessage := fmt.Sprintf("[%s][%s]: %s", timestamp, levelStr, message)
	// DEBUG messages won´t be printed
	if level != DEBUG {
		fmt.Println(formatedMessage)
	}
	// Writing in file the formated message with level and time
	logger.fStorage.WriteToFile(logger.logFilePath, formatedMessage, true)
}

// levelToString converts the of the log into an string
func levelToString(level int) string {
	switch level {
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case DEBUG:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}
