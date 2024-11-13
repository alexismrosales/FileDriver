package protocol

import (
	"github.com/alexismrosales/FileDriver/pkg/storage"
	"net"
)

// Trying handshake with server
func SendMessage(conn *net.UDPConn, message *Message, addr *net.UDPAddr) error {
	messageEncoded, err := EncodeMessage(*message)
	if err != nil {
		return err
	}
	// Sending message to server, waiting for response
	if conn.RemoteAddr() != nil {
		// If DialUP setted
		_, err = conn.Write(messageEncoded)
	} else {
		_, err = conn.WriteToUDP(messageEncoded, addr)
	}
	if err != nil {
		return err
	}
	return nil
}

func ReceiveMessage(conn *net.UDPConn) (*Message, *net.UDPAddr, error) {
	var n int
	var clientAddr *net.UDPAddr
	var err error
	buffer := make([]byte, 512)
	// Reading server answer
	if conn.RemoteAddr() != nil {
		n, err = conn.Read(buffer)

	} else {
		n, clientAddr, err = conn.ReadFromUDP(buffer)
	}

	if err != nil {
		return nil, nil, err
	}
	// Decoding message to json
	messageRecieved, err := DecodeMessage(buffer[:n])
	if err != nil {
		return nil, nil, err
	}
	return &messageRecieved, clientAddr, nil
}

func SendChunk(conn *net.UDPConn, chunk storage.Chunk, addr *net.UDPAddr) error {
	chunkEncoded, err := EncodeChunk(chunk)
	if err != nil {
		return err
	}
	if conn.RemoteAddr() != nil {
		// If DialUP setted
		_, err = conn.Write(chunkEncoded)
	} else {
		_, err = conn.WriteToUDP(chunkEncoded, addr)
	}
	if err != nil {
		return err
	}
	return nil
}

func ReceiveChunk(conn *net.UDPConn) (*storage.Chunk, *net.UDPAddr, error) {
	var n int
	var clientAddr *net.UDPAddr
	var err error
	buffer := make([]byte, storage.ChunkSize+500)
	// Reading server answer
	if conn.RemoteAddr() != nil {
		n, err = conn.Read(buffer)

	} else {
		n, clientAddr, err = conn.ReadFromUDP(buffer)
	}
	if err != nil {
		return nil, nil, err
	}

	chunkRecieved, err := DecodeChunk(buffer[:n])
	if err != nil {
		return nil, nil, err
	}
	return &chunkRecieved, clientAddr, nil
}
