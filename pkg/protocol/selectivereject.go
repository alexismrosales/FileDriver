package protocol

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/alexismrosales/FileDriver/pkg/storage"
)

type Answer struct {
	Message     *Message
	IsLastChunk bool
	ChunkIndex  int
}

// TODO: implements go routine to sent chunks of data with go routines

// SelectiveRejectSend implements the protocol of selective reject when the client sends
// data
func SelectiveRejectSend(chunks []storage.Chunk, windowSize int, conn *net.UDPConn, addr *net.UDPAddr) error {
	// variable to wait until the go routines are all finished
	var wg sync.WaitGroup

	leftPointer := 0
	rightPointer := windowSize - 1

	// continue while there are chunks to send
	for leftPointer < len(chunks) {
		// channel to collect number of packets
		packetsStatusQueue := make(chan struct {
			int
			bool
		})

		// adjust right pointer
		if rightPointer >= len(chunks) {
			rightPointer = len(chunks) - 1
		}

		// TODO: think about how to manage all packets, can i check if the packets in the window were sent? if not sent the rest
		// send all packets that windowSize allow
		for w := leftPointer; w <= rightPointer; w++ {
			wg.Add(1)
			go sendPacket(chunks[w], conn, addr, packetsStatusQueue)
		}
		// close channel
		close(packetsStatusQueue)
		// wait until all go routines end
		wg.Wait()

		// slide right window
		leftPointer += windowSize
		rightPointer += windowSize
	}
	return nil
}

// TODO: have in mind that there exists the possibility a state where connection is lost, add a timeout would be useful

// sendPacket to the other side and fill channel value expecting chunk index value
func sendPacket(chunk storage.Chunk, conn *net.UDPConn, addr *net.UDPAddr, status chan struct {
	int
	bool
}) {
	// Process every packet of the current window sequencially

	// Sent the chunk and wait answer
	err := SendAndWaitAck(conn, addr, chunk)

	// If package was sent succesully eval notSent chan
	if err != nil {
		//fmt.Errorf(err.Error())
		status <- struct {
			int
			bool
		}{chunk.ChunkIndex, true}
	} else {
		status <- struct {
			int
			bool
		}{chunk.ChunkIndex, false}
	}
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
		// TODO: create method to order chunks in an specific order
		chunks = append(chunks, *chunk)
	}
	return chunks, nil
}

// SendAndWaitAck sent a chunk of data and wait for an answer, if the answer is NAK
// packet was not send correctly, is necessary to resend it
func SendAndWaitAck(conn *net.UDPConn, addr *net.UDPAddr, chunk storage.Chunk) error {
	err := SendChunk(conn, chunk, addr)
	if err != nil {
		return errors.New("error while sending chunk")
	}
	conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	message, _, err := ReceiveMessage(conn)
	if err != nil {
		return errors.New("error  receiving packet")
	}

	if message.Type == MsgNak {
		return errors.New("nak received")
	}
	return nil
}
