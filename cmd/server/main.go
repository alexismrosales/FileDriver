package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/alexismrosales/FileDriver/pkg/logger"
)

const (
	COMMANDPARSERCTX = "COMMAND PARSER"
	SERVERCTX        = "SERVER"
)

func main() {
	ip, port, err := askPort()
	if err != nil {
		panic(err)
	}
	fmt.Println("ip", ip, "port", port)
	server, err := NewServerUDP(ip, port)
	if err != nil {
		panic(err)
	}
	server.logger.Print("Server created succesfully", internal.INFO)
	defer server.conn.Close()
	server.ListenPetitions()
}

func askPort() (string, int, error) {
	if len(os.Args) < 3 {
		return "", 0, errors.New("Wrong args, you must assign the ip and port.")
	}
	args := os.Args
	if len(args[2]) == 4 {
		port, err := strconv.Atoi(args[2])
		if err != nil {
			return "", 0, errors.New("Invalid port")
		}
		return args[1], port, nil
	}
	// Returning default port
	return "", 8080, errors.New("Direction not setted correctly")
}
