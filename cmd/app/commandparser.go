package main

import (
	"errors"
	"github.com/alexismrosales/FileDriver/pkg/logger"
	"strconv"
	"strings"
)

const (
	CmdDownload = "DOWNLOAD"
	CmdUpload   = "UPLOAD"
	CmdPwd      = "PWD"
	CmdMkdir    = "MKDIR"
	CmdLs       = "LS"
	CmdCd       = "CD"
	CmdMv       = "MV"
	CmdRm       = "RM"
)

// FileManger gets the current directory from the server
type CommandParser struct {
	client     *Client
	currentDir string
	flags      []string
	logger     *internal.Logger
}

func NewCommandParser(currentDir string, flags ...string) *CommandParser {
	return &CommandParser{
		currentDir: currentDir,
		flags:      flags,
		logger:     CreateLogger(COMMANDPARSERCTX),
	}
}

// Connect works to save the address and port to be use later
func (cp *CommandParser) FirstConnection(address, port string) error {
	intPort, err := strconv.Atoi(port)
	if err != nil {
		return err
	}
	client, err := NewServerConnection(address, intPort)
	if err != nil {
		return err
	}
	cp.client = client
	return nil
}

// Disconnect close the connection between client and server
func (cp *CommandParser) Disconnect() error {
	return cp.client.conn.Close()
}

// Commands send to execute tasks to the server
func (cp *CommandParser) Pwd() error {
	return cp.executeCommand(CmdPwd)
}

func (cp *CommandParser) Mkdir(paths ...string) error {
	return cp.executeCommand(CmdMkdir, paths...)
}

func (cp *CommandParser) Ls(paths ...string) error {
	return cp.executeCommand(CmdLs, paths...)
}

func (cp *CommandParser) Cd(path string) error {
	return cp.executeCommand(CmdCd, path)
}

func (cp *CommandParser) Mv(paths ...string) error {
	return cp.executeCommand(CmdMv, paths...)
}
func (cp *CommandParser) Rm(paths ...string) error {
	return cp.executeCommand(CmdRm, paths...)
}

// Upload send all the files to the server
func (cp *CommandParser) Upload(paths ...string) error {
	cp.logCommand(CmdUpload, paths...)
	client, err := cp.client.ServerConnection()
	if err != nil {
		return err
	}
	defer client.conn.Close()
	cp.client = client
	if err := cp.client.uploadFiles(paths...); err != nil {
		cp.logger.LogError(err)
		return err
	}
	cp.logger.Print("Upload succesfully completed", internal.INFO)
	return nil

}

// Download works for every selected file
func (cp *CommandParser) Download(paths ...string) error {
	cp.logCommand(CmdDownload, paths...)
	// Connect with client
	client, err := cp.client.ServerConnection()
	if err != nil {
		err = errors.New("Connection error with server.")
	}
	defer client.conn.Close()
	cp.client = client
	// Send message with command
	if err := cp.client.downloadFiles(paths...); err != nil {
		cp.logger.LogError(err)
		return err
	}
	cp.logger.Print("Download finished succesfully", internal.DEBUG)
	return nil

}

// executeCommand and log the command and tasks done
func (cp *CommandParser) executeCommand(command string, paths ...string) error {
	cp.logCommand(command, paths...)
	// Connect with client
	client, err := cp.client.ServerConnection()
	defer client.conn.Close()
	if err != nil {
		err = errors.New("Connection error with server.")
	}
	cp.client = client
	// Send message with command
	if err := cp.client.sendCommand(command, paths, cp.flags); err != nil {
		cp.logger.LogError(err)
		return err
	}
	cp.logger.Print("Command executed succesfully", internal.DEBUG)
	return nil
}

// logCommand show the command on log with formated text
func (cp *CommandParser) logCommand(command string, paths ...string) {
	mergedArgs := strings.Join(cp.flags, " ") + strings.Join(paths, " ")
	cp.logger.Print("Command ["+command+" "+mergedArgs+"]", internal.INFO)
}
