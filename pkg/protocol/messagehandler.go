package protocol

import (
	"encoding/json"
	"github.com/alexismrosales/FileDriver/pkg/storage"
)

func EncodeMessage(data Message) ([]byte, error) {
	return json.Marshal(data)
}

func DecodeMessage(data []byte) (Message, error) {
	var message Message
	err := json.Unmarshal(data, &message)
	return message, err
}

func EncodeChunk(chunk storage.Chunk) ([]byte, error) {
	return json.Marshal(chunk)
}

func DecodeChunk(data []byte) (storage.Chunk, error) {
	var chunk storage.Chunk
	err := json.Unmarshal(data, &chunk)
	return chunk, err
}
