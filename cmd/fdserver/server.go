package main

import (
	"github.com/alexismrosales/FileDriver/pkg/logger"
	"github.com/alexismrosales/FileDriver/pkg/protocol"
	"github.com/alexismrosales/FileDriver/pkg/storage"
	"net"
	"strings"
	"time"
)

var windowSize = 3

type Server struct {
	conn         *net.UDPConn
	addr         *net.UDPAddr
	logger       *internal.Logger
	chunkManager *storage.ChunkManager
}

func NewServerUDP(ip string, port int) (*Server, error) {
	addr := &net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	logger, err := internal.NewLogger(loggerFilePath, SERVERCTX)
	return &Server{
		conn:         conn,
		addr:         addr,
		logger:       logger,
		chunkManager: storage.NewChunkManager(),
	}, nil
}

// ListenPetitions start listening petitions from the client
func (server *Server) ListenPetitions() error {
	var timeoutCount int
	logger := server.logger
	// First Handshake with client
	server.handshake()
	// Create a directory called server/ in home dir
	err := createBaseDir()

	if err != nil {
		return err
	}
	err = saveStatus("/")
	if err != nil {
		return err
	}
	logger.LogMessage("Now listening petitions from: " + server.conn.LocalAddr().String())
	// Start listening petitions from the client, the loop receive messages from client
	// and answer depending of the command

	for {
		server.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		message, clientAddr, err := protocol.ReceiveMessage(server.conn)
		server.addr = clientAddr
		if err != nil {
			// Ignorar errores específicos de desconexión de cliente
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				timeoutCount++
				if timeoutCount == 5 { // Solo muestra el mensaje una vez si no hay peticiones
					logger.Print("Timeout reached, no new requests. Waiting...", internal.DEBUG)
				}
				continue
			} else if strings.Contains(err.Error(), "connection reset by peer") {
				logger.Print("Client disconnected unexpectedly, waiting for new connection...", internal.INFO)
				continue
			}
			logger.Print("Error receiving messages", internal.ERROR, err)
			continue
		}
		logger.LogMessage("New message recieved succesfully.")
		logger.LogMessage("Command '" + message.Type + "': to be executed")
		output, err := server.handleCommands(message)
		if err != nil {
			answer := &protocol.Message{
				Type:   protocol.MsgError, // Send Ack confirmation to client
				Output: output,
				Error:  err.Error(),
			}
			logger.Print("Error executing command: "+err.Error(), internal.ERROR, err)
			protocol.SendMessage(server.conn, answer, clientAddr)

		} else {
			answer := &protocol.Message{
				Type:   protocol.MsgAck, // Send Ack confirmation to client
				Output: output,
			}
			protocol.SendMessage(server.conn, answer, clientAddr)
			logger.LogMessage("Command executed and message sent to client.")
		}
	}
}

func (server *Server) handshake() error {
	logger := server.logger
	_, clientAddr, err := protocol.ReceiveMessage(server.conn)
	if err != nil {
		server.logger.Print("Error receiving messages", internal.ERROR)
		return err
	}
	answer := &protocol.Message{
		Type: protocol.MsgAck, // Send Ack confirmation to client
	}
	protocol.SendMessage(server.conn, answer, clientAddr)
	logger.LogMessage("Command executed and message sent to client.")
	return nil
}

func (s *Server) handleCommands(message *protocol.Message) (string, error) {
	currentPath, err := getStatusPath()
	if err != nil {
		return "", err
	}
	switch message.Type {
	case protocol.MsgDownload:
		return "", download(currentPath, message.Paths, windowSize, s.conn, s.addr)
	case protocol.MsgUpload:
		return "", upload(currentPath, windowSize, s.conn)
	case protocol.MsgPwd:
		return pwd(currentPath), nil
	case protocol.MsgMkdir:
		err = mkdir(currentPath, message.Paths...)
		return "", nil
	case protocol.MsgCd:
		// Get path using cd command
		newCurrentPath, err := cd(currentPath, message.Paths[0])
		if err != nil {
			return "", err
		}
		// Write the new path in file
		err = saveStatus(newCurrentPath)
		return newCurrentPath, err
	case protocol.MsgLs:
		return ls(currentPath, message.Flags)
	case protocol.MsgMv:
		return "", mv(currentPath, message.Paths...)
	case protocol.MsgRm:
		return "", rm(currentPath, message.Paths, message.Flags)
	}
	return "", nil
}
