package storage

import (
	"path"
	"path/filepath"
	"strings"
)

// Safely bytes size to send a datagram over the internet
const (
	ChunkSize = 850
)

type ChunkManager struct {
	fileStorage *FileStorage
}

type Chunk struct {
	FileName        string // Name of the file of the chunk
	TotalChunks     int    // Assigning total number of chunks of all files
	ChunkIndex      int    // Saving all the indexes of the chunks
	IsLastFileChunk bool   // Set flag to know if the chunk is the last one of the file
	Data            []byte // Chunk of data
}

func NewChunkManager() *ChunkManager {
	return &ChunkManager{fileStorage: NewFileStorage()}
}

// fragmentData creates an array of chunks, using the data structure "Chunk"
// every chunk is created depending of the numbers of files given by paths
// as a return the function gives the array of chunks of the bytes of the files
func (cm *ChunkManager) FragmentData(paths ...string) ([]Chunk, error) {
	var files [][]byte
	var chunks []Chunk
	// Saving all files converted in bytes
	for _, path := range paths {
		// Getting file in bytes format
		file, err := cm.fileStorage.ReadFromFile(path)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	chunkIndex := 0
	// Converting file by file into a chunk of arrays
	for index, file := range files {
		// Determine number of total chunks
		totalChunksFile := (len(file) + ChunkSize - 1) / ChunkSize

		for i := 0; i < totalChunksFile; i++ {
			// Setting start and end size of the chunk
			start := i * ChunkSize
			end := start + ChunkSize
			// Assign the rest of the file as the end in case the length of "end" is
			// greater than the length of the "file"
			if end > len(file) {
				end = len(file)
			}
			// Saving data in a single chunk with aditional information
			chunk := &Chunk{
				FileName:        filepath.Base(paths[index]),
				TotalChunks:     totalChunksFile,
				IsLastFileChunk: i == totalChunksFile-1,
				Data:            file[start:end],
			}
			// Increment index
			chunk.ChunkIndex = chunkIndex
			chunkIndex++

			// Saving every fragment into chunks of data
			chunks = append(chunks, *chunk)
		}
	}
	return chunks, nil
}

// DefragmentData receive an array SORTED of chunks.
// Is important to know that FileDriver Server APP will recieve datagrams
// in order due the UDP packet reciever algorithm used.
func (cm *ChunkManager) DefragmentData(chunks []Chunk, saveFilePath string) ([]*File, error) {
	var files []*File
	fileDataMap := make(map[string][]byte)
	// Use a temporal variable to stack an individual file on the array
	for _, chunk := range chunks {
		// Adding chunks to the data assumming are in order
		fileDataMap[chunk.FileName] = append(fileDataMap[chunk.FileName], chunk.Data...)
		if chunk.IsLastFileChunk {
			extension := path.Ext(chunk.FileName)
			name := strings.TrimSuffix(path.Base(chunk.FileName), extension)
			saveFilePath, err := GetShortPath(saveFilePath)
			if err != nil {
				return nil, err
			}
			file := &File{
				Name:      name,
				Extension: extension,
				Data:      fileDataMap[chunk.FileName],
				Path:      saveFilePath,
			}
			files = append(files, file)
			delete(fileDataMap, chunk.FileName)
		}
	}
	return files, nil
}
