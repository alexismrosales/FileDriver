package storage

type FileWriter interface {
	WriteToFile(filePath string, data []byte) error
}

// Reader provides read-only (without list) access to context data
type FileReader interface {
	ReadFromFile(filePath string) ([]byte, error)
}
