package internal

import (
	"bufio"
	"os"
)

type GenericWriter struct {
	writer *bufio.Writer
}

// NewWriter creates a new write buffer
func NewGWriter(file *os.File) *GenericWriter {
	writer := bufio.NewWriter(file)
	return &GenericWriter{writer: writer}
}

// Writer selects the type of data to writes it on the file
func (gWriter *GenericWriter) Writer(data any) {
	switch d := data.(type) {
	case []byte:
		gWriter.writer.Write(d)
	case byte:
		gWriter.writer.WriteByte(d)
	case string:
		gWriter.writer.WriteString(d)
	case rune:
		gWriter.writer.WriteRune(d)
	}
	gWriter.writer.WriteString("\n")
	gWriter.writer.Flush()
}
