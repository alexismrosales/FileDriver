package main

import (
	"cmd/client/internal"
	"log"
	"os"
)

func main() {
	// Here you can set your own path to save the logs of the APP, just
	// add a your preferred path
	logger, err := internal.NewLogger("~/.filedriver.log")
	if err != nil {
		log.Fatal("Error creating/reading loggerfile: ", err)
	}
	app := RunApp(logger)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
