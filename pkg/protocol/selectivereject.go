package protocol

import (
	"errors"
	"fmt"
	"github.com/alexismrosales/FileDriver/pkg/storage"
	"net"
	"sync"
	"time"
)

type Answer struct {
	Message     *Message
	Err         error
	IsLastChunk bool
}

// SelectiveRejectSend implements the protocol of selective reject when the client sends
// data
func SelectiveRejectSend(chunks []storage.Chunk, windowSize int, conn *net.UDPConn, addr *net.UDPAddr) error {
	var mu sync.Mutex
	packagesSent := make(map[int]struct{})
	leftPointer := 0
	rightPointer := windowSize - 1

	// Continuar mientras queden paquetes por enviar
	for leftPointer < len(chunks) {
		// Ajustar el puntero derecho si supera el límite de chunks
		if rightPointer >= len(chunks) {
			rightPointer = len(chunks) - 1
		}

		// Procesa cada chunk de la ventana actual secuencialmente
		for i := leftPointer; i <= rightPointer; i++ {
			chunk := chunks[i]

			// Verifica si el paquete ya fue enviado exitosamente
			mu.Lock()
			if _, chunkSent := packagesSent[chunk.ChunkIndex]; chunkSent {
				mu.Unlock()
				continue
			}
			mu.Unlock()

			// Envía el chunk y espera la respuesta
			answer := SentAndWaitAck(conn, addr, chunk)

			if answer.IsLastChunk && answer.Message.Type == MsgAck {
				break
			}
			if answer.Message != nil && answer.Message.Type == MsgAck {
				mu.Lock()
				packagesSent[chunk.ChunkIndex] = struct{}{}
				mu.Unlock()
			}

		}
		// Desplazar la ventana
		leftPointer += windowSize
		rightPointer += windowSize
	}
	return nil
}

func SelectiveRejectReceive(windowSize int, segmentSize int, conn *net.UDPConn) ([]storage.Chunk, error) {
	var chunks []storage.Chunk

	// The loop will stop when until all the packages have been received succesfully
	for len(chunks) < segmentSize {
		conn.SetReadDeadline(time.Now().Add(500 * time.Second))
		// Waiting for chunks
		chunk, clientAddr, err := ReceiveChunk(conn)
		// Waiting for chunks error
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, errors.New("Timeout, error recieving packages...")
		}
		if chunk.TotalChunks == 0 {
			break
		}
		// In case something wrong happen while receiving the chunk, NAK will be sent
		if err != nil {
			answerMessage := &Message{
				Type: MsgNak,
			}
			err := SendMessage(conn, answerMessage, clientAddr)
			if err != nil {
				continue
			}
		}
		if chunk == nil {
			return nil, errors.New("Received nil chunk")
		}
		// If the chunk is received succesfully a confirmation message will be sent
		answerMessage := &Message{
			Type:       MsgAck,
			ChunkIndex: chunk.ChunkIndex,
		}
		err = SendMessage(conn, answerMessage, clientAddr)
		if err != nil {
			continue
		}
		// Saving chunk without an specific order
		chunks = append(chunks, *chunk)
	}
	return chunks, nil
}

// SentAndWaitAck sent a chunk of data and wait for an answer, if the answer is NAK
// packet was not send correctly, is necessary to resend it
func SentAndWaitAck(conn *net.UDPConn, addr *net.UDPAddr, chunk storage.Chunk) *Answer {
	err := SendChunk(conn, chunk, addr)
	if err != nil {
		return &Answer{Message: nil, Err: err}
	}
	conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	message, _, err := ReceiveMessage(conn)
	if err != nil {
		fmt.Println("Error receiveing ")
		return &Answer{Message: nil, Err: err}
	}

	if message.Type == MsgNak {
		return &Answer{Message: message}
	}

	if message.Type == MsgAck {
		if chunk.TotalChunks == 0 {
			return &Answer{Message: message, IsLastChunk: true}
		}
		return &Answer{Message: message}
	}

	// If not message ack or nak
	return &Answer{Message: nil, Err: errors.New("unexpected message type")}
}
