package protocol

import (
	"cmd/client/internal/storage"
	"encoding/json"
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
