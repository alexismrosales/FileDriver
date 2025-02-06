package main

import (
	"github.com/alexismrosales/FileDriver/pkg/logger"
	"log"
	"os"
)

const (
	COMMANDPARSERCTX = "COMMAND PARSER"
	CLICTX           = "CLI"
	CLIENTCTX        = "CLIENT"
)

func main() {
	app := RunApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func CreateLogger(context string) *internal.Logger {
	// Here you can set your own path to save the logs of the APP, just
	// add a your preferred path
	logger, err := internal.NewLogger(loggerFilePath, context)
	if err != nil {
		log.Fatal("Error creating/reading loggerfile: ", err)
	}
	return logger
}
