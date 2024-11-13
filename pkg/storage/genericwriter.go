package storage

import (
	"bufio"
	"os"
)

type genericWriter struct {
	writer *bufio.Writer
}

// NewWriter creates a new write buffer
func NewGWriter(file *os.File) *genericWriter {
	writer := bufio.NewWriter(file)
	return &genericWriter{writer: writer}
}

// Writer selects the type of data to writes it on the file
func (gWriter *genericWriter) write(data any) error {
	var err error
	switch d := data.(type) {
	case []byte:
		_, err = gWriter.writer.Write(d)
	case byte:
		err = gWriter.writer.WriteByte(d)
	case string:
		_, err = gWriter.writer.WriteString(d)
		_, err = gWriter.writer.WriteString("\n")
	case rune:
		_, err = gWriter.writer.WriteRune(d)
	}
	err = gWriter.writer.Flush()
	return err
}
