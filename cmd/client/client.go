package main

import (
	"cmd/client/internal/logger"
	"cmd/client/internal/protocol"
	"cmd/client/internal/storage"
	"errors"
	"net"
	"path/filepath"
	"time"
)

const windowSize = 3

var handshakeFailed = errors.New("Error while connecting to the server...")

type Client struct {
	conn         *net.UDPConn
	logger       *internal.Logger
	chunkManager *storage.ChunkManager
}

func NewServerConnection(ip string, port int) (*Client, error) {
	// Create a new Client with UDP Protocol
	client, err := newClientUDP(ip, port)
	logger := client.logger
	if err != nil {
		logger.LogError(err)
	}
	// Message to confirm connection with server
	firstMessage := &protocol.Message{
		Type: protocol.MsgSyn,
	}
	// First time trying to stablish connection
	err = handshake(client.conn, firstMessage, client.logger)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	newServerInfo := info{
		Ip:   ip,
		Port: port,
	}
	err = saveInfo(newServerInfo)
	return client, err
}

func (c *Client) ServerConnection() (*Client, error) {
	serverInfo, err := getInfo()
	if err != nil {
		c.logger.LogError(errors.New("Error connecting to server..."))
		return nil, errors.New("Error connecting to server, check if configuration file exists if not try \"setaddr\" command")
	}
	return newClientUDP(serverInfo.Ip, serverInfo.Port)
}

func (c *Client) downloadFiles(paths ...string) error {
	logger := c.logger
	msgAns := &protocol.Message{
		Type:  protocol.MsgDownload,
		Paths: paths,
	}
	err := protocol.SendMessage(c.conn, msgAns, nil)
	if err != nil {
		return err
	}
	msgServerAns, _, err := protocol.ReceiveMessage(c.conn)
	if err != nil {
		return err
	}

	// Using protocol to receive all chuks of data from packets
	chunks, err := protocol.SelectiveRejectReceive(windowSize, msgServerAns.SegmentSize, c.conn)
	if err != nil {
		logger.LogError(errors.New("Error fragmenting data."))
		return err
	}

	files, err := c.chunkManager.DefragmentData(chunks, saveFilePath)
	if err != nil {
		logger.LogError(errors.New("Error defragmenting data."))
		return err
	}

	fileWriter := storage.NewFileStorage()
	for _, file := range files {
		path := filepath.Join(file.Path, file.Name+file.Extension)
		fileWriter.WriteToFile(path, file.Data, false)
	}
	return nil
}

func (c *Client) uploadFiles(paths ...string) error {
	msgAns := &protocol.Message{
		Type:  protocol.MsgUpload,
		Paths: paths,
	}
	err := protocol.SendMessage(c.conn, msgAns, nil)
	if err != nil {
		return err
	}
	err = getPathsFull(paths...)
	if err != nil {
		return err
	}
	// Fragment files and create an array of chunks
	chunks, err := c.chunkManager.FragmentData(paths...)
	if err != nil {
		return err
	}

	// Creating a message for the first time, letting the server know that
	// packets will be sent
	err = protocol.SendMessage(c.conn, &protocol.Message{Type: protocol.MsgUpload, SegmentSize: len(chunks)}, nil)
	if err != nil {
		return err
	}

	return protocol.SelectiveRejectSend(chunks, windowSize, c.conn, nil)
	// After confirm connection, start Selective Repeat Protocol
}

func (c *Client) sendCommand(command string, paths []string, flags []string) error {
	logger := c.logger
	messageCmd := &protocol.Message{
		Type:  command,
		Paths: paths,
		Flags: flags,
	}
	err := protocol.SendMessage(c.conn, messageCmd, nil)
	if err != nil {
		logger.LogError(errors.New("Error sending command to server."))
		return err
	}
	c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	answer, _, err := protocol.ReceiveMessage(c.conn)
	if err != nil {
		return err
	}

	if answer.Type == protocol.MsgError {
		return errors.New(answer.Error)
	}
	logger.LogMessage("Output recieved:" + answer.Output)
	return nil
}

// NewClient constructor using ip and port
func newClientUDP(ip string, port int) (*Client, error) {
	severAddr := net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	}
	// Creating connection with udp protocol
	conn, err := net.DialUDP("udp", nil, &severAddr)
	if err != nil {
		return nil, err
	}

	//  Creating a new logger for CLIENT context
	logger, err := internal.NewLogger(loggerFilePath, CLIENTCTX)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:         conn,
		logger:       logger,
		chunkManager: storage.NewChunkManager(),
	}, nil
}

func handshake(conn *net.UDPConn, message *protocol.Message, logger *internal.Logger) error {
	// Sending SYN to the sever to try HandShake
	err := protocol.SendMessage(conn, message, nil)
	if err != nil {
		return err
	}
	// Wait connection for 5 seconds
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	messageRecieved, _, err := protocol.ReceiveMessage(conn)
	if err != nil {
		return err
	}
	// In case the response is wrong or different of the desired
	// The expeceted message has a type like ACK and SYN-ACK
	if messageRecieved.Type != protocol.MsgAck {
		return handshakeFailed
	}
	// Modifying segment size
	message.SegmentSize = messageRecieved.SegmentSize
	logger.LogMessage("Handshake succesfully with server!")
	return nil
}
