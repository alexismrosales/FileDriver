package main

import internal "cmd/client/internal/logger"

const (
	COMMANDPARSERCTX = "COMMAND PARSER"
	SERVERCTX        = "SERVER"
)

func main() {
	server, err := NewServerUDP(askPort())
	if err != nil {
		return
	}
	server.logger.Print("Server created succesfully", internal.INFO)
	defer server.conn.Close()
	server.ListenPetitions()
}

func askPort() (string, int) {
	return "127.0.0.1", 8080
}
